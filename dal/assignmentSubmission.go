package dal

import (
	"context"
	"learning-assistant/dal/schema"
	"time"
)

func SubmitAssignment(ctx context.Context, assignmentId, studentId uint, content string) error {
	sub := &schema.AssignmentSubmission{
		AssignmentID: assignmentId,
		StudentID:    studentId,
		Content:      content,
		SubmittedAt:  time.Now(),
	}
	return DB.WithContext(ctx).Create(sub).Error
}

func GetMySubmissions(ctx context.Context, studentID uint) ([]schema.AssignmentSubmission, error) {
	var list []schema.AssignmentSubmission
	err := DB.WithContext(ctx).
		Where("student_id = ?", studentID).
		Preload("Assignment").
		Find(&list).Error
	return list, err
}

func GetSubmissionsByAssignment(ctx context.Context, assignmentID uint) ([]schema.AssignmentSubmission, error) {
	var list []schema.AssignmentSubmission
	err := DB.WithContext(ctx).
		Where("assignment_id = ?", assignmentID).
		Preload("Assignment").
		Preload("Student").
		Find(&list).Error
	return list, err
}
