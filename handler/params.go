package handler

import (
	"ar-app-api/model"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	TokenInfo *model.TokenInfo `json:"token_info"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	// todo 其他参数
}

type CurrentUserResp struct {
	User *model.User `json:"user"`
}
