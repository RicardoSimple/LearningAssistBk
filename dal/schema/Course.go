package schema

import (
	"time"
)

// 科目表（如：语文、数学）
type Subject struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null;unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// 课程表（可绑定班级和教师）
type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	SubjectID   uint      `json:"subject_id"` // 外键
	TeacherID   uint      `json:"teacher_id"` // 外键，关联用户表
	ClassID     uint      `json:"class_id"`   // 外键，关联班级
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
