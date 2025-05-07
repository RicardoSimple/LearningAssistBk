package dal

import (
	"context"
	"learning-assistant/dal/schema"
	"learning-assistant/model"
	"time"
)

func CreateAssignment(ctx context.Context, title, content string, courseID, teacherID uint, classId uint, due time.Time) (*schema.Assignment, error) {
	a := &schema.Assignment{
		Title:     title,
		Content:   content,
		CourseID:  courseID,
		TeacherID: teacherID,
		DueDate:   due,
		ClassID:   classId,
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

func GetAssignmentsByClassId(ctx context.Context, classID uint) ([]schema.Assignment, error) {
	var list []schema.Assignment
	err := DB.Preload("Teacher").
		Preload("Course").
		Where("class_id = ?", classID).Find(&list).Error
	return list, err
}

func DeleteAssignment(ctx context.Context, assignmentID uint) error {
	return DB.Delete(&schema.Assignment{}, assignmentID).Error
}
func GetAssignmentsByClassIdPage(ctx context.Context, classID uint, page, pageSize int) ([]*model.Assignment, int64, error) {
	var (
		list  []schema.Assignment
		total int64
	)
	if err := DB.Model(&schema.Assignment{}).Where("class_id = ?", classID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := DB.Preload("Teacher").
		Preload("Course").
		Where("class_id = ?", classID).
		Order("due_date DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	assignments := make([]*model.Assignment, 0, len(list))
	for _, a := range list {
		assignments = append(assignments, a.ToType())
	}
	return assignments, total, nil
}
