/*
@Author: 梦无矶小仔
@Date:   2024/1/30 11:40
*/
package sys

import (
	"time"
)

type SysApis struct {
	ID          uint      `gorm:"primarykey;comment:主键ID" json:"id" structs:"-"` // 主键ID
	CreatedAt   time.Time `gorm:"type:datetime(0);autoCreateTime;comment:创建时间" json:"createdAt" structs:"-"`
	UpdatedAt   time.Time `gorm:"type:datetime(0);autoUpdateTime;comment:更新时间" json:"updatedAt" structs:"-"`
	Title       string    `json:"title" gorm:"comment:api路径名称"`          // api路径
	Path        string    `json:"path" gorm:"comment:api路径"`             // api路径
	Description string    `json:"description" gorm:"comment:api中文描述"`    // api中文描述
	ParentId    uint      `json:"parentId" gorm:"comment:隶属于父ID"`        // api组
	Method      string    `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	Code        string    `json:"code" gorm:"comment:权限代号"`              // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	// 忽略该字段，- 表示无读写，-:migration 表示无迁移权限，-:all 表示无读写迁移权限
	Children []*SysApis `gorm:"-" json:"children"`
}

func (s *SysApis) TableName() string {
	return "sys_apis"
}
