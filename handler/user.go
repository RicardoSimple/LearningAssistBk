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
			LastLogin: func() string {
				if user.LastLogin.IsZero() {
					return "从未"
				} else {
					return user.LastLogin.Format("2006-01-02 15:04:05")
				}
			}(),
			UserType:   user.UserType,
			ClassStage: user.ClassStage,
			Name:       user.Name,
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
	var classId uint = 0
	if req.ClassNum != "" {
		// 校验班级码
		class, err := service.GetClassByClassNum(c, req.ClassNum)
		if err != nil {
			basic.RequestFailure(c, "<UNK>")
			return
		}
		if class == nil {
			basic.RequestFailure(c, "课堂邀请码错误")
			return
		}
		classId = uint(class.ID)
	}

	hashPassword, err := util.HashPassword(consts.DefaultPassword)

	user, err := service.CreateUser(c, &model.User{
		Username:    req.Username,
		Password:    hashPassword,
		Email:       req.Email,
		PhoneNumber: req.Phone,
		UserType:    req.UserType,
		Name:        req.Name,
	})

	if err != nil {
		basic.RequestFailure(c, "service error create user")
		return
	}
	if req.ClassNum != "" && classId != 0 {
		err = service.BindUserToClass(c, user.ID, classId)
		if err != nil {
			basic.RequestFailure(c, "绑定班级错误")
			return
		}
	}
	basic.Success(c, user)
}

// GetUserByType 按照类型获取用户
// @Summary 获取所有教师
// @Tags User
// @Success 200 {object} basic.Resp{data=[]UserResp}
// @Router /api/v1/user/byType [get]
func GetUserByType(c *gin.Context) {
	t := c.Query("type")
	if _, b := consts.UserTypeToIntMap[t]; !b {
		basic.RequestParamsFailure(c)
		return
	}
	teachers, err := service.GetUsersByType(c, t)
	if err != nil {
		basic.RequestFailure(c, "获取用户列表失败："+err.Error())
		return
	}

	result := make([]UserResp, 0, len(teachers))

	for _, user := range teachers {

		user.LastLogin.IsZero()

		result = append(result, UserResp{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Status:      user.Status,
			CreatedAt:   user.CreatedAt.Format("2006-01-02 15:04:05"),
			LastLogin: func() string {
				if user.LastLogin.IsZero() {
					return "从未"
				} else {
					return user.LastLogin.Format("2006-01-02 15:04:05")
				}
			}(),
			UserType:   user.UserType,
			ClassStage: user.ClassStage,
			Name:       user.Name,
		})
	}

	basic.Success(c, result)
}

// UpdateUserByAdmin 更新用户信息
// @Summary 管理员更新用户信息
// @Tags User
// @Param req body UserResp true "更新信息"
// @Success 200 {object} basic.Resp
// @Router /api/v1/user/update [post]
func UpdateUserByAdmin(c *gin.Context) {
	var req UserResp
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		basic.RequestParamsFailure(c)
		return
	}
	if _, b := consts.UserTypeToIntMap[req.UserType]; !b {
		basic.RequestParamsFailure(c)
		return
	}
	// 持久化更新
	err := service.UpdateUserByAdmin(c, req.ID, req.Email, req.PhoneNumber, req.UserType, req.Name)
	if err != nil {
		basic.RequestFailure(c, "更新失败："+err.Error())
		return
	}

	basic.Success(c, "更新成功")
}

type UserInfoResp struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

// GetUserInfoById
// @Summary 管理员获取用户信息
// @Tags User
// @Param id query string true
// @Success 200 {object} basic.Resp
// @Router /api/v1/user/info [get]
func GetUserInfoById(c *gin.Context) {
	id, err := util.GetQueryUint(c, "id")
	if err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	user, err := service.GetUserById(c, id)
	if err != nil {
		basic.RequestFailure(c, "获取用户信息错误")
		return
	}
	basic.Success(c, user)
}

// DeleteUserByAdmin 更新用户信息
// @Summary 管理员删除用户
// @Tags User
// @Param id query string true
// @Success 200 {object} basic.Resp
// @Router /api/v1/user/delete [post]
func DeleteUserByAdmin(c *gin.Context) {
	userId, err := util.GetQueryUint(c, "id")
	if err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	// todo
	service.DeleteUserByAdmin(c, userId)

	basic.Success(c, "删除成功")
}
