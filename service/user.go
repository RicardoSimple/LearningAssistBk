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
	suser, err := dal.CreateUser(ctx, util.ToUserSchema(user))
	if err != nil {
		return nil, err
	}
	return suser.ToType(), nil
}

func UpdateUserByAdmin(ctx context.Context, id uint, email, phoneNumber, userType, name string) error {
	// 获取原始用户数据
	oldUser, err2 := dal.GetUserByID(ctx, id)
	if err2 != nil {
		return err2
	}

	// 更新字段

	oldUser.Name = name
	oldUser.Email = email
	oldUser.PhoneNumber = phoneNumber
	oldUser.UserType = consts.UserTypeToIntMap[userType]

	return dal.UpdateUser(ctx, oldUser)
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
func GetUserListPage(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	usersPage, i, err := dal.GetUsersPage(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	res := make([]model.User, 0, len(usersPage))
	for _, user := range usersPage {
		res = append(res, *user.ToType())
	}
	return res, i, nil
}
func GetUserById(ctx context.Context, id uint) (*model.User, error) {
	user, err := dal.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToType(), nil
}
func GetUsersByType(ctx context.Context, userType string) ([]*model.User, error) {
	users, err := dal.GetUsersByType(ctx, userType)
	if err != nil {
		return nil, err
	}
	res := make([]*model.User, 0, len(users))
	for _, user := range users {
		res = append(res, user.ToType())
	}
	return res, nil
}
func BindUserToClass(ctx context.Context, userID, classID uint) error {
	user, err := dal.GetUserByID(ctx, userID)
	s := consts.UserTypeToStringMap[user.UserType]
	// 教师用户增加关联表
	if s == consts.UserTypeTeacher {
		err := BindTeacherToClass(ctx, userID, classID)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	user.ClassId = classID
	user.AddClassTime = time.Now()
	return dal.UpdateUser(ctx, user)
}
func DeleteUserByAdmin(ctx context.Context, id uint) error {
	user, err := dal.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	// 删除关联数据

	// 删除作业提交记录
	// 删除收藏
	return nil
}
