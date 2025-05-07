package handler

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/handler/basic"
	"learning-assistant/model"
	"learning-assistant/service"
	"learning-assistant/util"
)

type SubmitAssignmentReq struct {
	AssignmentID uint   `json:"assignment_id" binding:"required"`
	Content      string `json:"content" binding:"required"`
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
	sub := &model.AssignmentSubmission{
		AssignmentId: req.AssignmentID,
		StudentId:    user.ID,
		Content:      req.Content,
	}
	err = service.SubmitAssignment(c, sub)
	if err != nil {
		basic.RequestFailure(c, "提交失败："+err.Error())
		return
	}
	basic.Success(c, "提交成功")
}
