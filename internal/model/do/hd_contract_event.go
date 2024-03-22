// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// HdContractEvent is the golang structure of table hd_contract_event for DAO operations like Where/Data.
type HdContractEvent struct {
	g.Meta          `orm:"table:hd_contract_event, do:true"`
	Id              interface{} // id
	ContractName    interface{} // 合约名
	ContractAddress interface{} // 合约地址
	TxHash          interface{} // 交易哈希
	EventHash       interface{} // 事件名
	EventId         interface{} // 事件id
	BlockNumber     interface{} // 区块编号
	BlockHash       interface{} // 交易哈希
	EventTopics     []byte      // 事件头
	EventData       []byte      // event数据
	State           interface{} // 0:待处理 10:已处理
	CheckState      interface{} // 链上状态: 0:待处理 10:已确认 20:确认异常
	CheckedBlock    interface{} // 已确认区块
	CreatedAt       *gtime.Time // 创建时间
	UpdatedAt       *gtime.Time // 更新时间
}
