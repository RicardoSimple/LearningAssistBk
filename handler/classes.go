package handler

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/handler/basic"
	"learning-assistant/model"
	"learning-assistant/service"
	"learning-assistant/util"
)

type MyClassResp struct {
	Classes []*model.Class `json:"classes"`
}

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

// DeleteClassHandler 删除班级
// @Summary 删除班级
// @Tags Class
// @Param id query int true "班级ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/class/delete [POST]
func DeleteClassHandler(c *gin.Context) {
	id, err := util.GetQueryUint(c, "id")
	if err != nil || id == 0 {
		basic.RequestParamsFailure(c)
		return
	}

	// 可选：校验权限，例如仅管理员可删

	err = service.DeleteClassByID(c, id)
	if err != nil {
		basic.RequestFailure(c, "删除班级失败："+err.Error())
		return
	}
	basic.Success(c, "删除成功")
}

// GetClassListHandler 分页获取班级列表
// @Summary 获取班级列表（分页）
// @Tags Class
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Success 200 {object} basic.Resp{data=ClassPageResp}
// @Router /api/v1/class/list [get]
func GetClassListHandler(c *gin.Context) {
	page, pageSize := util.GetPageParams(c)

	classes, total, err := service.GetClassListPage(c, page, pageSize)
	if err != nil {
		basic.RequestFailure(c, "获取班级列表失败："+err.Error())
		return
	}
	basic.Success(c, gin.H{
		"list":     classes,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetAllClassHandler 获取所有班级（不分页）
// @Summary 获取所有班级
// @Tags Class
// @Success 200 {object} basic.Resp{data=[]schema.Class}
// @Router /api/v1/class/all [get]
func GetAllClassHandler(c *gin.Context) {
	classes, err := service.GetAllClassList(c)
	if err != nil {
		basic.RequestFailure(c, "获取所有班级失败："+err.Error())
		return
	}
	basic.Success(c, classes)
}

// GetMyClassHandler 教师端获取自己管理的班级信息
// @Summary 教师获取所属班级
// @Tags Class
// @Success 200 {object} basic.Resp{data=schema.Class}
// @Router /api/v1/class/my [get]
func GetMyClassHandler(c *gin.Context) {
	user, err := util.GetUserFromGinContext(c)
	if err != nil {
		basic.RequestFailure(c, "用户信息获取失败")
		return
	}
	classModels, err := service.GetClassByTeacherID(c, user.ID)
	if err != nil {
		basic.RequestFailure(c, "查询班级失败")
		return
	}
	resp := MyClassResp{Classes: classModels}

	basic.Success(c, resp)
}

type MyStudentsResp struct {
	Students []*model.User `json:"students"`
}

// GetMyClassStudentsHandler 教师端获取当前班级学生
// @Summary 教师获取所属班级学生
// @Tags Class
// @Params classId query true
// @Success 200 {object} basic.Resp{data=[]schema.User}
// @Router /api/v1/class/my/students [get]
func GetMyClassStudentsHandler(c *gin.Context) {
	classId := c.Query("classId")
	if classId == "" {
		basic.RequestParamsFailure(c)
		return
	}
	students, err := service.GetStudentsByClassId(c, classId)
	if err != nil {
		basic.RequestFailure(c, "查询学生失败")
		return
	}
	resp := &MyStudentsResp{Students: students}
	basic.Success(c, resp)
}

// BindTeacherToClassHandler 教师绑定到某个班级
// @Summary 教师绑定班级
// @Tags Class
// @Param classId body uint true "班级ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/class/bind [post]
func BindTeacherToClassHandler(c *gin.Context) {
	var req struct {
		ClassID   uint `json:"classId" binding:"required"`
		TeacherId uint `json:"teacherId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.ClassID == 0 {
		basic.RequestParamsFailure(c)
		return
	}

	// 调用 service 层绑定逻辑
	if err := service.BindTeacherToClass(c, req.TeacherId, req.ClassID); err != nil {
		basic.RequestFailure(c, "绑定失败："+err.Error())
		return
	}

	basic.Success(c, "绑定成功")
}

// BindUserToClassHandler 用户绑定班级
// @Summary 用户绑定班级（如学生）
// @Tags Class
// @Param req body struct{ UserID uint `json:"userId" binding:"required"`; ClassID uint `json:"classId" binding:"required"` } true "绑定数据"
// @Success 200 {object} basic.Resp
// @Router /api/v1/class/user/bind [post]
func BindUserToClassHandler(c *gin.Context) {
	var req struct {
		UserID  uint `json:"userId" binding:"required"`
		ClassID uint `json:"classId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.UserID == 0 || req.ClassID == 0 {
		basic.RequestParamsFailure(c)
		return
	}

	err := service.BindUserToClass(c, req.UserID, req.ClassID)
	if err != nil {
		basic.RequestFailure(c, "绑定失败："+err.Error())
		return
	}

	basic.Success(c, "绑定成功")
}
