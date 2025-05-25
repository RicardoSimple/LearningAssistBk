package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"learning-assistant/handler/basic"
	"learning-assistant/service"
	"learning-assistant/service/algo"
	"learning-assistant/util"
	"strconv"
)

type ChatReq struct {
	ConversationId uint   `json:"conversation_id"`
	Question       string `json:"message"`
}

// ChatAssistant 智能对话（支持流式返回）
// @Summary 智能学习助手对话
// @Tags Chat
// @Accept json
// @Produce text/event-stream
// @Param body  ChatReq true
// @Success 200 {string} string "流式文本输出"
// @Router /api/v1/chat/assistant [post]
func ChatAssistant(c *gin.Context) {
	// 设置为 SSE 流式响应
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 获取用户信息（假设有中间件设置）
	user, err := util.GetUserFromGinContext(c)
	if err != nil {
		basic.AuthFailure(c)
		return
	}
	req := &ChatReq{}
	if err := c.ShouldBindJSON(&req); err != nil || req.Question == "" {
		basic.RequestParamsFailure(c)
		return
	}

	// 若无对话 ID，则新建
	conversationID := req.ConversationId
	if conversationID == 0 {
		conversationID, err = service.CreateNewConversation(c, user.ID, req.Question)
		if err != nil {
			basic.RequestFailure(c, "创建对话失败")
			return
		}
	}

	// 保存用户提问
	err = service.SaveMessage(c, conversationID, algo.UserRole, req.Question)
	if err != nil {
		c.String(500, "保存提问失败")
		return
	}
	// 获取上下文消息
	history, err := service.GetLastNMessagesForContext(c, conversationID, 6)
	if err != nil {
		c.String(500, "获取上下文失败")
		return
	}

	// 构建为 LLM 所需的格式
	var messages []algo.ChatMessage
	messages = append(messages, algo.ChatMessage{
		Role:    algo.SystemRole,
		Content: algo.System_Assistant_Prompt,
	})
	for _, m := range history {
		messages = append(messages, algo.ChatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	messages = append(messages, algo.ChatMessage{
		Role:    algo.UserRole,
		Content: req.Question,
	})

	// 调用流式响应
	var answer string
	client := algo.GetClient()
	err = client.ChatStream(messages, algo.ChatModel, func(token string) {
		answer += token
		fmt.Fprintf(c.Writer, "data: %s\n\n", token)
		c.Writer.Flush()
	})
	if err != nil {
		c.String(500, "生成失败："+err.Error())
		return
	}

	// 保存 assistant 回复
	_ = service.SaveMessage(c, conversationID, algo.AssistantRole, answer)
}

// GetConversations 获取对话列表
// @Summary 智能学习助手对话
// @Tags Chat
// @Accept json
// @Produce text/event-stream
// @Success 200 {string} string "流式文本输出"
// @Router /api/v1/chat/conversations [get]
func GetConversations(c *gin.Context) {
	user, err := util.GetUserFromGinContext(c)
	if err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	conversations, err := service.GetConversationsByUserId(c, user.ID)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, conversations)
}

// GetMessages 获取对话记录
// @Summary 智能学习助手对话
// @Tags Chat
// @Accept json
// @Param query conversation_id string
// @Produce text/event-stream
// @Success 200 {string} string "流式文本输出"
// @Router /api/v1/chat/messages [get]
func GetMessages(c *gin.Context) {
	idStr := c.Query("conversation_id")
	csId, err := strconv.Atoi(idStr)
	if err != nil {
		basic.RequestParamsFailure(c)
		return
	}

	msgs, err := service.GetLastNMessagesForContext(c, uint(csId), -1)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, msgs)
}

// DeleteConversationHandler 删除指定对话及其所有消息
// @Summary 删除对话
// @Tags Chat
// @Param id query int true "对话 ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/chat/conversation/delete [post]
func DeleteConversationHandler(c *gin.Context) {
	id, err := util.GetQueryUint(c, "id")
	if err != nil || id == 0 {
		basic.RequestParamsFailure(c)
		return
	}

	if err := service.DeleteConversationWithMessages(c, id); err != nil {
		basic.RequestFailure(c, "删除对话失败："+err.Error())
		return
	}
	basic.Success(c, "删除成功")
}

type SmartEvaluateReq struct {
	AssignmentId uint `json:"assignmentId"`
	SubmissionId uint `json:"submissionId"`
}

// SmartEvaluateAssignment 智能评估作业
// @Summary 作业评估
// @Tags Chat
// @Param SmartEvaluateReq body int true "作业 ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/assignment/algo/evaluate [post]
func SmartEvaluateAssignment(c *gin.Context) {
	req := &SmartEvaluateReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}

	// todo 根据两个id确定缓存 减少开销
	result, err := service.SmartEvaluateAssign(c, req.AssignmentId, req.SubmissionId)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, result)
}

// SmartCourseDetail 智能获取课程详情及路线
// @Summary 内容生成
// @Tags Chat
// @Param GenerateCourseReq body int true "作业 ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/algo/detail [post]
func SmartCourseDetail(c *gin.Context) {
	req := &GenerateCourseReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	// todo

	result, err := service.SmartCourseDetail(c, req.Name, req.Description, req.SubjectIDs, req.CoursePrompt)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, result)
}

// SmartNewCourses 个性化推荐课程
// @Summary 内容生成
// @Tags Chat
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/algo/hot [post]
func SmartNewCourses(c *gin.Context) {
	n, err := util.GetQueryUint(c, "n")
	if err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	user, _ := util.GetUserFromGinContext(c)
	var userId uint = 0
	if user != nil {
		// 已登录
		userId = user.ID
	}
	list, err := service.SmartLLMHotList(c, userId, n)
	if err != nil {
		basic.RequestFailure(c, err.Error())
		return
	}
	basic.Success(c, list)
}

// SmartAssignmentDetail 智能生成作业内容
// @Summary 内容生成
// @Tags Chat
// @Param CreateCourseReq body int true "作业 ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/algo/detail [post]
func SmartAssignmentDetail(c *gin.Context) {

}
