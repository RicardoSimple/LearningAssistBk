package schema

import (
	"learning-assistant/model"
	"time"
)

type Assignment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:255;not null" json:"title"`      // 作业标题
	Content   string    `gorm:"type:text;not null" json:"content"`   // 作业内容描述
	CourseID  uint      `gorm:"not null" json:"course_id"`           // 所属课程
	Course    Course    `gorm:"foreignKey:CourseID" json:"course"`   // 👈 预加载用
	TeacherID uint      `gorm:"not null" json:"teacher_id"`          // 发布教师
	Teacher   User      `gorm:"foreignKey:TeacherID" json:"teacher"` // 👈 添加这个字段以支持预加载
	DueDate   time.Time `json:"due_date"`                            // 截止日期
	CreatedAt time.Time `json:"created_at"`
	ClassID   uint      `gorm:"not null" json:"class_id"` // 班级id
}

func (assignment Assignment) ToType() *model.Assignment {
	return &model.Assignment{
		ID:        assignment.ID,
		Title:     assignment.Title,
		Content:   assignment.Content,
		CourseId:  assignment.CourseID,
		Course:    assignment.Course.ToType(),
		TeacherId: assignment.TeacherID,
		Teacher:   *assignment.Teacher.ToType(),
		DueDate:   assignment.DueDate.Format(time.DateTime),
		CreatedAt: assignment.CreatedAt.Format(time.DateTime),
		ClassID:   assignment.ClassID,
	}
}
