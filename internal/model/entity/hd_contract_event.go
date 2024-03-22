// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// HdContractEvent is the golang structure for table hd_contract_event.
type HdContractEvent struct {
	Id              *uint64      `json:"id"              ` // id
	ContractName    string      `json:"contractName"    ` // 合约名
	ContractAddress string      `json:"contractAddress" ` // 合约地址
	TxHash          string      `json:"txHash"          ` // 交易哈希
	EventHash       string      `json:"eventHash"       ` // 事件名
	EventId         int64       `json:"eventId"         ` // 事件id
	BlockNumber     int64       `json:"blockNumber"     ` // 区块编号
	BlockHash       string      `json:"blockHash"       ` // 交易哈希
	EventTopics     []byte      `json:"eventTopics"     ` // 事件头
	EventData       []byte      `json:"eventData"       ` // event数据
	State           int         `json:"state"           ` // 0:待处理 10:已处理
	CheckState      int         `json:"checkState"      ` // 链上状态: 0:待处理 10:已确认 20:确认异常
	CheckedBlock    uint64      `json:"checkedBlock"    ` // 已确认区块
	CreatedAt       *gtime.Time `json:"createdAt"       ` // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       ` // 更新时间
}
