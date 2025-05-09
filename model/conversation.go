package model

import "time"

type Conversation struct {
	Id        uint      `json:"id"`
	UserId    uint      `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatMessage struct {
	Id             uint      `json:"id"`
	ConversationId uint      `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}
