// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// HdSystemConfigDao is the data access object for table hd_system_config.
type HdSystemConfigDao struct {
	table   string                // table is the underlying table name of the DAO.
	group   string                // group is the database configuration group name of current DAO.
	columns HdSystemConfigColumns // columns contains all the column names of Table for convenient usage.
}

// HdSystemConfigColumns defines and stores column names for table hd_system_config.
type HdSystemConfigColumns struct {
	Id          string // id
	Name        string // 配置名称
	Key         string // 配置key
	Value       string // 配置值
	Description string // 描述
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// hdSystemConfigColumns holds the columns for table hd_system_config.
var hdSystemConfigColumns = HdSystemConfigColumns{
	Id:          "id",
	Name:        "name",
	Key:         "key",
	Value:       "value",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewHdSystemConfigDao creates and returns a new DAO object for table data access.
func NewHdSystemConfigDao() *HdSystemConfigDao {
	return &HdSystemConfigDao{
		group:   "default",
		table:   "hd_system_config",
		columns: hdSystemConfigColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *HdSystemConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *HdSystemConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *HdSystemConfigDao) Columns() HdSystemConfigColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *HdSystemConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *HdSystemConfigDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *HdSystemConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
