package data

import (
	"ar-app-api/dal/schema"
	"gorm.io/gorm"
)

func CreateClass(db *gorm.DB, name string, subjectID *uint) (*schema.Class, error) {
	class := &schema.Class{Name: name, SubjectID: subjectID}
	if err := db.Create(class).Error; err != nil {
		return nil, err
	}
	return class, nil
}
func AssignTeacherToClass(db *gorm.DB, teacherID uint, classID uint) error {
	relation := &schema.ClassTeacher{
		TeacherID: teacherID,
		ClassID:   classID,
	}
	return db.Create(relation).Error
}
func GetClassesByTeacherID(db *gorm.DB, teacherID uint) ([]schema.Class, error) {
	var classes []schema.Class
	err := db.Table("classes").
		Select("classes.*").
		Joins("JOIN class_teachers ON classes.id = class_teachers.class_id").
		Where("class_teachers.teacher_id = ?", teacherID).
		Scan(&classes).Error
	return classes, err
}
func GetTeachersByClassID(db *gorm.DB, classID uint) ([]uint, error) {
	var teacherIDs []uint
	err := db.Model(&schema.ClassTeacher{}).
		Where("class_id = ?", classID).
		Pluck("teacher_id", &teacherIDs).Error
	return teacherIDs, err
}
func RemoveTeacherFromClass(db *gorm.DB, teacherID uint, classID uint) error {
	return db.Where("teacher_id = ? AND class_id = ?", teacherID, classID).
		Delete(&schema.ClassTeacher{}).Error
}
func GetAllClasses(db *gorm.DB) ([]schema.Class, error) {
	var classes []schema.Class
	err := db.Find(&classes).Error
	return classes, err
}
