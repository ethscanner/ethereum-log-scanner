package scanner

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
)

type checker struct {
	name string
	//保存日志的存储器
	logStorage         DbLogStorage
	SegmentationLength int //分段查询数量
	DelayBlocks        uint64
	LookbackBlocks     uint64 // 回溯区块数量
}

func NewChecker(name string, logStorage DbLogStorage) *checker {
	return &checker{
		name:               name,
		logStorage:         logStorage,
		DelayBlocks:        15,
		SegmentationLength: 1000,
		LookbackBlocks:     10000000,
	}
}

func (s *checker) CheckAllStroage(ctx context.Context, client *ethclient.Client, blockNumber uint64) (scannedBlockNum uint64, err error) {
	var CheckState int = int(CHECK_STATE_PENDING)
	if blockNumber < s.DelayBlocks {
		return 0, nil
	}
	time := gtime.Now().Add(-time.Hour * 24 * 30)
	ltBlockNumber := int64(blockNumber - s.DelayBlocks)
	query := LogQuery{
		ContractName:  s.name,
		CheckState:    &CheckState,
		Limit:         s.SegmentationLength,
		BlockNumberLt: &ltBlockNumber,
		OrderBy:       "block_number",
		Desc:          false,
		CreatedAtGt:   time,
	}
	checkedBlockSet := gset.NewStrSet()
	errBlockSet := gset.NewStrSet()
	if logs, err := s.logStorage.QueryLogs(ctx, query); err != nil {
		return 0, err
	} else {
		for _, v := range logs {
			//startTime := gtime.TimestampMilli()
			check, err := s.CheckBlockHash(ctx, client, v)
			if err != nil {
				return 0, err
			}
			v.CheckedBlock = blockNumber
			if check {
				checkedBlockSet.Add(v.BlockHash.Hex())
			} else {
				errBlockSet.Add(v.BlockHash.Hex())
			}
			//costTime := gtime.TimestampMilli() - startTime
			//g.Log().Infof(ctx, "检查%v区块%v耗时%vms", v.BlockHash.Hex(), blockNumber, costTime)
		}
	}
	checkedBlockSlice := checkedBlockSet.Slice()
	errBlockSlice := errBlockSet.Slice()

	if len(checkedBlockSlice) > 0 {
		if err := s.logStorage.UpdateBlockCheckState(ctx, checkedBlockSlice, blockNumber, 1); err != nil {
			return 0, err
		}
	}
	if len(errBlockSlice) > 0 {
		if err := s.logStorage.UpdateBlockCheckState(ctx, errBlockSlice, blockNumber, 2); err != nil {
			return 0, err
		}
	}
	return blockNumber, nil
}

func (s *checker) CheckBlockHash(ctx context.Context, client *ethclient.Client, log Elog) (success bool, err error) {
	cacheKey := "hash_block_number_" + log.BlockHash.Hex()
	cacheValue, err := gcache.GetOrSetFunc(ctx, cacheKey, func(context.Context) (interface{}, error) {
		_, err := client.BlockByHash(ctx, log.BlockHash)
		if err == ethereum.NotFound {
			return false, nil
		} else if err != nil {
			return false, err
		}
		return true, nil
	}, 5*time.Minute)
	if err != nil {
		return false, err
	}
	return cacheValue.Bool(), nil
}
