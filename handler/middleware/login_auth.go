package middleware

import (
	"context"
	"learning-assistant/consts"
	"learning-assistant/service"
	"strings"

	"github.com/gin-gonic/gin"

	"learning-assistant/handler/basic"
	"learning-assistant/util"
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

		util.SetUserToGinContext(c, claims.ID, claims.UserName, claims.Email, claims.UserType, claims.ClassId)
		// 更新用户登录时间+状态
		go service.UpdateLoginStatus(context.Background(), claims.ID)
		c.Next()
	}
}

// AuthAlwaysAllow 登录鉴权中间件 始终放行
func AuthAlwaysAllow() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		tokenString := c.GetHeader(consts.Authotization_Header)
		if tokenString == "" {
			// websocket token信息
			if strings.Contains(c.Request.RequestURI, "/ws") {
				tokenString = c.GetHeader(consts.WebSocketAuthorization)
			}
			if tokenString == "" {
				c.Next()
			} else {
				c.Writer.Header().Add(consts.WebSocketAuthorization, tokenString)
			}
		}

		// 解析Token
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := util.ParseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		util.SetUserToGinContext(c, claims.ID, claims.UserName, claims.Email, claims.UserType, claims.ClassId)
		// 更新用户登录时间+状态
		go service.UpdateLoginStatus(context.Background(), claims.ID)
		c.Next()
	}
}

// AuthMiddlewareRequireRoles 要求用户角色属于指定列表
func AuthMiddlewareRequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if token == "" {
			basic.RequestFailureWithCode(c, "未提供 token", 401)
			c.Abort()
			return
		}

		userInfo, err := util.ParseToken(token)
		if err != nil {
			basic.RequestFailureWithCode(c, "token 无效或已过期", 401)
			c.Abort()
			return
		}

		// 权限判断
		userRole := userInfo.UserType // 假设 token 里解析出的字段
		authorized := false
		for _, r := range allowedRoles {
			if userRole == r {
				authorized = true
				break
			}
		}

		if !authorized {
			basic.RequestFailureWithCode(c, "无权限访问", 403)
			c.Abort()
			return
		}

		// 放行并将用户信息保存
		util.SetUserToGinContext(c, userInfo.ID, userInfo.UserName, userInfo.Email, userRole, userInfo.ClassId)

		c.Next()
	}
}
