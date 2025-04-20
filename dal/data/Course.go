package data

import (
	"ar-app-api/dal/schema"
	"gorm.io/gorm"
)

// 创建科目
func CreateSubject(db *gorm.DB, name string) (*schema.Subject, error) {
	subject := &schema.Subject{Name: name}
	if err := db.Create(subject).Error; err != nil {
		return nil, err
	}
	return subject, nil
}

// 获取所有科目（用于下拉菜单/管理）
func GetAllSubjects(db *gorm.DB) ([]schema.Subject, error) {
	var subjects []schema.Subject
	err := db.Find(&subjects).Error
	return subjects, err
}

// 创建课程
func CreateCourse(db *gorm.DB, name string, subjectID, teacherID, classID uint, description string) (*schema.Course, error) {
	course := &schema.Course{
		Name:        name,
		SubjectID:   subjectID,
		TeacherID:   teacherID,
		ClassID:     classID,
		Description: description,
	}
	if err := db.Create(course).Error; err != nil {
		return nil, err
	}
	return course, nil
}

// 获取指定班级的课程（学生视角）
func GetCoursesByClassID(db *gorm.DB, classID uint) ([]schema.Course, error) {
	var courses []schema.Course
	err := db.Where("class_id = ?", classID).Find(&courses).Error
	return courses, err
}

// 获取指定教师的课程（教师视角）
func GetCoursesByTeacherID(db *gorm.DB, teacherID uint) ([]schema.Course, error) {
	var courses []schema.Course
	err := db.Where("teacher_id = ?", teacherID).Find(&courses).Error
	return courses, err
}

// 获取课程详情（含科目名）
func GetCourseWithSubject(db *gorm.DB, courseID uint) (*schema.Course, *schema.Subject, error) {
	var course schema.Course
	var subject schema.Subject
	if err := db.First(&course, courseID).Error; err != nil {
		return nil, nil, err
	}
	if err := db.First(&subject, course.SubjectID).Error; err != nil {
		return &course, nil, err
	}
	return &course, &subject, nil
}
