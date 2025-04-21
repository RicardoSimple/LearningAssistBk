package handler

import (
	"ar-app-api/handler/basic"
	"ar-app-api/service"
	"ar-app-api/util"
	"github.com/gin-gonic/gin"
)

// CreateClassHandler 新建班级
// @Summary 创建班级
// @Tags Class
// @Param req body CreateClassReq true "班级数据"
// @Success 200 {object} basic.Resp{data=schema.Class}
// @Router /api/v1/class/create [post]
func CreateClassHandler(c *gin.Context) {
	var req CreateClassReq
	if err := c.ShouldBindJSON(&req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	// todo 校验权限

	// 判断年级
	if !util.IsValidGrade(req.Grade) {
		basic.RequestParamsFailure(c)
	}
	clazz, err := service.CreateClass(c, req.Name, req.Grade)
	if err != nil {
		basic.RequestFailure(c, "创建班级失败："+err.Error())
		return
	}
	basic.Success(c, clazz)
}
