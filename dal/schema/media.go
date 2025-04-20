package schema

import (
	"gorm.io/gorm"
)

type Media struct {
	gorm.Model
	Content string `gorm:"type:text"`
	Type    string `gorm:"size:10;default:'text'"`
	Url     string `gorm:"size:255"`
}
