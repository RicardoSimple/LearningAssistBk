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
	ViewCount        uint      `gorm:"default:0" json:"view_count"` // 👈 新增：点击量字段
	CreatedAt        time.Time `json:"created_at"`
	CourseDetail     string    `gorm:"type:text" json:"course_detail"`
	FavoriteBy       []*User   `gorm:"many2many:user_course_favorites;"` // 被哪些用户收藏

	Subjects []Subject `gorm:"many2many:course_subjects;" json:"subjects"`
}

// 中间表：多对多映射关系（可加上课程名索引或唯一约束）
type CourseSubject struct {
	CourseID  uint `gorm:"primaryKey"`
	SubjectID uint `gorm:"primaryKey"`
}

func (course Course) ToType() *model.Course {

	favoriteBy := make([]*model.User, 0, len(course.FavoriteBy))
	for _, user := range course.FavoriteBy {
		favoriteBy = append(favoriteBy, user.ToType())
	}

	return &model.Course{
		ID:           course.ID,
		Name:         course.Name,
		Subjects:     makeSubjects2Map(course.Subjects),
		Cover:        course.PageURL,
		Description:  course.Description,
		Duration:     duration2Str(course.TotalTimeMinutes),
		TeacherId:    course.TeacherID,
		ClassId:      course.ClassID,
		ViewCount:    course.ViewCount,
		Date:         course.CreatedAt,
		CourseDetail: course.CourseDetail,
		FavoriteBy:   favoriteBy,
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
