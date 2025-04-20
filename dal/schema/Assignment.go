package schema

import (
	"time"
)

type Assignment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`    // 作业标题
	Content   string    `gorm:"type:text;not null" json:"content"` // 作业内容描述
	CourseID  uint      `gorm:"not null" json:"course_id"`         // 所属课程
	TeacherID uint      `gorm:"not null" json:"teacher_id"`        // 发布教师
	DueDate   time.Time `json:"due_date"`                          // 截止日期
	CreatedAt time.Time `json:"created_at"`
}
