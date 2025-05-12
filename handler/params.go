package handler

import (
	"learning-assistant/model"
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

type SubjectResp struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type CourseResp struct {
	Id          int           `json:"id"`
	Cover       string        `json:"cover"`
	Name        string        `json:"name"`
	Subjects    []SubjectResp `json:"subjects"`
	Description string        `json:"description"`
	Duration    string        `json:"duration"`
	Date        string        `json:"date"`
	ViewCount   uint          `json:"view_count"`
}

type CourseDetailResp struct {
	Id           int           `json:"id"`
	Cover        string        `json:"cover"`
	Name         string        `json:"name"`
	Subjects     []SubjectResp `json:"subjects"`
	Description  string        `json:"description"`
	Duration     string        `json:"duration"`
	Date         string        `json:"date"`
	ViewCount    uint          `json:"view_count"`
	CourseDetail string        `json:"course_detail"`
	IsFavorite   bool          `json:"is_favorite"`
	FavoriteNum  uint          `json:"favorite_num"`
}

type CoursePageResp struct {
	Courses  []CourseResp `json:"courses"`
	Total    int          `json:"total"`
	PageSize int          `json:"page_size"`
	PageNum  int          `json:"page_num"`
}

type CreateSubjectReq struct {
	Name string `json:"name" binding:"required"`
}

type CreateCourseReq struct {
	Id           int    `json:"id"` // 用于更新接口
	Name         string `json:"name" binding:"required"`
	Cover        string `json:"cover"`
	Description  string `json:"description"`
	Duration     string `json:"duration"`
	Date         string `json:"date"` // 格式 "2006-01-02 15:04:05"
	TeacherID    uint   `json:"teacher_id"`
	ClassID      uint   `json:"class_id"`
	SubjectIDs   []uint `json:"subject_ids"` // 多个科目 ID
	CourseDetail string `json:"course_detail"`
}

type CreateClassReq struct {
	Name  string `json:"name" binding:"required"`
	Grade string `json:"grade"`
}
