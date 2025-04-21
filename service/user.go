package service

import (
	"context"
	"errors"
	"learning-assistant/consts"
	"learning-assistant/dal"
	"time"

	"gorm.io/gorm"

	"learning-assistant/model"
	"learning-assistant/util"
)

func CreateUser(ctx context.Context, user *model.User) (*model.User, error) {

	// todo 自动加入班级

	suser, err := dal.CreateUser(ctx, util.ToUserSchema(user))
	if err != nil {
		return nil, err
	}
	return suser.ToType(), nil
}

func UpdateUser(ctx context.Context, user model.User) {

}

func GetUserByUserName(ctx context.Context, userName string) (*model.User, error) {
	if userName == "" {
		return nil, errors.New("userName is empty")
	}
	user, err := dal.GetUserByUsername(ctx, userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user.ToType(), nil
}

func UpdateLoginStatus(ctx context.Context, id uint) {
	user, err := dal.GetUserByID(ctx, id)
	if err != nil {
		return
	}
	user.LastLogin = time.Now()
	user.Status = consts.OnLine
	dal.UpdateUser(ctx, user)
}
