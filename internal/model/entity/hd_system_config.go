// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// HdSystemConfig is the golang structure for table hd_system_config.
type HdSystemConfig struct {
	Id          uint64      `json:"id"          ` // id
	Name        string      `json:"name"        ` // 配置名称
	Key         string      `json:"key"         ` // 配置key
	Value       string      `json:"value"       ` // 配置值
	Description string      `json:"description" ` // 描述
	CreatedAt   *gtime.Time `json:"createdAt"   ` // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   ` // 更新时间
}
