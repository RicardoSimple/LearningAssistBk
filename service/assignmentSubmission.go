package service

import (
	"context"
	"learning-assistant/dal"
	"learning-assistant/model"
	"time"
)

// SubmitAssignment 学生提交作业
func SubmitAssignment(ctx context.Context, sub *model.AssignmentSubmission) error {
	// 可以在这里添加防重复提交等校验
	return dal.SubmitAssignment(ctx, sub.AssignmentId, sub.StudentId, sub.Content, sub.Title)
}

// GetAssignmentSubmissionsPage 获取提交记录（分页）
func GetAssignmentSubmissionsPage(ctx context.Context, assignmentID uint, page, pageSize int) ([]*model.AssignmentSubmission, int64, error) {
	list, total, err := dal.GetAssignmentSubmissionsPage(ctx, assignmentID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 model 层结构体
	result := make([]*model.AssignmentSubmission, 0, len(list))
	for _, s := range list {
		result = append(result, s.ToType())
	}
	return result, total, nil
}
func EvaluateAssignmentSubmission(ctx context.Context, submissionID uint, score float64, feedback string) error {
	now := time.Now()
	return dal.UpdateAssignmentSubmissionEvaluation(ctx, submissionID, score, feedback, &now)
}
