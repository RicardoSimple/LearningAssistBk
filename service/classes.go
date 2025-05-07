package service

import (
	"context"
	"learning-assistant/consts"
	"learning-assistant/dal"
	"learning-assistant/dal/schema"
	"learning-assistant/model"
	"learning-assistant/util"
	"strconv"
)

func CreateClass(ctx context.Context, name string, grade string) (*schema.Class, error) {
	classNum := util.GenerateInviteCode(consts.ClassNumLength)

	class, err := dal.CreateClass(ctx, name, grade, classNum)
	if err != nil {
		return nil, err
	}
	return class, nil
}
func GetClassListPage(ctx context.Context, page, pageSize int) ([]model.Class, int64, error) {
	classesPage, i, err := dal.GetClassesPage(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res := make([]model.Class, 0, len(classesPage))
	for _, v := range classesPage {
		res = append(res, *v.ToType())
	}
	return res, i, nil
}
func GetAllClassList(ctx context.Context) ([]*model.Class, error) {
	schemaClasses, err := dal.GetAllClasses(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*model.Class, 0, len(schemaClasses))
	for _, cls := range schemaClasses {
		res = append(res, cls.ToType())
	}
	return res, nil
}

func DeleteClassByID(ctx context.Context, id uint) error {
	return dal.DeleteClassByID(ctx, id)
}

func GetClassByTeacherID(ctx context.Context, id uint) ([]*model.Class, error) {
	clss, err := dal.GetClassesByTeacherID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := make([]*model.Class, 0, len(clss))
	for _, v := range clss {
		res = append(res, v.ToType())
	}
	return res, nil
}

func GetStudentsByClassId(ctx context.Context, classId string) ([]*model.User, error) {
	atoi, err := strconv.Atoi(classId)
	if err != nil {
		return nil, err
	}
	users, err := dal.GetUsersByClassID(ctx, uint(atoi))
	if err != nil {
		return nil, err
	}
	res := make([]*model.User, 0, len(users))
	for _, v := range users {
		res = append(res, v.ToType())
	}
	return res, nil
}
func BindTeacherToClass(ctx context.Context, teacherID, classID uint) error {
	return dal.AssignTeacherToClass(ctx, teacherID, classID)
}
