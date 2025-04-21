package handler

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/handler/basic"
	"learning-assistant/model"
	"learning-assistant/service"
	"time"
)

// GetCoursesHandler 获取所有课程
// @Summary 获取课程列表
// @Tags Course
// @Success 200 []CourseResp basic.Resp{data=[]CourseResp}
// @Router /api/v1/course/courses [get]
func GetCoursesHandler(c *gin.Context) {
	courses, err := service.GetAllCourses(c)
	if err != nil {
		basic.RequestFailure(c, "获取课程失败："+err.Error())
		return
	}
	res := make([]CourseResp, 0, len(courses))
	for _, c := range courses {
		res = append(res, CourseResp{
			Id:          int(c.ID),
			Cover:       c.Cover,
			Name:        c.Name,
			Subjects:    mapToSubjectResp(c.Subjects),
			Description: c.Description,
			Duration:    c.Duration,
			Date:        c.Date.Format("2006-01-02 15:04:05"),
		})
	}
	basic.Success(c, courses)
}

// CreateSubjectHandler 创建科目
// @Summary 创建科目
// @Tags Course
// @Param req body CreateSubjectReq true "科目名称"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/subject/create [post]
func CreateSubjectHandler(c *gin.Context) {
	var req CreateSubjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	sub, err := service.CreateSubject(c, req.Name)
	if err != nil {
		basic.RequestFailure(c, "创建科目失败："+err.Error())
		return
	}
	basic.Success(c, sub)
}

// CreateCourseHandler 创建课程
// @Summary 创建课程
// @Tags Course
// @Param req body CreateCourseReq true "课程数据"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/create [post]
func CreateCourseHandler(c *gin.Context) {
	var req CreateCourseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		basic.RequestParamsFailure(c)
		return
	}

	// 时间解析
	parsedTime, err := time.Parse("2006-01-02 15:04:05", req.Date)
	if err != nil {
		basic.RequestFailure(c, "时间格式错误")
		return
	}
	// todo 校验权限
	err = service.CreateCourse(c, &model.Course{
		Name:        req.Name,
		Subjects:    nil,
		Cover:       req.Cover,
		Description: req.Description,
		Duration:    "",
		Date:        parsedTime,
	}, req.SubjectIDs, req.Duration)
	if err != nil {
		basic.RequestFailure(c, "创建课程失败："+err.Error())
		return
	}
	basic.Success(c, nil)
}

func mapToSubjectResp(subs map[int]string) []SubjectResp {
	res := make([]SubjectResp, 0, len(subs))
	for k, v := range subs {
		res = append(res, SubjectResp{
			Id:   k,
			Name: v,
		})
	}
	return res
}
