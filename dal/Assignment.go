package dal

import (
	"context"
	"learning-assistant/dal/schema"
	"time"
)

func CreateAssignment(ctx context.Context, title, content string, courseID, teacherID uint, due time.Time) (*schema.Assignment, error) {
	a := &schema.Assignment{
		Title:     title,
		Content:   content,
		CourseID:  courseID,
		TeacherID: teacherID,
		DueDate:   due,
	}
	if err := DB.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func GetAllAssignments(ctx context.Context) ([]schema.Assignment, error) {
	var list []schema.Assignment
	err := DB.Preload("Teacher").Preload("Course").Find(&list).Error
	return list, err
}

func GetAssignmentsByCourseID(ctx context.Context, courseID uint) ([]schema.Assignment, error) {
	var list []schema.Assignment
	err := DB.Preload("Teacher").
		Preload("Course").
		Where("course_id = ?", courseID).
		Order("due_date ASC").
		Find(&list).Error
	return list, err
}

func GetAssignmentsByTeacherID(ctx context.Context, teacherID uint) ([]schema.Assignment, error) {
	var list []schema.Assignment
	err := DB.Preload("Teacher").
		Preload("Course").
		Where("teacher_id = ?", teacherID).Find(&list).Error
	return list, err
}

func DeleteAssignment(ctx context.Context, assignmentID uint) error {
	return DB.Delete(&schema.Assignment{}, assignmentID).Error
}
