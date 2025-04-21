package handler

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/consts"
	"learning-assistant/handler/aerrors"
	"learning-assistant/handler/basic"
	"learning-assistant/model"
	"learning-assistant/service"
	"learning-assistant/util"
)

type UserResp struct {
	ID          uint   `json:"ID,omitempty"`
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"createdAt"`
	LastLogin   string `json:"lastLogin"`
	UserType    string `json:"userType,omitempty"`
	ClassStage  string `json:"class_stage,omitempty"`
	Name        string `json:"name,omitempty"`
}

type CreateUserResp struct {
	Username string `json:"username" binding:"required"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ClassNum string `json:"class_num"`
	UserType string `json:"user_type" binding:"required"`
}

// GetUserListHandler 分页获取用户列表
// @Summary 获取用户列表（分页）
// @Tags User
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Success 200 {object} basic.Resp{data=UserPageResp}
// @Router /api/v1/user/list [get]
func GetUserListHandler(c *gin.Context) {
	page, pageSize := util.GetPageParams(c)

	users, total, err := service.GetUserListPage(c, page, pageSize)
	if err != nil {
		basic.RequestFailure(c, "获取用户列表失败："+err.Error())
		return
	}
	result := make([]UserResp, 0, len(users))
	for _, user := range users {
		result = append(result, UserResp{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Status:      user.Status,
			CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
			LastLogin:   user.LastLogin.Format("2006-01-02 15:04:05"),
			UserType:    user.UserType,
			ClassStage:  user.ClassStage,
		})
	}

	basic.Success(c, gin.H{
		"list":     result,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateUserByAdmin 创建用户接口
// @Summary 用户注册
// @Tag Account.Auth(注册)
// @Param req body CreateUserResp
// @Success 200 {object} basic.Resp
// @Router /user/create [POST]
func CreateUserByAdmin(c *gin.Context) {
	//  参数校验
	req := &CreateUserResp{}
	if err := c.ShouldBindJSON(req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}

	name, _ := service.GetUserByUserName(c, req.Username)
	if name != nil {
		basic.RequestFailureWithCode(c, "用户名已存在", aerrors.ParamsError)
		return
	}

	user, err := service.CreateUser(c, &model.User{
		Username:    req.Username,
		Password:    consts.DefaultPassword,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		UserType:    req.UserType,
		Name:        req.Name,
	})
	if err != nil {
		basic.RequestFailure(c, "service error create user")
		return
	}
	basic.Success(c, user)
}
