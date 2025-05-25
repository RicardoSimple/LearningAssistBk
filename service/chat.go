package service

import (
	"context"
	"encoding/json"
	"learning-assistant/dal"
	"learning-assistant/model"
	"learning-assistant/service/algo"
	"learning-assistant/util/log"
	"strings"
)

func CreateNewConversation(ctx context.Context, userId uint, title string) (uint, error) {
	id, err := dal.CreateConversation(ctx, userId, title)
	return id, err
}
func SaveMessage(ctx context.Context, conversationID uint, role, message string) error {
	err := dal.CreateChatMessage(ctx, conversationID, role, message)
	return err
}
func GetLastNMessagesForContext(ctx context.Context, conversationId uint, limit int) ([]*model.ChatMessage, error) {
	msgs, err := dal.GetLastNMessagesByConversationID(ctx, conversationId, limit)
	if err != nil {
		return nil, err
	}
	res := make([]*model.ChatMessage, 0, len(msgs))
	for _, msg := range msgs {
		res = append(res, msg.ToType())
	}
	return res, nil
}
func GetConversationsByUserId(ctx context.Context, userId uint) ([]*model.Conversation, error) {
	conersations, _, err := dal.GetConversationsByUserID(ctx, userId, 1, -1)
	if err != nil {
		return nil, err
	}
	res := make([]*model.Conversation, 0, len(conersations))
	for _, con := range conersations {
		res = append(res, con.ToType())
	}
	return res, nil
}

func DeleteConversationWithMessages(ctx context.Context, convID uint) error {
	return dal.DeleteConversationWithMessages(ctx, convID)
}

func SmartEvaluateAssign(ctx context.Context, assignmentId, submissionId uint) (model.AssignmentEvaluate, error) {
	result := model.AssignmentEvaluate{}
	assignDal, err := dal.GetAssignmentById(ctx, assignmentId)
	if err != nil {
		return result, err
	}
	assignment := assignDal.ToType()
	submissionDal, err := dal.GetSubmissionById(ctx, submissionId)
	if err != nil {
		return result, err
	}
	submission := submissionDal.ToType()

	// 拼接提示词
	msg := make([]algo.ChatMessage, 0, 2)
	msg = append(msg, algo.ChatMessage{
		Role:    algo.SystemRole,
		Content: algo.System_Evaluate_Score_Prompt,
	})
	msg = append(msg, algo.ChatMessage{
		Role:    algo.UserRole,
		Content: algo.BuildLLMEvaluationPrompt(assignment.Content, assignment.Title, submission.Content, submission.Title),
	})
	resp, err := algo.GetClient().Chat(msg, algo.ChatModel, true)
	log.Info("smart evaluate", resp)
	if err != nil {
		return result, err
	}
	// 转换json
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func SmartCourseDetail(ctx context.Context, name, description string, subjectIds []uint, prompt string) (string, error) {
	subjectsName, err := dal.GetSubjectNamesByIDs(ctx, subjectIds)
	if err != nil {
		return "", err
	}
	subjects := strings.Join(subjectsName, ", ")
	// 拼接提示词
	msg := make([]algo.ChatMessage, 0, 2)
	msg = append(msg, algo.ChatMessage{
		Role:    algo.SystemRole,
		Content: algo.BuildCourseDetailPrompt(name, description, subjects),
	})
	msg = append(msg, algo.ChatMessage{
		Role:    algo.UserRole,
		Content: prompt,
	})
	resp, err := algo.GetClient().Chat(msg, algo.ChatModel, false)
	log.Info("smart detail", resp)
	if err != nil {
		return "", err
	}
	return resp, nil
}

type TinyCourse struct {
	N string        `json:"n"`
	D string        `json:"d"`
	S []string      `json:"s"`
	C *model.Course `json:"-"`
}

type HotResp struct {
	RecommendedCourses []int `json:"recommendedCourses"`
}

func SmartLLMHotList(ctx context.Context, userId, n uint) ([]*model.Course, error) {

	allCourses, e := dal.GetAllCoursesWithSubjects(ctx)
	if e != nil {
		return nil, e
	}

	favoriteCourses, err := dal.GetFavoriteCoursesByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	courseMap := make(map[uint]TinyCourse)
	for _, course := range allCourses {
		f := course.ToType()
		subjects := make([]string, 0, len(f.Subjects))
		for _, subject := range f.Subjects {
			subjects = append(subjects, subject)
		}
		courseMap[f.ID] = TinyCourse{
			N: f.Name,
			D: f.Description,
			S: subjects,
			C: f,
		}
	}
	favoriteList := make([]uint, 0, len(favoriteCourses))
	for _, fav := range favoriteCourses {
		favoriteList = append(favoriteList, fav.ID)
	}

	courseMapJson, _ := json.Marshal(courseMap)
	listJson, _ := json.Marshal(favoriteList)
	// 拼接提示词

	msg := make([]algo.ChatMessage, 0, 2)
	msg = append(msg, algo.ChatMessage{
		Role:    algo.SystemRole,
		Content: algo.BuildHotCoursePrompt(n),
	})
	msg = append(msg, algo.ChatMessage{
		Role:    algo.UserRole,
		Content: algo.BuildHotCourseInput(string(courseMapJson), string(listJson)),
	})

	resp, err := algo.GetClient().Chat(msg, algo.ChatModel, true)
	log.Info("hot detail", resp)
	if err != nil {
		return nil, err
	}
	req := &HotResp{}

	json.Unmarshal([]byte(resp), &req)
	// findIdFromAllCourse

	result := make([]*model.Course, 0, n)

	for _, id := range req.RecommendedCourses {
		log.Info("hot detail", id)
		result = append(result, courseMap[uint(id)].C)
	}
	return result, nil
}
