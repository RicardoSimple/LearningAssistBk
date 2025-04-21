package util

import (
	"learning-assistant/consts"
	"math/rand"
	"time"
)

// GenerateInviteCode 生成班级邀请码（默认长度为6，可指定）
func GenerateInviteCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	// 用当前时间做种子，确保随机
	rand.Seed(time.Now().UnixNano())
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
func IsValidGrade(grade string) bool {
	for _, g := range consts.GradeOptions {
		if g == grade {
			return true
		}
	}
	return false
}
