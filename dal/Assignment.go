package dal

import (
	"ar-app-api/dal/schema"
	"gorm.io/gorm"
	"time"
)

func CreateAssignment(db *gorm.DB, title, content string, courseID, teacherID uint, due time.Time) (*schema.Assignment, error) {
	a := &schema.Assignment{
		Title:     title,
		Content:   content,
		CourseID:  courseID,
		TeacherID: teacherID,
		DueDate:   due,
	}
	if err := db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func GetAssignmentsByCourseID(db *gorm.DB, courseID uint) ([]schema.Assignment, error) {
	var list []schema.Assignment
	err := db.Where("course_id = ?", courseID).
		Order("due_date ASC").
		Find(&list).Error
	return list, err
}

func GetAssignmentsByTeacherID(db *gorm.DB, teacherID uint) ([]schema.Assignment, error) {
	var list []schema.Assignment
	err := db.Where("teacher_id = ?", teacherID).Find(&list).Error
	return list, err
}

func DeleteAssignment(db *gorm.DB, assignmentID uint) error {
	return db.Delete(&schema.Assignment{}, assignmentID).Error
}
