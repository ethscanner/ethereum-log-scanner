// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// HdContractEventDao is the data access object for table hd_contract_event.
type HdContractEventDao struct {
	table   string                 // table is the underlying table name of the DAO.
	group   string                 // group is the database configuration group name of current DAO.
	columns HdContractEventColumns // columns contains all the column names of Table for convenient usage.
}

// HdContractEventColumns defines and stores column names for table hd_contract_event.
type HdContractEventColumns struct {
	Id              string // id
	ContractName    string // 合约名
	ContractAddress string // 合约地址
	TxHash          string // 交易哈希
	EventHash       string // 事件名
	EventId         string // 事件id
	BlockNumber     string // 区块编号
	BlockHash       string // 交易哈希
	EventTopics     string // 事件头
	EventData       string // event数据
	State           string // 0:待处理 10:已处理
	CheckState      string // 链上状态: 0:待处理 10:已确认 20:确认异常
	CheckedBlock    string // 已确认区块
	CreatedAt       string // 创建时间
	UpdatedAt       string // 更新时间
}

// hdContractEventColumns holds the columns for table hd_contract_event.
var hdContractEventColumns = HdContractEventColumns{
	Id:              "id",
	ContractName:    "contract_name",
	ContractAddress: "contract_address",
	TxHash:          "tx_hash",
	EventHash:       "event_hash",
	EventId:         "event_id",
	BlockNumber:     "block_number",
	BlockHash:       "block_hash",
	EventTopics:     "event_topics",
	EventData:       "event_data",
	State:           "state",
	CheckState:      "check_state",
	CheckedBlock:    "checked_block",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
}

// NewHdContractEventDao creates and returns a new DAO object for table data access.
func NewHdContractEventDao() *HdContractEventDao {
	return &HdContractEventDao{
		group:   "default",
		table:   "hd_contract_event",
		columns: hdContractEventColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *HdContractEventDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *HdContractEventDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *HdContractEventDao) Columns() HdContractEventColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *HdContractEventDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *HdContractEventDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *HdContractEventDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
