package service

import (
	"ar-app-api/dal"
	"ar-app-api/dal/schema"
	"ar-app-api/model"
	"context"
	"fmt"
)

func GetAllCourses(ctx context.Context) ([]*model.Course, error) {
	coursesData, err := dal.GetAllCoursesWithSubjects(ctx)
	if err != nil {
		return nil, err
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

	return courses, nil
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
func makeSubjects2Map(subjects []schema.Subject) map[int]string {
	m := make(map[int]string, len(subjects))
	for i, s := range subjects {
		m[i] = s.Name
	}
	return m
}
func duration2Str(totalTimeMinutes uint) string {
	return fmt.Sprintf("%02d:%02d", totalTimeMinutes/60, totalTimeMinutes%60)
}
