package model

type TokenInfo struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ID           uint   `json:"id"`
	UserName     string `json:"user_name"`
}
