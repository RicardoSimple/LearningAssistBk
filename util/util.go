package util

import (
	"github.com/gin-gonic/gin"
	"learning-assistant/consts"
	"math/rand"
	"strconv"
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
func GetPageParams(c *gin.Context) (page int, pageSize int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, _ = strconv.Atoi(pageStr)
	pageSize, _ = strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return
}
func GetQueryUint(c *gin.Context, key string) (uint, error) {
	val := c.Query(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
