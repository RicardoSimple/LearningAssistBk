package handler

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/handler/basic"
	"learning-assistant/model"
	"learning-assistant/service"
	"learning-assistant/util"
	"strconv"
)

type SubmitAssignmentReq struct {
	AssignmentID string `json:"assignment_id" binding:"required"`
	Content      string `json:"content" binding:"required"`
	Title        string `json:"title" binding:"required"`
}

// SubmitAssignmentHandler 学生提交作业
// @Summary 作业管理
// @Tags Assignment
// @Success 200 {object} basic.Resp{}
// @Router /api/v1/assignment/submit [post]
func SubmitAssignmentHandler(c *gin.Context) {
	user, err := util.GetUserFromGinContext(c)
	if err != nil || user.UserType != "student" {
		basic.AuthFailure(c)
		return
	}
	var req SubmitAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	as_id, _ := strconv.Atoi(req.AssignmentID)
	sub := &model.AssignmentSubmission{
		AssignmentId: uint(as_id),
		StudentId:    user.ID,
		Content:      req.Content,
		Title:        req.Title,
	}
	err = service.SubmitAssignment(c, sub)
	if err != nil {
		basic.RequestFailure(c, "提交失败："+err.Error())
		return
	}
	basic.Success(c, "提交成功")
}

type AssignmentSubmissionPageResp struct {
	Submissions []*model.AssignmentSubmission `json:"list"`
	Total       int64                         `json:"total"`
	Page        int                           `json:"page"`
	PageSize    int                           `json:"page_size"`
}

// GetAssignmentSubmissionsHandler 查询作业提交记录（分页）
// @Summary 获取提交记录
// @Tags Assignment
// @Param assignment_id query int false "作业ID，若为空则查询全部"
// @Param page query int false "页码"
// @Param page_size query int false "每页条数"
// @Success 200 {object} basic.Resp{data=AssignmentSubmissionPageResp}
// @Router /api/v1/assignment/submissions [get]
func GetAssignmentSubmissionsHandler(c *gin.Context) {
	assignmentID, _ := strconv.Atoi(c.Query("assignment_id"))
	page, pageSize := util.GetPageParams(c)

	submissions, total, err := service.GetAssignmentSubmissionsPage(c, uint(assignmentID), page, pageSize)
	if err != nil {
		basic.RequestFailure(c, "获取提交记录失败："+err.Error())
		return
	}

	resp := AssignmentSubmissionPageResp{
		Submissions: submissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
	}
	basic.Success(c, resp)
}

type EvaluateAssignmentReq struct {
	SubmissionID uint   `json:"submission_id" binding:"required"`
	Score        int    `json:"score" binding:"required"`
	Feedback     string `json:"feedback"`
}

// EvaluateAssignmentSubmissionHandler 教师评价作业提交
// @Summary 教师评价作业
// @Tags Assignment
// @Param req body EvaluateAssignmentReq true "评价信息"
// @Success 200 {object} basic.Resp
// @Router /api/v1/assignment/evaluate [post]
func EvaluateAssignmentSubmissionHandler(c *gin.Context) {
	var req EvaluateAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil || req.SubmissionID == 0 {
		basic.RequestParamsFailure(c)
		return
	}

	// 获取当前登录用户（建议后续加权限判断）
	user, err := util.GetUserFromGinContext(c)
	if err != nil {
		basic.AuthFailure(c)
		return
	}

	// 可选：验证是否为教师
	if user.UserType != "teacher" {
		basic.RequestFailure(c, "仅教师可以评价作业")
		return
	}

	err = service.EvaluateAssignmentSubmission(c, req.SubmissionID, req.Score, req.Feedback)
	if err != nil {
		basic.RequestFailure(c, "评价失败："+err.Error())
		return
	}

	basic.Success(c, "评价成功")
}
