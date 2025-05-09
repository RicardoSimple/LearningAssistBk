package service

import (
	"context"
	"encoding/json"
	"learning-assistant/dal"
	"learning-assistant/model"
	"learning-assistant/service/algo"
	"learning-assistant/util/log"
)

func CreateNewConversation(ctx context.Context, userId uint, title string) (uint, error) {
	id, err := dal.CreateConversation(ctx, userId, title)
	return id, err
}
func SaveMessage(ctx context.Context, conversationID uint, role, message string) error {
	err := dal.CreateChatMessage(ctx, conversationID, role, message)
	return err
}
func GetLastNMessagesForContext(ctx context.Context, conversationId uint, limit int) ([]*model.ChatMessage, error) {
	msgs, err := dal.GetLastNMessagesByConversationID(ctx, conversationId, limit)
	if err != nil {
		return nil, err
	}
	res := make([]*model.ChatMessage, 0, len(msgs))
	for _, msg := range msgs {
		res = append(res, msg.ToType())
	}
	return res, nil
}
func GetConversationsByUserId(ctx context.Context, userId uint) ([]*model.Conversation, error) {
	conersations, _, err := dal.GetConversationsByUserID(ctx, userId, 1, -1)
	if err != nil {
		return nil, err
	}
	res := make([]*model.Conversation, 0, len(conersations))
	for _, con := range conersations {
		res = append(res, con.ToType())
	}
	return res, nil
}

func DeleteConversationWithMessages(ctx context.Context, convID uint) error {
	return dal.DeleteConversationWithMessages(ctx, convID)
}

func SmartEvaluateAssign(ctx context.Context, assignmentId, submissionId uint) (model.AssignmentEvaluate, error) {
	result := model.AssignmentEvaluate{}
	assignDal, err := dal.GetAssignmentById(ctx, assignmentId)
	if err != nil {
		return result, err
	}
	assignment := assignDal.ToType()
	submissionDal, err := dal.GetSubmissionById(ctx, submissionId)
	if err != nil {
		return result, err
	}
	submission := submissionDal.ToType()

	// 拼接提示词
	msg := make([]algo.ChatMessage, 0, 2)
	msg = append(msg, algo.ChatMessage{
		Role:    algo.SystemRole,
		Content: algo.System_Evaluate_Score_Prompt,
	})
	msg = append(msg, algo.ChatMessage{
		Role:    algo.UserRole,
		Content: algo.BuildLLMEvaluationPrompt(assignment.Content, assignment.Title, submission.Content, submission.Title),
	})
	resp, err := algo.GetClient().Chat(msg, algo.ChatModel, true)
	log.Info("smart evaluate", resp)
	if err != nil {
		return result, err
	}
	// 转换json
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}
