package scanner

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// 日志id
type ELogId struct {
	TxHash    common.Hash `json:"transactionHash" `
	BlockHash common.Hash `json:"blockHash"`
	Index     uint        `json:"logIndex" `
}

type Elog struct {
	types.Log
	Id           *uint64 `json:"id"              ` // id
	ContractName string  `json:"contractName"    ` // 合约名
	CheckState   int     `json:"checkState"      ` // 链上状态: 0:待处理 10:已确认 20:确认异常
	CheckedBlock uint64  `json:"checkedBlock"    ` // 已确认区块

}

type LogQuery struct {
	IdGt            int64  `json:"idGt"         `    // id大于
	ContractName    string `json:"contractName"    ` // 合约名
	ContractAddress string `json:"contractAddress" ` // 合约地址
	TxHash          string `json:"txHash"          ` // 交易哈希
	EventHash       string `json:"eventHash"       ` // 事件名
	EventId         int64  `json:"eventId"         ` // 事件id
	BlockNumber     int64  `json:"blockNumber"     ` // 区块编号
	CheckState      int    `json:"checkState"      ` // 链上状态: 0:待处理 10:已确认 20:确认异常
	State           int    `json:"state"      `      // 链上状态: 0:待处理 10:已处理
	Limit           int
}

// 日志存储器
type LogStorage interface {
	//保存日志
	SaveLogs(ctx context.Context, name string, logs []types.Log) error
}

type DbLogStorage interface {
	//保存日志
	SaveLogs(ctx context.Context, name string, logs []types.Log) error
	//批量标记已处理
	MarkAsProcessed(ctx context.Context, name string, logs []ELogId) error

	QueryLogs(ctx context.Context, query LogQuery) (logs []Elog, err error)

	UpdateBlockCheckState(ctx context.Context, log Elog) error

	GetLogByElogId(ctx context.Context, id ELogId) (log Elog, err error)
}

// 扫描存储器
type ScannerStorage interface {
	UpdateUint64(ctx context.Context, key string, val uint64) error
	GetUint64ByKey(ctx context.Context, key string) (val uint64, err error)
	InsertUint64(ctx context.Context, key string, val uint64) error
}
