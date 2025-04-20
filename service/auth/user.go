package auth

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"ar-app-api/dal/data"
	"ar-app-api/model"
	"ar-app-api/util"
)

func CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	suser, err := data.CreateUser(ctx, util.ToUserSchema(user))
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
	user, err := data.GetUserByUsername(ctx, userName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user.ToType(), nil
}

func UpdateLoginStatus(ctx context.Context, id uint) {
	user, err := data.GetUserByID(ctx, id)
	if err != nil {
		return
	}
	user.LastLogin = time.Now()
	user.Status = util.OnLine
	data.UpdateUser(ctx, user)
}
