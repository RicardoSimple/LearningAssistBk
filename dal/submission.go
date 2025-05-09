package dal

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"learning-assistant/dal/schema"
	"time"
)

func SubmitAssignment(ctx context.Context, assignmentId, studentId uint, content, title string) error {
	sub := &schema.AssignmentSubmission{
		AssignmentID: assignmentId,
		StudentID:    studentId,
		Content:      content,
		Title:        title,
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

func GetSubmissionById(ctx context.Context, submissionId uint) (*schema.AssignmentSubmission, error) {
	var sub schema.AssignmentSubmission
	err := DB.WithContext(ctx).
		Where("id = ?", submissionId).
		First(&sub).Error
	return &sub, err
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
func GetSubmissionByAssignmentAndUser(ctx context.Context, assignmentID, userID uint) (*schema.AssignmentSubmission, error) {
	var s schema.AssignmentSubmission
	err := DB.
		Where("assignment_id = ? AND student_id = ?", assignmentID, userID).
		Order("submitted_at DESC").
		First(&s).Error

	if err != nil {
		// 如果没有记录，不返回错误，只返回 nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}
func GetAssignmentSubmissionsPage(ctx context.Context, assignmentID uint, page, pageSize int) ([]schema.AssignmentSubmission, int64, error) {
	var list []schema.AssignmentSubmission
	var total int64
	db := DB.WithContext(ctx).Model(&schema.AssignmentSubmission{})

	if assignmentID > 0 {
		db = db.Where("assignment_id = ?", assignmentID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Order("submitted_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&list).Error
	return list, total, err
}
func UpdateAssignmentSubmissionEvaluation(ctx context.Context, id uint, score int, feedback string, reviewedAt *time.Time) error {
	return DB.WithContext(ctx).
		Model(&schema.AssignmentSubmission{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"score":       score,
			"feedback":    feedback,
			"reviewed_at": reviewedAt,
		}).Error
}
