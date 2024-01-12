/*
@Author: 梦无矶小仔
@Date:   2024/1/12 18:23
*/
package jwt

import (
	"gorm.io/gorm"
	"time"
)

type JwtBlacklist struct {
	ID        uint           `gorm:"primarykey;comment:主键ID"` // 主键ID
	CreatedAt time.Time      `gorm:"type:datetime(0);comment:创建时间" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"type:datetime(0);comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"` // 删除时间
	Jwt       string         `gorm:"type:text;comment:jwt"`
}
