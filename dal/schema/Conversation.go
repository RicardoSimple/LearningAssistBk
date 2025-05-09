package schema

import (
	"learning-assistant/model"
	"time"
)

type Conversation struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"` // 关联用户
	Title     string    `gorm:"size:128"` // 对话标题（可选）
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time
}

type ChatMessage struct {
	ID             uint   `gorm:"primaryKey"`
	ConversationID uint   `gorm:"index;not null"` // 关联对话
	Role           string `gorm:"size:16"`        // "user" 或 "assistant"
	Content        string `gorm:"type:text"`      // 消息内容（支持富文本或Markdown）
	//Tokens         int       `gorm:"default:0"`      // 消耗的token（可选）
	CreatedAt time.Time // 发送时间
}

func (cs *Conversation) ToType() *model.Conversation {
	return &model.Conversation{
		Id:        cs.ID,
		UserId:    cs.UserID,
		Title:     cs.Title,
		CreatedAt: cs.CreatedAt,
	}
}

func (cm *ChatMessage) ToType() *model.ChatMessage {
	return &model.ChatMessage{
		Id:             cm.ID,
		ConversationId: cm.ConversationID,
		Role:           cm.Role,
		Content:        cm.Content,
		CreatedAt:      cm.CreatedAt,
	}
}
