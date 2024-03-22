// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// HdSystemConfig is the golang structure of table hd_system_config for DAO operations like Where/Data.
type HdSystemConfig struct {
	g.Meta      `orm:"table:hd_system_config, do:true"`
	Id          interface{} // id
	Name        interface{} // 配置名称
	Key         interface{} // 配置key
	Value       interface{} // 配置值
	Description interface{} // 描述
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
