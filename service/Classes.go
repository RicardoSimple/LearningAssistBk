package service

import (
	"context"
	"learning-assistant/consts"
	"learning-assistant/dal"
	"learning-assistant/dal/schema"
	"learning-assistant/util"
)

func CreateClass(ctx context.Context, name string, grade string) (*schema.Class, error) {
	classNum := util.GenerateInviteCode(consts.ClassNumLength)

	class, err := dal.CreateClass(ctx, name, grade, classNum)
	if err != nil {
		return nil, err
	}
	return class, nil
}
