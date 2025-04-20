package model

import (
	"time"
)

type Message struct {
	ID          uint      `json:"ID,omitempty"`
	SenderID    uint      `json:"senderID,omitempty"`    // 发送者的用户 ID
	ReceiverID  uint      `json:"receiverID,omitempty"`  // 接收者的用户 ID（如果是私聊，否则为 0）
	RoomID      uint      `json:"roomID,omitempty"`      // 聊天室 ID（如果是群聊）
	Content     string    `json:"content,omitempty"`     // 消息内容
	MessageType string    `json:"messageType,omitempty"` // 消息类型 (例如：text, image, video)
	Timestamp   time.Time `json:"timestamp"`             // 消息发送时间
	IsRead      bool      `json:"isRead,omitempty"`      // 消息是否已读
	IsSend      bool      `json:"isSend,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ChatGroup struct {
	ID          uint        `json:"ID,omitempty"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	CreatorID   uint        `json:"creatorID,omitempty"` // 群聊创建者的用户 ID
	CreatedAt   time.Time   `json:"createdAt"`           // 群聊创建时间
	IsPrivate   bool        `json:"isPrivate,omitempty"` // 是否为私有群聊
	UpdatedAt   time.Time   `json:"updatedAt"`
	Members     []*TinyUser `json:"users,omitempty"`
}

// TinyUser 基础用户信息 展示用
type TinyUser struct {
	ID          uint      `json:"ID,omitempty"`
	Username    string    `json:"username,omitempty"`
	Email       string    `json:"email,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Status      string    `json:"status,omitempty"`
	LastLogin   time.Time `json:"lastLogin"`
}
