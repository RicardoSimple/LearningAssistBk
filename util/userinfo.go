package util

import (
	"context"
	"errors"
)

type contextKey string

type UserInfo struct {
	ID       uint   `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}

const userContextKey contextKey = "user-key"

// GetUserFromContext 从上下文获取用户信息
func GetUserFromContext(ctx context.Context) (*UserInfo, error) {
	user, ok := ctx.Value(userContextKey).(*UserInfo)
	if !ok {
		return nil, errors.New("no user found in context")
	}
	return user, nil
}

// SetUserInContext 将用户信息存入上下文
func SetUserInContext(ctx context.Context, id uint, userName, email string) context.Context {
	return context.WithValue(ctx, userContextKey, &UserInfo{
		ID:       id,
		UserName: userName,
		Email:    email,
	})
}
