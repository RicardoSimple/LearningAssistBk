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

// GetCoursesHandler 获取所有课程
// @Summary 获取课程列表
// @Tags Course
// @Success 200 []CourseResp basic.Resp{data=[]CourseResp}
// @Router /api/v1/course/courses [get]
func GetCoursesHandler(c *gin.Context) {
	courses, _, err := service.GetCourses(c, -1, -1)
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
	basic.Success(c, res)
}

// GetCoursesByPage 获取所有课程，分页查询
// @Summary 获取课程列表
// @Tags Course
// @Param page query string
// @Param pageSize query string
// @Success 200 CoursePageResp basic.Resp{data=CoursePageResp}
// @Router /api/v1/course/get [get]
func GetCoursesByPage(c *gin.Context) {
	page, pageSize := util.GetPageParams(c)
	courses, _, err := service.GetCourses(c, page, pageSize)
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
	resp := CoursePageResp{
		Courses:  res,
		Total:    len(courses),
		PageNum:  page,
		PageSize: pageSize,
	}
	basic.Success(c, resp)
}

// GetSubjects 获取所有科目
// @Summary 获取课程列表
// @Tags Course
// @Success 200 []SubjectResp basic.Resp{data=[]SubjectResp}
// @Router /api/v1/course/subject/getAll [get]
func GetSubjects(c *gin.Context) {
	subsMap, err := service.GetSubjects(c)
	if err != nil {
		basic.RequestFailure(c, "error"+err.Error())
	}
	basic.Success(c, mapToSubjectResp(subsMap))
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
	duration, err := strconv.Atoi(req.Duration)
	err = service.CreateCourse(c, &model.Course{
		Name:        req.Name,
		Subjects:    nil,
		Cover:       req.Cover,
		Description: req.Description,
		Duration:    "",
		Date:        parsedTime,
	}, req.SubjectIDs, uint(duration))
	if err != nil {
		basic.RequestFailure(c, "创建课程失败："+err.Error())
		return
	}
	basic.Success(c, nil)
}

// GetCourseDetailHandler 获取课程详情
// @Summary 获取课程详情
// @Tags Course
// @Param id query int true "课程ID"
// @Success 200 {object} basic.Resp{data=CourseResp}
// @Router /api/v1/course/detail [get]
func GetCourseDetailHandler(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		basic.RequestParamsFailure(c)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		basic.RequestFailure(c, "课程ID格式错误")
		return
	}

	course, err := service.GetCourseDetail(c, uint(id))
	if err != nil {
		basic.RequestFailure(c, "获取课程详情失败："+err.Error())
		return
	}

	res := CourseDetailResp{
		Id:          int(course.ID),
		Cover:       course.Cover,
		Name:        course.Name,
		Subjects:    mapToSubjectResp(course.Subjects),
		Description: course.Description,
		Duration:    course.Duration,
		Date:        course.Date.Format("2006-01-02 15:04:05"),
	}
	basic.Success(c, res)
}

// DeleteCourseHandler 删除课程
// @Summary 删除课程
// @Tags Course
// @Param id query int true "课程ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/delete [post]
func DeleteCourseHandler(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		basic.RequestParamsFailure(c)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		basic.RequestFailure(c, "课程ID格式错误")
		return
	}

	err = service.DeleteCourse(c, uint(id))
	if err != nil {
		basic.RequestFailure(c, "删除课程失败："+err.Error())
		return
	}
	basic.Success(c, "删除成功")
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
