package service

import (
	"ar-app-api/consts"
	"ar-app-api/dal"
	"ar-app-api/dal/schema"
	"ar-app-api/util"
	"context"
)

func CreateClass(ctx context.Context, name string, grade string) (*schema.Class, error) {
	classNum := util.GenerateInviteCode(consts.ClassNumLength)

	class, err := dal.CreateClass(ctx, name, grade, classNum)
	if err != nil {
		return nil, err
	}
	return class, nil
}
