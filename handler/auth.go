package handler

import (
	"github.com/gin-gonic/gin"

	"ar-app-api/handler/aerrors"
	"ar-app-api/handler/basic"
	"ar-app-api/model"
	"ar-app-api/service/auth"
	"ar-app-api/util"
)

// Login 登录接口
// @Summary 用户登录
// @Tag Account.Auth(登录)
// @Param req body LoginReq true "上传图片文件"
// @Success 200 {object} basic.Resp{data=LoginResp}
// @Router /account/auth/login [POST]
func Login(c *gin.Context) {
	req := &LoginReq{}
	if err := c.BindJSON(req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	// 1. 查询用户
	u, _ := auth.GetUserByUserName(c, req.Username)
	if u == nil {
		basic.RequestFailureWithCode(c, "用户名错误", aerrors.RecordNotFind)
		return
	}
	// 2. 校验密码
	err := util.CheckPassword(u.Password, req.Password)
	if err != nil {
		basic.RequestFailureWithCode(c, "密码错误", aerrors.ParamsError)
		return
	}
	// 3. 组装token
	token, err := util.GenerateTokens(u)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, &LoginResp{
		TokenInfo: token,
	})
	return
}

// Register 注册接口
// @Summary 用户注册
// @Tag Account.Auth(注册)
// @Param req body RegisterReq true "上传图片文件"
// @Success 200 {object} basic.Resp{data=RegisterResp}
// @Router /account/auth/register [POST]
func Register(c *gin.Context) {
	req := &RegisterReq{}
	if err := c.BindJSON(req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	// 1. 普通参数校验
	if !util.IsValidEmail(req.Email) {
		basic.RequestFailureWithCode(c, "邮箱格式错误", aerrors.ParamsError)
		return
	}
	if !util.IsValidPhoneNumber(req.Phone) {
		basic.RequestFailureWithCode(c, "手机号格式错误", aerrors.ParamsError)
		return
	}
	name, _ := auth.GetUserByUserName(c, req.Username)
	if name != nil {
		basic.RequestFailureWithCode(c, "用户名已存在", aerrors.ParamsError)
		return
	}
	// jiami1
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	user, err := auth.CreateUser(c, &model.User{
		Username:    req.Username,
		Password:    hashPassword,
		PhoneNumber: req.Phone,
		Email:       req.Email,
	})
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	// 自动登录
	tokens, err := util.GenerateTokens(user)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, &LoginResp{
		TokenInfo: tokens,
	})
	return
}

// CurrentUser 获取当前用户信息
// @Summary 当前用户信息
// @Tag Account.Auth(当前用户信息)
// @Param Authorization header string true "API TOKEN"
// @Success 200 {object} basic.Resp{data=CurrentUserResp}
// @Router /account/auth/current [GET]
func CurrentUser(c *gin.Context) {
	userInfo, err := util.GetUserFromContext(c)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	u, err := auth.GetUserByUserName(c, userInfo.UserName)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, &CurrentUserResp{
		User: u,
	})
}
