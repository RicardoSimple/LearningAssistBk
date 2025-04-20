package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"ar-app-api/handler/basic"
	"ar-app-api/service/auth"
	"ar-app-api/util"
)

// AuthMiddleware 登录鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		tokenString := c.GetHeader(util.Authotization_Header)
		if tokenString == "" {
			// websocket token信息
			if strings.Contains(c.Request.RequestURI, "/ws") {
				tokenString = c.GetHeader(util.WebSocketAuthorization)
			}
			if tokenString == "" {
				basic.AuthFailure(c)
				return
			}
			c.Writer.Header().Add(util.WebSocketAuthorization, tokenString)
		}

		// 解析Token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := util.ParseToken(tokenString)
		if err != nil {
			basic.AuthFailure(c)
			return
		}

		util.SetUserInContext(c, claims.ID, claims.UserName, claims.Email)
		// 更新用户登录时间+状态
		go auth.UpdateLoginStatus(context.Background(), claims.ID)
		c.Next()
	}
}
