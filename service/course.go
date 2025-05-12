package service

import (
	"context"
	"learning-assistant/dal"
	"learning-assistant/dal/schema"
	"learning-assistant/model"
)

func GetCourses(ctx context.Context, page, pageSize int) ([]*model.Course, int, error) {
	coursesData, i, err := dal.GetCoursesPage(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	courses := make([]*model.Course, len(coursesData))
	for i, s := range coursesData {
		courses[i] = s.ToType()
	}

	return courses, i, nil
}

func GetSubjects(ctx context.Context) (map[int]string, error) {
	subjects, err := dal.GetAllSubjects(ctx)
	if err != nil {
		return nil, err
	}
	return makeSubjects2Map(subjects), nil
}

func CreateCourse(ctx context.Context, course *model.Course, subIds []uint, duration uint) error {

	_, err := dal.CreateCourseWithSubjects(ctx, course.Name, course.TeacherId, course.ClassId,
		course.Description, course.Cover, subIds, duration, course.CourseDetail)
	return err
}
func CreateSubject(ctx context.Context, name string) (uint, error) {
	subject, err := dal.CreateSubject(ctx, name)
	if err != nil {
		return 0, err
	}
	return subject.ID, err
}

func GetCourseDetail(ctx context.Context, id uint) (*model.Course, error) {
	c, err := dal.GetCourseWithSubjects(ctx, id)
	if err != nil {
		return nil, err
	}
	return c.ToType(), nil
}

func DeleteCourse(ctx context.Context, id uint) error {
	return dal.DeleteCourseByID(ctx, id)
}
func makeSubjects2Map(subjects []schema.Subject) map[int]string {
	m := make(map[int]string, len(subjects))
	for _, s := range subjects {
		m[int(s.ID)] = s.Name
	}
	return m
}

// UpdateCourse 更新课程信息及其科目关联
func UpdateCourse(ctx context.Context, course *model.Course, subjectIDs []uint, duration uint) error {
	return dal.UpdateCourseWithSubjects(ctx, course.ID, course.Name, course.Description, course.Cover, subjectIDs, duration)
}
func IncrementCourseView(ctx context.Context, courseID uint) error {
	return dal.IncrementCourseView(ctx, courseID)
}

// AddCourseFavorite 添加收藏
func AddCourseFavorite(ctx context.Context, userID uint, courseID uint) error {

	favorite, _ := dal.GetFavorite(ctx, userID, courseID)
	if favorite != nil && favorite.CourseID == courseID && favorite.UserID == userID {
		return nil
	}

	return dal.AddFavorite(ctx, userID, courseID)
}

// RemoveCourseFavorite 取消收藏
func RemoveCourseFavorite(ctx context.Context, userID uint, courseID uint) error {
	return dal.RemoveFavorite(ctx, userID, courseID)
}

func GetHotNCourses(ctx context.Context, limit int) ([]*model.Course, error) {

	courses, err := dal.GetTopViewedCourses(ctx, limit)
	if err != nil {
		return nil, err
	}
	res := make([]*model.Course, 0, len(courses))
	for _, course := range courses {
		res = append(res, course.ToType())
	}
	return res, nil
}
