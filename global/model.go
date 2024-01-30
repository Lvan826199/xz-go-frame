/*
@Author: 梦无矶小仔
@Date:   2024/1/30 16:28
*/
package global

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

/*修改逻辑删除的默认规则
go get gorm.io/plugin/soft_delete
*/

type GVA_MODEL struct {
	ID        uint      `gorm:"primarykey;comment:主键ID" json:"id" structs:"-"` // 主键ID
	CreatedAt time.Time `gorm:"type:datetime(0);autoCreateTime;comment:创建时间" json:"createdAt" structs:"-"`
	UpdatedAt time.Time `gorm:"type:datetime(0);autoUpdateTime;comment:更新时间" json:"updatedAt" structs:"-"`
	// gorm会默认有一个DeletedAt字段，可以修改为IsDeleted,用0和1来表示是否删除
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt;default:0" json:"isDeleted" structs:"is_deleted"`
}
