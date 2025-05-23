package util

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type contextKey string

type UserInfo struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	ClassID  uint   `json:"class_id"`
}

const userContextKey contextKey = "user-key"

// SetUserToGinContext 存入 gin.Context（推荐在中间件中使用）
func SetUserToGinContext(c *gin.Context, id uint, userName, email string, userType string, classID uint) {
	c.Set(string(userContextKey), &UserInfo{
		ID:       id,
		UserName: userName,
		Email:    email,
		UserType: userType,
		ClassID:  classID,
	})
}

// GetUserFromGinContext 从 gin.Context 获取
func GetUserFromGinContext(c *gin.Context) (*UserInfo, error) {
	val, exists := c.Get(string(userContextKey))
	if !exists {
		return nil, errors.New("user not found in gin context")
	}
	user, ok := val.(*UserInfo)
	if !ok {
		return nil, errors.New("user data format invalid")
	}
	return user, nil
}
