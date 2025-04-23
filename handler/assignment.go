package handler

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/handler/basic"
	"learning-assistant/model"
	"learning-assistant/service"
	"learning-assistant/util"
	"strconv"
	"time"
)

type CreateAssignmentReq struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	CourseID  uint   `json:"course_id" binding:"required"`
	TeacherID uint   `json:"teacher_id" binding:"required"`
	DueDate   string `json:"due_date" binding:"required"` // 格式：2025-05-01 23:59:00
}
type AssignmentPageResp struct {
	Courses  []*model.Assignment `json:"list"`
	Total    int                 `json:"total"`
	PageSize int                 `json:"page_size"`
	PageNum  int                 `json:"page_num"`
}

// CreateAssignmentHandler
// @Summary 创建作业
// @Tags Course
// @Param req body CreateAssignmentReq true
// @Success 200 {object} basic.Resp
// @Router /api/v1/assignment/create [post]
func CreateAssignmentHandler(c *gin.Context) {
	var req CreateAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}

	dueTime, err := time.Parse("2006-01-02 15:04:05", req.DueDate)
	if err != nil {
		basic.RequestFailure(c, "截止时间格式错误，应为：2006-01-02 15:04:05")
		return
	}

	assignment, err := service.CreateAssignment(c, req.Title, req.Content, req.CourseID, req.TeacherID, dueTime)
	if err != nil {
		basic.RequestFailure(c, "创建作业失败："+err.Error())
		return
	}
	basic.Success(c, assignment)
}

// GetAssignments
// @Summary 获取课程下的作业
// @Tags Course
// @Param
// @Success 200 {object} basic.Resp
// GET /api/v1/assignment/list
func GetAssignments(c *gin.Context) {

	userInfo, err := util.GetUserFromGinContext(c)
	if err != nil {
		basic.AuthFailure(c)
		return
	}
	userType := userInfo.UserType
	// todo 根据userType 自动调整获取内容
	// todo pageable
	assignments, err := service.GetAssignmentList(c, userType)
	if err != nil {
		basic.RequestFailure(c, "获取失败："+err.Error())
		return
	}
	resp := AssignmentPageResp{
		Courses: assignments,
		Total:   len(assignments),
	}
	basic.Success(c, resp)
}

// GetAssignmentsByCourseHandler
// @Summary 获取课程下的作业
// @Tags Course
// @Param course_id query
// @Success 200 {object} basic.Resp
// GET /api/v1/assignment/listByCourse?course_id=1
func GetAssignmentsByCourseHandler(c *gin.Context) {
	courseID, _ := strconv.Atoi(c.Query("course_id"))
	if courseID <= 0 {
		basic.RequestParamsFailure(c)
		return
	}
	assignments, err := service.GetAssignmentsByCourseID(c, uint(courseID))
	if err != nil {
		basic.RequestFailure(c, "获取失败："+err.Error())
		return
	}
	basic.Success(c, assignments)
}

// GetAssignmentsByTeacherHandler
// @Summary 获取教师创建的作业
// @Tags Course
// @Param teacher_id query
// @Success 200 {object} basic.Resp
// GET /api/v1/assignment/listByTeacher?teacher_id=1
func GetAssignmentsByTeacherHandler(c *gin.Context) {
	tid, _ := strconv.Atoi(c.Query("teacher_id"))
	if tid <= 0 {
		basic.RequestParamsFailure(c)
		return
	}
	assignments, err := service.GetAssignmentsByTeacherID(c, uint(tid))
	if err != nil {
		basic.RequestFailure(c, "获取失败："+err.Error())
		return
	}
	basic.Success(c, assignments)
}

// GetAssignmentDetailHandler
// @Summary 获取作业详情
// @Tags Course
// @Param id query
// @Success 200 {object} basic.Resp
// GET /api/v1/assignment/detail?id=1
func GetAssignmentDetailHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		basic.RequestParamsFailure(c)
		return
	}
	assignment, err := service.GetAssignmentByID(c, uint(id))
	if err != nil {
		basic.RequestFailure(c, "获取详情失败："+err.Error())
		return
	}
	basic.Success(c, assignment)
}

// DeleteAssignmentHandler
// @Summary 删除
// @Tags Course
// @Param id query
// @Success 200 {object} basic.Resp
// POST /api/v1/assignment/delete?id=1
func DeleteAssignmentHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	if id <= 0 {
		basic.RequestParamsFailure(c)
		return
	}
	err := service.DeleteAssignment(c, uint(id))
	if err != nil {
		basic.RequestFailure(c, "删除失败："+err.Error())
		return
	}
	basic.Success(c, "已删除")
}
