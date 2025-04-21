package schema

import (
	"time"

	"gorm.io/gorm"

	"learning-assistant/model"
)

type Message struct {
	gorm.Model
	SenderID    uint      // 发送者的用户 ID
	ReceiverID  uint      // 接收者的用户 ID（如果是私聊，否则为 0）
	RoomID      uint      // 聊天室 ID（如果是群聊）
	Content     string    `gorm:"type:text"` // 消息内容
	MessageType string    `gorm:"size:20"`   // 消息类型 (例如：text, image, video)
	Timestamp   time.Time // 消息发送时间
	IsRead      bool      // 消息是否已读
	IsSend      bool      // 消息是否发送
}

// ChatGroup represents a chat group in the system. todo
type ChatGroup struct {
	gorm.Model
	Name        string `gorm:"size:255;not null"` // 群聊名称
	Description string `gorm:"size:1024"`         // 群聊描述
	CreatorID   uint   // 群聊创建者的用户 ID
	IsPrivate   bool   // 是否为私有群聊
	// 使用 GORM 的 many2many 特性建立与 User 的多对多关系
	Members []*User `gorm:"many2many:group_members;"`
}

func (c *ChatGroup) ToType() *model.ChatGroup {
	members := make([]*model.TinyUser, 0, len(c.Members))
	for _, user := range c.Members {
		members = append(members, &model.TinyUser{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			Gender:      user.Gender,
			DateOfBirth: user.DateOfBirth,
			Status:      user.Status,
			LastLogin:   user.LastLogin,
		})
	}
	return &model.ChatGroup{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatorID:   c.CreatorID,
		CreatedAt:   c.CreatedAt,
		IsPrivate:   c.IsPrivate,
		UpdatedAt:   c.UpdatedAt,
		Members:     members,
	}
}

func (m *Message) ToType() *model.Message {
	return &model.Message{
		ID:          m.ID,
		SenderID:    m.SenderID,
		ReceiverID:  m.ReceiverID,
		RoomID:      m.RoomID,
		Content:     m.Content,
		MessageType: m.MessageType,
		Timestamp:   m.Timestamp,
		IsRead:      m.IsRead,
		IsSend:      m.IsSend,
		UpdatedAt:   m.UpdatedAt,
		CreatedAt:   m.CreatedAt,
	}
}
