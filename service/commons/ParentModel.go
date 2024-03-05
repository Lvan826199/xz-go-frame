/*
* @Author: 梦无矶小仔
* @Date: 2024/3/5 15:27
 */
package commons

type ParentModel[D any] struct {
	ID D `gorm:"primarykey;comment:主键ID" json:"id" structs:"-"` // 主键ID
}
