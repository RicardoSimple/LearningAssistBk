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

// schema/course.go
type Course struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:100;not null" json:"name"`
	TeacherID        uint      `json:"teacher_id"`
	ClassID          uint      `json:"class_id"`
	PageURL          string    `json:"page_url"`
	Description      string    `gorm:"type:text" json:"description"`
	TotalTimeMinutes uint      `json:"total_time"`
	CreatedAt        time.Time `json:"created_at"`

	Subjects []Subject `gorm:"many2many:course_subjects;" json:"subjects"`
}

// 中间表：多对多映射关系（可加上课程名索引或唯一约束）
type CourseSubject struct {
	CourseID  uint `gorm:"primaryKey"`
	SubjectID uint `gorm:"primaryKey"`
}
