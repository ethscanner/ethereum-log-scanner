package scanner

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/frame/g"
)

const Max_Concurrent_Threads = 3

type scanner struct {
	name string
	//保存日志的存储器
	logStorage         LogStorage
	storage            ScannerStorage
	addresses          []common.Address
	scannedBlockNumkey string
	queryTopics        [][]common.Hash
	//重复扫块的偏移值(OverridePerScan和DelayBlocks其中有一项打开就可以避免reorg时事件丢失)
	OverridePerScan uint64
	//延迟扫块的偏移值(OverridePerScan和DelayBlocks其中有一项打开就可以避免reorg时事件丢失)
	DelayBlocks uint64
	//每次扫块请求的最大间隔数
	IntervalPerScan    uint64
	SegmentationLength uint64 //分段查询数量
}

type queryResult struct {
	from uint64
	to   uint64
	logs []types.Log
	err  error
}

func NewScanner(name string, addresses []common.Address, storage ScannerStorage, logStorage LogStorage, queryTopics [][]common.Hash) *scanner {
	var IntervalPerScan uint64 = 500
	return &scanner{
		name:               name,
		scannedBlockNumkey: name + "BlockNumber",
		addresses:          addresses,
		OverridePerScan:    0,
		IntervalPerScan:    IntervalPerScan,
		DelayBlocks:        uint64(3),
		queryTopics:        queryTopics,
		storage:            storage,
		logStorage:         logStorage,
		SegmentationLength: IntervalPerScan * Max_Concurrent_Threads,
	}
}

// 获取扫块from区块和to区块
func (s *scanner) GetScanBlockNumbersByStroage(ctx context.Context, expectToBlockNum uint64) (from uint64, to uint64, err error) {
	if expectToBlockNum >= s.DelayBlocks {
		to = expectToBlockNum - s.DelayBlocks
	}
	if s.storage == nil {
		return 0, 0, errors.New("storage is nil")
	}
	if from, err = s.storage.GetUint64ByKey(ctx, s.scannedBlockNumkey); err != nil {
		return
	}

	if from == 0 {
		from = to
		s.InitScannedBlockNum(ctx, expectToBlockNum)
	}
	if from > s.OverridePerScan {
		from = from - s.OverridePerScan
	}
	return from, to, nil

}

func (s *scanner) ScanToStroage(ctx context.Context, client *ethclient.Client, to uint64) (scannedBlockNum uint64, err error) {
	//如果to传入为0,则获取最新区块
	if to == 0 {
		if to, err = client.BlockNumber(ctx); err != nil {
			return 0, err
		}
	}
	//获取偏移后的from和to
	from, to, err := s.GetScanBlockNumbersByStroage(ctx, to)
	if err != nil {
		return 0, err
	}
	if from == to {
		return 0, err
	}
	if _, err = s.SegmentationScanToStroage(ctx, client, from, to); err != nil {
		return 0, err
	}
	return 0, err
}

// 分段扫描(不包含from)并保存日志
func (s *scanner) SegmentationScanToStroage(ctx context.Context, client *ethclient.Client, from, to uint64) (scannedBlockNum uint64, err error) {
	g.Log().Infof(ctx, "开始扫描 SegmentationScanToStroage %v:%d - %d", s.name, from, to)
	var logs []types.Log
	var expectTo uint64 = from
	for expectTo < to {
		expectTo += s.SegmentationLength
		from = expectTo - s.SegmentationLength
		if expectTo > to {
			expectTo = to
		}
		g.Log().Infof(ctx, "开始扫描 %v:%d - %d", s.name, from, expectTo)
		logs, scannedBlockNum, err = s.ConcurrentThreadQuery(ctx, client, from, expectTo)
		if err != nil {
			return 0, err
		} else if scannedBlockNum != expectTo {
			return 0, fmt.Errorf("ConcurrentThreadQuery error: %w ", err)
		}
		//保存日志
		if err := s.saveLogs(ctx, logs); err != nil {
			return 0, err
		}
		//保存查询
		if err := s.saveScannedBlockNum(ctx, scannedBlockNum); err != nil {
			return 0, err
		}
	}
	return scannedBlockNum, nil
}

// 分段扫描(不包含from)并返回,注意数据可能很大
func (s *scanner) SegmentationScan(ctx context.Context, client *ethclient.Client, from, to uint64) (allLogs []types.Log, scannedBlockNum uint64, err error) {
	g.Log().Infof(ctx, "开始扫描 %v:%d - %d", s.name, from, to)
	var logs []types.Log
	var expectTo uint64 = from
	for expectTo < to {
		expectTo += s.SegmentationLength
		from = expectTo - s.SegmentationLength
		if expectTo > to {
			expectTo = to
		}
		logs, scannedBlockNum, err = s.ConcurrentThreadQuery(ctx, client, from, expectTo)
		if err != nil {
			return nil, 0, err
		} else if scannedBlockNum != expectTo {
			return nil, 0, fmt.Errorf("ConcurrentThreadQuery error: %w ", err)
		}
		allLogs = append(allLogs, logs...)
	}
	return allLogs, scannedBlockNum, nil
}

// 多线程查询(不包含from区块)
func (s *scanner) ConcurrentThreadQuery(ctx context.Context, client *ethclient.Client, from, to uint64) (logs []types.Log, scannedBlockNum uint64, err error) {
	ranges := to - from
	if ranges < s.IntervalPerScan {
		logs, err = s.FilterQuery(ctx, client, from+1, to)
		return logs, to, err
	} else if ranges > Max_Concurrent_Threads*s.IntervalPerScan { //如果大于最大线程数则截取
		return nil, 0, errors.New("超出最大限制数")
	}

	//计算需要用到的线程数
	remainder := ranges % s.IntervalPerScan
	threadCount := int(ranges / s.IntervalPerScan)
	if remainder != 0 {
		threadCount += 1
	}
	//并发操作
	wg := sync.WaitGroup{}
	wg.Add(threadCount)
	//通道
	ch := make(chan queryResult, threadCount)
	var expectTo uint64 = from
	for expectTo < to {
		expectTo += s.IntervalPerScan
		from = expectTo - s.IntervalPerScan + 1
		if expectTo > to {
			expectTo = to
		}
		go func(from, to uint64) {
			logs, err := s.FilterQuery(ctx, client, from, to)
			ret := queryResult{from, to, logs, err}
			ch <- ret
			wg.Done()
		}(from, expectTo)
	}
	wg.Wait()
	close(ch)

	//更改
	for v := range ch {
		if v.err != nil {
			return nil, 0, v.err
		}
		if v.to > scannedBlockNum {
			scannedBlockNum = v.to
		}
		logs = append(logs, v.logs...)
	}

	return
}

func (s *scanner) FilterQuery(ctx context.Context, client *ethclient.Client, from, to uint64) (logs []types.Log, err error) {
	query := ethereum.FilterQuery{
		ToBlock:   new(big.Int).SetUint64(to),
		FromBlock: new(big.Int).SetUint64(from),
		Addresses: s.addresses,
		Topics:    s.queryTopics,
	}
	logs, err = client.FilterLogs(ctx, query)
	return
}

// 保存日志
func (s *scanner) saveLogs(ctx context.Context, logs []types.Log) error {
	if s.logStorage != nil {
		return s.logStorage.SaveLogs(ctx, s.name, logs)
	}
	return nil
}

// 记录已经扫描完毕的区块位置
func (s *scanner) saveScannedBlockNum(ctx context.Context, to uint64) error {
	if s.storage != nil {
		return s.storage.UpdateUint64(ctx, s.scannedBlockNumkey, to)
	}
	return nil
}

func (s *scanner) InitScannedBlockNum(ctx context.Context, to uint64) error {
	if s.storage != nil {
		return s.storage.InsertUint64(ctx, s.scannedBlockNumkey, to)
	}
	return nil
}
