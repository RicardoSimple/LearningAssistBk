package service

import (
	"context"
	"learning-assistant/dal"
	"learning-assistant/model"
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
