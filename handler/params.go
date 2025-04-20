package handler

import (
	"ar-app-api/model"
	"mime/multipart"
)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	TokenInfo *model.TokenInfo `json:"token_info"`
}

type RegisterReq struct {
	Username   string                `form:"username" binding:"required"`
	Phone      string                `form:"phone"`
	NickName   string                `form:"nickname"`
	Gender     string                `form:"gender"`
	Name       string                `form:"name"`
	Email      string                `form:"email"`
	Password   string                `form:"password" binding:"required"`
	ClassNum   string                `form:"class_num"`
	ClassStage string                `form:"class_stage"`
	UserType   string                `form:"user_type" binding:"required"`
	Number     string                `form:"number"` // 学号、工号等
	Avatar     *multipart.FileHeader `form:"avatar"` // 头像上传文件
}

type CurrentUserResp struct {
	User *model.User `json:"user"`
}
