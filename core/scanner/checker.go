package scanner

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
)

type checker struct {
	name string
	//保存日志的存储器
	logStorage         DbLogStorage
	SegmentationLength int //分段查询数量
	DelayBlocks        uint64
}

func NewChecker(name string, logStorage DbLogStorage) *checker {
	return &checker{
		name:               name,
		logStorage:         logStorage,
		DelayBlocks:        15,
		SegmentationLength: 500,
	}
}

func (s *checker) CheckAllStroage(ctx context.Context, client *ethclient.Client, blockNumber uint64) (scannedBlockNum uint64, err error) {
	//如果to传入为0,则获取最新区块
	if blockNumber == 0 {
		if blockNumber, err = client.BlockNumber(ctx); err != nil {
			return 0, err
		}
	}
	var CheckState int = 0
	query := LogQuery{
		ContractName: s.name,
		CheckState:   &CheckState,
		Limit:        s.SegmentationLength,
	}
	if logs, err := s.logStorage.QueryLogs(ctx, query); err != nil {
		return 0, err
	} else {
		for _, v := range logs {
			if err, success := s.CheckLog(ctx, client, v); err != nil {
				return blockNumber, err
			} else if success {
				v.CheckedBlock = blockNumber
				if v.BlockNumber+s.DelayBlocks < v.CheckedBlock {
					v.CheckState = 1
				}
			} else {
				v.CheckedBlock = blockNumber
				v.CheckState = 2
			}
			if err := s.logStorage.UpdateBlockCheckState(ctx, v); err != nil {
				return blockNumber, err
			}
		}
	}
	return blockNumber, err
}

func (s *checker) CheckLog(ctx context.Context, client *ethclient.Client, log Elog) (err error, success bool) {
	if tx, err := client.TransactionReceipt(ctx, log.TxHash); err != nil {
		if err == ethereum.NotFound {
			return nil, false
		}
		return err, false
	} else if tx.Status == 1 && tx.BlockHash.Hex() == log.BlockHash.Hex() {
		return nil, true
	} else {
		return nil, false
	}
}
