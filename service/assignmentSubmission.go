package service

import (
	"context"
	"learning-assistant/dal"
	"learning-assistant/model"
)

// SubmitAssignment 学生提交作业
func SubmitAssignment(ctx context.Context, sub *model.AssignmentSubmission) error {
	// 可以在这里添加防重复提交等校验
	return dal.SubmitAssignment(ctx, sub.AssignmentId, sub.StudentId, sub.Content)
}
