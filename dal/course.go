package dal

import (
	"context"
	"gorm.io/gorm"
	"learning-assistant/dal/schema"
	"time"
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
	courseDetail string,
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
		CourseDetail:     courseDetail,
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
	err := DB.Preload("Subjects").
		Preload("FavoriteBy").
		First(&course, courseID).Error
	return &course, err
}

func DeleteCourseByID(ctx context.Context, id uint) error {
	return DB.Delete(&schema.Course{}, id).Error
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

// GetCourseByID 查询课程
func GetCourseByID(ctx context.Context, id uint) (*schema.Course, error) {
	var course schema.Course
	err := DB.WithContext(ctx).Preload("Subjects").First(&course, id).Error
	return &course, err
}

// UpdateCourse 更新课程及其关联关系
func UpdateCourseWithSubjects(ctx context.Context,
	id uint,
	name string,
	description, pageURL string,
	subjectIDs []uint,
	totalMinutes uint,
) error {
	// 查询id
	course, err := GetCourseByID(ctx, id)
	if err != nil {
		return err
	}

	// update
	// 获取所有 subject 实体
	var subjects []schema.Subject
	if err := DB.Where("id IN ?", subjectIDs).Find(&subjects).Error; err != nil {
		return err
	}

	course.Name = name
	course.Description = description
	course.PageURL = pageURL
	course.Subjects = subjects
	course.TotalTimeMinutes = totalMinutes

	return DB.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(course).Error
}

func IncrementCourseView(ctx context.Context, courseID uint) error {
	return DB.WithContext(ctx).
		Model(&schema.Course{}).
		Where("id = ?", courseID).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func AddFavorite(ctx context.Context, userID, courseID uint) error {
	return DB.WithContext(ctx).Create(&schema.UserCourseFavorite{
		UserID:    userID,
		CourseID:  courseID,
		CreatedAt: time.Now(),
	}).Error
}
func GetFavorite(ctx context.Context, userID, courseID uint) (*schema.UserCourseFavorite, error) {
	var course schema.UserCourseFavorite
	err := DB.WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).First(&course).Error
	return &course, err
}
func RemoveFavorite(ctx context.Context, userID, courseID uint) error {
	return DB.WithContext(ctx).Where("user_id = ? AND course_id = ?", userID, courseID).
		Delete(&schema.UserCourseFavorite{}).Error
}

// GetTopViewedCourses 获取浏览量前 N 的课程，若浏览量相同则按创建时间倒序
func GetTopViewedCourses(ctx context.Context, limit int) ([]schema.Course, error) {
	var courses []schema.Course
	err := DB.WithContext(ctx).
		Model(&schema.Course{}).
		Preload("Subjects").
		Order("view_count DESC").
		Order("created_at DESC").
		Limit(limit).
		Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
