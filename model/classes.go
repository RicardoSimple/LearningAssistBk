package model

// 班级表
type Class struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	InviteCode string `json:"invite_code"`
	Grade      string `json:"grade"`
	CreatedAt  string `json:"created_at"`
}
