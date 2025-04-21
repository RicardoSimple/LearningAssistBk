package handler

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/handler/basic"
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
