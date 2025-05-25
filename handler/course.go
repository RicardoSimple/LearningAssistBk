package handler

import (
	"github.com/RicardoSimple/hao-tool/lists"
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
			ViewCount:   c.ViewCount,
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
	courses, total, err := service.GetCourses(c, page, pageSize)
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
			ViewCount:   c.ViewCount,
		})
	}
	resp := CoursePageResp{
		Courses:  res,
		Total:    total,
		PageNum:  page,
		PageSize: pageSize,
	}
	basic.Success(c, resp)
}

// GetCourseById 获取课程
// @Summary 获取课程列表
// @Tags Course
// @Param id query uint
// @Success 200 CoursePageResp basic.Resp{data=CoursePageResp}
// @Router /api/v1/course/byId/get [get]
func GetCourseById(c *gin.Context) {
	id, err := util.GetQueryUint(c, "id")
	if err != nil {
		basic.RequestFailure(c, "<UNK>id<UNK>"+err.Error())
		return
	}
	courseschema, err := service.GetCourseById(c, id)
	if err != nil {
		basic.RequestFailure(c, "<UNK>id<UNK>"+err.Error())
		return
	}
	subjectsIds := make([]uint, 0, len(courseschema.Subjects))
	for _, s := range courseschema.Subjects {
		subjectsIds = append(subjectsIds, uint(s.ID))
	}
	basic.Success(c, CreateCourseReq{
		Id:           int(courseschema.ID),
		Name:         courseschema.Name,
		Cover:        courseschema.PageURL,
		Description:  courseschema.Description,
		Duration:     strconv.Itoa(int(courseschema.TotalTimeMinutes)),
		Date:         courseschema.ToType().Date.Format("2006-01-02 15:04:05"),
		TeacherID:    courseschema.TeacherID,
		ClassID:      courseschema.ClassID,
		SubjectIDs:   subjectsIds,
		CourseDetail: courseschema.CourseDetail,
	})
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
		Name:         req.Name,
		Subjects:     nil,
		Cover:        req.Cover,
		Description:  req.Description,
		Duration:     "",
		Date:         parsedTime,
		CourseDetail: req.CourseDetail,
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

	// 忽略未登录状态
	user, _ := util.GetUserFromGinContext(c)
	isFavorite := false
	if user != nil && user.ID > 0 {
		isFavorite = lists.InListWithFn(course.FavoriteBy, &model.User{ID: user.ID}, func(a, b *model.User) bool {
			return a.ID == b.ID
		})
	}

	res := CourseDetailResp{
		Id:           int(course.ID),
		Cover:        course.Cover,
		Name:         course.Name,
		Subjects:     mapToSubjectResp(course.Subjects),
		Description:  course.Description,
		Duration:     course.Duration,
		Date:         course.Date.Format("2006-01-02 15:04:05"),
		ViewCount:    course.ViewCount,
		CourseDetail: course.CourseDetail,
		IsFavorite:   isFavorite,
		FavoriteNum:  uint(len(course.FavoriteBy)),
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

// UpdateCourseHandler 更新课程
// @Summary 更新课程
// @Tags Course
// @Param req body CreateCourseReq true "更新后的课程数据（需携带ID）"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/update [post]
func UpdateCourseHandler(c *gin.Context) {
	var req CreateCourseReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Id == 0 {
		basic.RequestParamsFailure(c)
		return
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", req.Date)
	if err != nil {
		basic.RequestFailure(c, "时间格式错误")
		return
	}
	duration, err := strconv.Atoi(req.Duration)
	if err != nil {
		basic.RequestFailure(c, "时长格式错误")
		return
	}

	course := &model.Course{
		ID:           uint(req.Id),
		Name:         req.Name,
		Cover:        req.Cover,
		Description:  req.Description,
		Duration:     req.Duration,
		Date:         parsedTime,
		CourseDetail: req.CourseDetail,
	}

	err = service.UpdateCourse(c, course, req.SubjectIDs, uint(duration))
	if err != nil {
		basic.RequestFailure(c, "更新课程失败："+err.Error())
		return
	}
	basic.Success(c, "更新成功")
}

// IncrementCourseViewHandler 增加课程点击量
// @Summary 课程点击统计
// @Tags Course
// @Param id query int true "课程ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/view/increase [post]
func IncrementCourseViewHandler(c *gin.Context) {

	id, err := util.GetQueryUint(c, "id")
	if err != nil {
		basic.RequestFailure(c, "课程ID格式错误")
		return
	}

	if err := service.IncrementCourseView(c, id); err != nil {
		basic.RequestFailure(c, "增加点击量失败："+err.Error())
		return
	}
	basic.Success(c, "点击量+1")
}

// FavoriteCourseHandler 用户收藏课程
// @Summary 收藏课程
// @Tags Course
// @Param course_id query int true "课程ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/favorite [post]
func FavoriteCourseHandler(c *gin.Context) {
	user, err := util.GetUserFromGinContext(c)
	if err != nil {
		basic.AuthFailure(c)
		return
	}
	courseID, err := util.GetQueryUint(c, "course_id")
	if err != nil || courseID == 0 {
		basic.RequestParamsFailure(c)
		return
	}
	err = service.AddCourseFavorite(c, user.ID, courseID)
	if err != nil {
		basic.RequestFailure(c, "收藏失败："+err.Error())
		return
	}
	basic.Success(c, "收藏成功")
}

// UnfavoriteCourseHandler 用户取消收藏课程
// @Summary 取消收藏课程
// @Tags Course
// @Param course_id query int true "课程ID"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/unfavorite [post]
func UnfavoriteCourseHandler(c *gin.Context) {
	user, err := util.GetUserFromGinContext(c)
	if err != nil {
		basic.AuthFailure(c)
		return
	}
	courseID, err := util.GetQueryUint(c, "course_id")
	if err != nil || courseID == 0 {
		basic.RequestParamsFailure(c)
		return
	}
	err = service.RemoveCourseFavorite(c, user.ID, courseID)
	if err != nil {
		basic.RequestFailure(c, "取消收藏失败："+err.Error())
		return
	}
	basic.Success(c, "取消收藏成功")
}

// FindHotNCourse 获取点击量前n的课程
// @Summary hot课程
// @Tags Course
// @Param n query int true "n"
// @Success 200 {object} basic.Resp
// @Router /api/v1/course/hot [get]
func FindHotNCourse(c *gin.Context) {
	n, err := util.GetQueryUint(c, "n")
	if err != nil {
		basic.RequestParamsFailure(c)
		return
	}
	if n <= 0 {
		n = 10
	}
	courses, err := service.GetHotNCourses(c, int(n))
	if err != nil {
		basic.RequestFailure(c, "获取失败")
		return
	}
	basic.Success(c, courses)
}
