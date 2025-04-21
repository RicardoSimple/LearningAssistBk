package dal

import (
	"context"
	"learning-assistant/dal/schema"
)

// 创建科目
func CreateSubject(ctx context.Context, name string) (*schema.Subject, error) {
	subject := &schema.Subject{Name: name}
	if err := DB.Create(subject).Error; err != nil {
		return nil, err
	}
	return subject, nil
}

// 获取所有科目（用于下拉菜单/管理）
func GetAllSubjects(ctx context.Context) ([]schema.Subject, error) {
	var subjects []schema.Subject
	err := DB.Find(&subjects).Error
	return subjects, err
}

func CreateCourseWithSubjects(
	ctx context.Context,
	name string,
	teacherID, classID uint,
	description, pageURL string,
	subjectIDs []uint,
	totalMinutes uint,
) (*schema.Course, error) {
	// 获取所有 subject 实体
	var subjects []schema.Subject
	if err := DB.Where("id IN ?", subjectIDs).Find(&subjects).Error; err != nil {
		return nil, err
	}

	course := &schema.Course{
		Name:             name,
		TeacherID:        teacherID,
		ClassID:          classID,
		Description:      description,
		PageURL:          pageURL,
		Subjects:         subjects,
		TotalTimeMinutes: totalMinutes,
	}
	if err := DB.Create(course).Error; err != nil {
		return nil, err
	}
	return course, nil
}

// 获取指定班级的课程（学生视角）
func GetCoursesByClassID(ctx context.Context, classID uint) ([]schema.Course, error) {
	var courses []schema.Course
	err := DB.Where("class_id = ?", classID).Find(&courses).Error
	return courses, err
}

// 获取指定教师的课程（教师视角）
func GetCoursesByTeacherID(ctx context.Context, teacherID uint) ([]schema.Course, error) {
	var courses []schema.Course
	err := DB.Where("teacher_id = ?", teacherID).Find(&courses).Error
	return courses, err
}

func GetCourseWithSubjects(ctx context.Context, courseID uint) (*schema.Course, error) {
	var course schema.Course
	err := DB.Preload("Subjects").First(&course, courseID).Error
	return &course, err
}
func GetCoursesBySubjectID(ctx context.Context, subjectID uint) ([]schema.Course, error) {
	var courses []schema.Course
	err := DB.Joins("JOIN course_subjects ON course_subjects.course_id = courses.id").
		Where("course_subjects.subject_id = ?", subjectID).
		Find(&courses).Error
	return courses, err
}

func GetAllCoursesWithSubjects(ctx context.Context) ([]schema.Course, error) {
	var courses []schema.Course
	err := DB.Preload("Subjects").Find(&courses).Error
	return courses, err
}
func GetCoursesPage(ctx context.Context, page, pageSize int) ([]schema.Course, int, error) {
	if page < 0 || pageSize < 0 {
		subjects, err := GetAllCoursesWithSubjects(ctx)
		return subjects, len(subjects), err
	}
	var (
		courses []schema.Course
		total   int64
	)

	// 查询总数
	if err := DB.Model(&schema.Course{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询分页数据并预加载 Subjects
	offset := (page - 1) * pageSize
	err := DB.Preload("Subjects").
		Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&courses).Error
	if err != nil {
		return nil, 0, err
	}

	return courses, int(total), nil
}
