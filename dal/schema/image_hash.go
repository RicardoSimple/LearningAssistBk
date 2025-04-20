package schema

import (
	"gorm.io/gorm"
)

// ImageHash 是图片哈希表的模型
type ImageHash struct {
	gorm.Model
	URL     string `gorm:"size:255"`          // 图片URL，唯一索引
	Hash    string `gorm:"size:100;not null"` // 图片哈希
	Desc    string `gorm:"size:255"`          // 描述
	MediaId uint   `gorm:"unique"`
}
