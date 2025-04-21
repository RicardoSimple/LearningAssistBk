package dal

import (
	"context"

	"gorm.io/gorm"

	"learning-assistant/dal/schema"
)

// CreateUser 新增用户
func CreateUser(ctx context.Context, user *schema.User) (*schema.User, error) {
	err := DB.WithContext(ctx).Create(user).Error
	return user, err
}

// UpdateUser 更新用户信息
func UpdateUser(ctx context.Context, user *schema.User) error {
	return DB.WithContext(ctx).Save(user).Error
}
func UpdateUserWithGroups(ctx context.Context, user *schema.User) error {
	// 使用事务保证数据一致性
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 更新 User
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		// 处理多对多关系
		if err := tx.Model(user).Association("ChatGroups").Replace(user.ChatGroups); err != nil {
			return err
		}

		return nil
	})
}

// GetUserByID 通过ID查询用户
func GetUserByID(ctx context.Context, id uint) (*schema.User, error) {
	var user schema.User
	if err := DB.WithContext(ctx).Preload("ChatGroups").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 通过用户名查询用户
func GetUserByUsername(ctx context.Context, username string) (*schema.User, error) {
	var user schema.User
	if err := DB.WithContext(ctx).Preload("ChatGroups").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser 删除用户
func DeleteUser(ctx context.Context, id uint) error {
	return DB.WithContext(ctx).Delete(&schema.User{}, id).Error
}
