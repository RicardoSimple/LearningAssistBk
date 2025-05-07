package service

import (
	"context"
	"learning-assistant/dal"
	"learning-assistant/dal/schema"
	"learning-assistant/model"
	"time"
)

func CreateAssignment(ctx context.Context, title, content string, courseID, teacherID, classId uint, due time.Time) (*model.Assignment, error) {
	assignment, err := dal.CreateAssignment(ctx, title, content, courseID, teacherID, classId, due)
	if err != nil {
		return nil, err
	}
	return assignment.ToType(), nil
}

func GetAssignmentList(ctx context.Context, userType string) ([]*model.Assignment, error) {

	assignments, err := dal.GetAllAssignments(ctx)

	if err != nil {
		return nil, err
	}
	res := make([]*model.Assignment, 0, len(assignments))
	for _, assignment := range assignments {
		res = append(res, assignment.ToType())
	}
	return res, nil
}

func GetAssignmentsByCourseID(ctx context.Context, courseID uint) ([]*model.Assignment, error) {
	assigns, err := dal.GetAssignmentsByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	res := make([]*model.Assignment, 0, len(assigns))

	for _, assign := range assigns {
		res = append(res, assign.ToType())
	}
	return res, nil
}

func GetAssignmentsByTeacherID(ctx context.Context, teacherID uint) ([]*model.Assignment, error) {
	assigns, err := dal.GetAssignmentsByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	res := make([]*model.Assignment, 0, len(assigns))

	for _, assign := range assigns {
		res = append(res, assign.ToType())
	}
	return res, nil
}

func GetAssignmentByID(ctx context.Context, id uint) (*model.Assignment, error) {
	var a schema.Assignment
	err := dal.DB.First(&a, id).Error
	return a.ToType(), err
}

func DeleteAssignment(ctx context.Context, id uint) error {
	return dal.DeleteAssignment(ctx, id)
}
func GetAssignmentsByClassIdPage(ctx context.Context, classID uint, page, pageSize int) ([]*model.Assignment, int64, error) {
	return dal.GetAssignmentsByClassIdPage(ctx, classID, page, pageSize)
}
