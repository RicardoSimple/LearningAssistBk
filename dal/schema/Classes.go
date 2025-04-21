package schema

import "time"

// 班级表
type Class struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Grade     string    `gorm:"size:100;not null" json:"grade"`
	ClassNum  string    `gorm:"size:50;not null" json:"class_num"`
	CreatedAt time.Time `json:"created_at"`
}

// 教师-班级关系表（多对多）
type ClassTeacher struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	TeacherID uint `gorm:"not null" json:"teacher_id"` // 外键: User
	ClassID   uint `gorm:"not null" json:"class_id"`   // 外键: Class
}
