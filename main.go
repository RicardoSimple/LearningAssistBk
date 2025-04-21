package main

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ar-app-api/client"
	"ar-app-api/conf"
	"ar-app-api/dal"
	"ar-app-api/handler"
	"ar-app-api/handler/hash"
	"ar-app-api/handler/middleware"
	"ar-app-api/service/cron"
	"ar-app-api/util"
	"ar-app-api/util/log"
	"ar-app-api/ws"
)

var db = make(map[string]string)

func init() {
	ctx := context.Background()
	// 加载logger
	log.InitLogging("", log.DEBUG)
	log.Info("[INIT] main init")
	conf.Init(ctx)
	dal.Init(ctx)
	util.Init(ctx)
	client.Init(ctx)
	go cron.InitCronJob(ctx)
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowWebSockets = true
	//corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "Origin, Authorization, Content-Type")
	r.Use(cors.New(corsCfg))

	// ping 接口
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// ping websocket接口
	webSocket := r.Group("/ws", middleware.AuthMiddleware())
	{
		webSocket.GET("/ping", ws.PingWS)
		webSocket.GET("/getclient", ws.GetPushNews)
		webSocket.GET("/deleteclient", ws.DeleteClient)
	}

	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/swagger/doc.json")))
	// 监听接口 无需鉴权
	listener := r.Group("/listener")
	{
		listener.GET("")
	}
	apiNoAuth := r.Group("/api/v1")
	{
		account := apiNoAuth.Group("/account")
		{
			account.GET("/auth/current", middleware.AuthMiddleware(), handler.CurrentUser)
			account.POST("/auth/login", handler.Login)
			account.POST("/auth/register", handler.Register)
			account.GET("/auth/check", middleware.AuthMiddleware(), handler.CheckToken)
		}
		course := apiNoAuth.Group("/course")
		{
			course.GET("/courses", handler.GetCoursesHandler)
		}
	}
	api := r.Group("/api/v1", middleware.AuthMiddleware())
	{
		imageHash := api.Group("/image")
		{
			imageHash.POST("/hash/bind", hash.BindImageHash)
			imageHash.GET("/hash/similar", hash.SimilarImage)
		}
		course := api.Group("/course")
		{
			course.POST("/create", handler.CreateCourseHandler)
			course.POST("/subject/create", handler.CreateSubjectHandler)
		}
		class := api.Group("/class")
		{
			class.POST("/create", handler.CreateClassHandler)
		}
	}
	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8998")
}
