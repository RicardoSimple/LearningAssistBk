package dal

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"learning-assistant/dal/schema"
)

func CreateClass(ctx context.Context, name string, grade string, classNum string) (*schema.Class, error) {
	class := &schema.Class{Name: name, Grade: grade, ClassNum: classNum}
	if err := DB.Create(class).Error; err != nil {
		return nil, err
	}
	return class, nil
}
func AssignTeacherToClass(ctx context.Context, teacherID uint, classID uint) error {
	relation := &schema.ClassTeacher{
		TeacherID: teacherID,
		ClassID:   classID,
	}
	return DB.Create(relation).Error
}
func GetClassesByTeacherID(ctx context.Context, teacherID uint) ([]schema.Class, error) {
	var classes []schema.Class
	err := DB.Table("classes").
		Select("classes.*").
		Joins("JOIN class_teachers ON classes.id = class_teachers.class_id").
		Where("class_teachers.teacher_id = ?", teacherID).
		Scan(&classes).Error
	return classes, err
}
func GetClassByClassNum(ctx context.Context, classNum string) (*schema.Class, error) {
	var class schema.Class
	err := DB.Where("class_num = ?", classNum).First(&class).Error
	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &class, nil
}
func GetTeachersByClassID(ctx context.Context, classID uint) ([]uint, error) {
	var teacherIDs []uint
	err := DB.Model(&schema.ClassTeacher{}).
		Where("class_id = ?", classID).
		Pluck("teacher_id", &teacherIDs).Error
	return teacherIDs, err
}
func RemoveTeacherFromClass(ctx context.Context, teacherID uint, classID uint) error {
	return DB.Where("teacher_id = ? AND class_id = ?", teacherID, classID).
		Delete(&schema.ClassTeacher{}).Error
}
func GetAllClasses(ctx context.Context) ([]schema.Class, error) {
	var classes []schema.Class
	err := DB.Find(&classes).Error
	return classes, err
}
func GetClassesPage(ctx context.Context, page, pageSize int) ([]schema.Class, int64, error) {
	var (
		classes []schema.Class
		total   int64
	)

	if err := DB.Model(&schema.Class{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := DB.
		Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&classes).Error
	if err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}
func DeleteClassByID(ctx context.Context, id uint) error {
	return DB.Delete(&schema.Class{}, id).Error
}
