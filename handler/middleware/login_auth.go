package middleware

import (
	"ar-app-api/consts"
	"ar-app-api/service"
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"ar-app-api/handler/basic"
	"ar-app-api/util"
)

// AuthMiddleware 登录鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		tokenString := c.GetHeader(consts.Authotization_Header)
		if tokenString == "" {
			// websocket token信息
			if strings.Contains(c.Request.RequestURI, "/ws") {
				tokenString = c.GetHeader(consts.WebSocketAuthorization)
			}
			if tokenString == "" {
				basic.AuthFailure(c)
				return
			}
			c.Writer.Header().Add(consts.WebSocketAuthorization, tokenString)
		}

		// 解析Token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := util.ParseToken(tokenString)
		if err != nil {
			basic.AuthFailure(c)
			return
		}

		util.SetUserToGinContext(c, claims.ID, claims.UserName, claims.Email)
		// 更新用户登录时间+状态
		go service.UpdateLoginStatus(context.Background(), claims.ID)
		c.Next()
	}
}
