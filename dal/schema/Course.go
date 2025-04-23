package schema

import (
	"fmt"
	"learning-assistant/model"
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

func (course Course) ToType() *model.Course {
	return &model.Course{
		ID:          course.ID,
		Name:        course.Name,
		Subjects:    makeSubjects2Map(course.Subjects),
		Cover:       course.PageURL,
		Description: course.Description,
		Duration:    duration2Str(course.TotalTimeMinutes),
		TeacherId:   course.TeacherID,
		ClassId:     course.ClassID,
		Date:        course.CreatedAt,
	}
}

func makeSubjects2Map(subjects []Subject) map[int]string {
	m := make(map[int]string, len(subjects))
	for _, s := range subjects {
		m[int(s.ID)] = s.Name
	}
	return m
}
func duration2Str(totalTimeMinutes uint) string {
	return fmt.Sprintf("%02d小时%02d分钟", totalTimeMinutes/60, totalTimeMinutes%60)
}
