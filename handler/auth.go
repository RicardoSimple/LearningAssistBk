package handler

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/consts"
	"learning-assistant/service"

	"learning-assistant/handler/aerrors"
	"learning-assistant/handler/basic"
	"learning-assistant/model"
	"learning-assistant/util"
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
	u, _ := service.GetUserByUserName(c, req.Username)
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
	if err := c.ShouldBind(req); err != nil {
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
	if !util.IsValidGrade(req.ClassStage) {
		basic.RequestFailureWithCode(c, "年级错误", aerrors.ParamsError)
	}
	name, _ := service.GetUserByUserName(c, req.Username)
	if name != nil {
		basic.RequestFailureWithCode(c, "用户名已存在", aerrors.ParamsError)
		return
	}

	if req.UserType == consts.UserTypeStudent {
		// todo 校验学生端参数
	} else if req.UserType == consts.UserTypeTeacher {
		// todo 校验教师端参数
	}

	// jiami1
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	user, err := service.CreateUser(c, &model.User{
		Username:    req.Username,
		Password:    hashPassword,
		PhoneNumber: req.Phone,
		Email:       req.Email,
		UserType:    req.UserType,
		ClassStage:  req.ClassStage,
		// ClassNum:    req.ClassNum,
		Gender: req.Gender,
		Name:   req.Name,
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
	userInfo, err := util.GetUserFromGinContext(c)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	u, err := service.GetUserByUserName(c, userInfo.UserName)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, &CurrentUserResp{
		User: u,
	})
}

// CheckToken 检查 Token 是否有效
// @Summary 检查Token
// @Tag Account.Auth(Token检查)
// @Param Authorization header string true "API TOKEN"
// @Success 200 {object} basic.Resp{data=string}
// @Failure 401 {object} basic.Resp
// @Router /account/auth/check [GET]
func CheckToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		basic.RequestFailureWithCode(c, "未提供Token", aerrors.NoAuth)
		return
	}

	_, err := util.GetUserFromGinContext(c)

	if err != nil {
		c.Writer.Header().Set("Authorization", "")
		basic.RequestFailureWithCode(c, "token不正确或已过期", aerrors.NoAuth)
		return
	}

	// 成功，返回 token（或你也可以返回 tokenInfo）
	c.Writer.Header().Set("Authorization", token)
	basic.Success(c, token)
}
