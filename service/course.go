package service

import (
	"context"
	"fmt"
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
		courses[i] = &model.Course{
			ID:          s.ID,
			Name:        s.Name,
			Description: s.Description,
			Cover:       s.PageURL,
			Date:        s.CreatedAt,
			Subjects:    makeSubjects2Map(s.Subjects),
			Duration:    duration2Str(s.TotalTimeMinutes),
		}
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
		course.Description, course.Cover, subIds, duration)
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
	return &model.Course{
		ID:          c.ID,
		Name:        c.Name,
		Subjects:    makeSubjects2Map(c.Subjects),
		Cover:       c.PageURL,
		Description: c.Description,
		Duration:    duration2Str(c.TotalTimeMinutes),
		TeacherId:   c.TeacherID,
		ClassId:     c.ClassID,
		Date:        c.CreatedAt,
	}, nil
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
func duration2Str(totalTimeMinutes uint) string {
	return fmt.Sprintf("%02d小时%02d分钟", totalTimeMinutes/60, totalTimeMinutes%60)
}
