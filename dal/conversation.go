package dal

import (
	"context"
	"gorm.io/gorm"
	"learning-assistant/dal/schema"
)

// CreateConversation 创建对话
func CreateConversation(ctx context.Context, userId uint, title string) (uint, error) {
	c := &schema.Conversation{
		UserID: userId,
		Title:  title,
	}
	err := DB.WithContext(ctx).Create(c).Error
	return c.ID, err
}

// GetConversationByID 根据ID获取单个对话
func GetConversationByID(ctx context.Context, id uint) (*schema.Conversation, error) {
	var c schema.Conversation
	err := DB.WithContext(ctx).
		First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// CreateChatMessage 插入一条聊天记录
func CreateChatMessage(ctx context.Context, conversationId uint, role, content string) error {
	msg := &schema.ChatMessage{
		ConversationID: conversationId,
		Role:           role,
		Content:        content,
	}
	return DB.WithContext(ctx).Create(msg).Error
}

func GetConversationsByUserID(ctx context.Context, userID uint, page, pageSize int) ([]schema.Conversation, int64, error) {
	var list []schema.Conversation
	var total int64
	db := DB.WithContext(ctx).Model(&schema.Conversation{}).Where("user_id = ?", userID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("updated_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func GetMessagesByConversationID(ctx context.Context, conversationID uint, page, pageSize int) ([]schema.ChatMessage, int64, error) {
	var list []schema.ChatMessage
	var total int64
	db := DB.WithContext(ctx).Model(&schema.ChatMessage{}).Where("conversation_id = ?", conversationID)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at asc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
func GetLastNMessagesByConversationID(ctx context.Context, conversationID uint, limit int) ([]schema.ChatMessage, error) {
	var list []schema.ChatMessage
	err := DB.WithContext(ctx).
		Model(&schema.ChatMessage{}).
		Where("conversation_id = ?", conversationID).
		Order("created_at desc").
		Limit(limit).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	// 为了保证按时间顺序返回（升序），反转列表
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	return list, nil
}
func DeleteConversationWithMessages(ctx context.Context, convID uint) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除消息
		if err := tx.Where("conversation_id = ?", convID).Delete(&schema.ChatMessage{}).Error; err != nil {
			return err
		}
		// 删除对话
		if err := tx.Delete(&schema.Conversation{}, convID).Error; err != nil {
			return err
		}
		return nil
	})
}
