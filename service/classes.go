package service

import (
	"context"
	"learning-assistant/consts"
	"learning-assistant/dal"
	"learning-assistant/dal/schema"
	"learning-assistant/model"
	"learning-assistant/util"
	"time"
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
		res = append(res, model.Class{
			ID:         int(v.ID),
			Name:       v.Name,
			InviteCode: v.ClassNum,
			Grade:      v.Grade,
			CreatedAt:  v.CreatedAt.Format(time.DateTime),
		})
	}
	return res, i, nil
}
