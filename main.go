package main

import (
	"context"
	"learning-assistant/handler/hash"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"learning-assistant/client"
	"learning-assistant/conf"
	"learning-assistant/dal"
	"learning-assistant/handler"
	"learning-assistant/handler/middleware"
	"learning-assistant/service/cron"
	"learning-assistant/util"
	"learning-assistant/util/log"
	"learning-assistant/ws"
)

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
	corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "Origin, Authorization, Content-Type")
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
	api := r.Group("/api/v1")
	{

		assignment := api.Group("/assignment")
		{
			assignment.POST("/create", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.CreateAssignmentHandler)
			assignment.GET("/listByCourse", middleware.AuthAlwaysAllow(), handler.GetAssignmentsByCourseHandler)
			assignment.GET("/listByTeacher", middleware.AuthAlwaysAllow(), handler.GetAssignmentsByTeacherHandler)
			assignment.GET("/detail", middleware.AuthAlwaysAllow(), handler.GetAssignmentDetailHandler)
			assignment.GET("/list", middleware.AuthAlwaysAllow(), handler.GetAssignments)
			assignment.POST("/delete", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.DeleteAssignmentHandler)
			assignment.GET("/my", middleware.AuthMiddlewareRequireRoles("student"), handler.GetCurrentUserAssignmentHandler)
			assignment.POST("/submit", middleware.AuthMiddlewareRequireRoles("student"), handler.SubmitAssignmentHandler)
			assignment.GET("/detail/full", middleware.AuthMiddlewareRequireRoles("admin", "student", "teacher"), handler.GetAssignmentDetailWithSubmissionHandler)
			assignment.GET("/submissions", middleware.AuthMiddlewareRequireRoles("student", "teacher", "admin"), handler.GetAssignmentSubmissionsHandler)
			assignment.POST("/evaluate", middleware.AuthMiddlewareRequireRoles("teacher"), handler.EvaluateAssignmentSubmissionHandler)
			assignment.POST("/algo/evaluate", middleware.AuthMiddlewareRequireRoles("admin", "teacher", "student"), handler.SmartEvaluateAssignment)
		}

		class := api.Group("/class")
		{
			class.POST("/create", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.CreateClassHandler)
			class.GET("/list", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.GetClassListHandler)
			class.GET("/all", middleware.AuthMiddlewareRequireRoles("admin", "teacher", "student"), handler.GetAllClassHandler)
			class.POST("/delete", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.DeleteClassHandler)
			class.POST("/bind", middleware.AuthMiddlewareRequireRoles("admin"), handler.BindTeacherToClassHandler)
			class.GET("/my", middleware.AuthMiddlewareRequireRoles("teacher"), handler.GetMyClassHandler)
			class.GET("/my/students", middleware.AuthMiddlewareRequireRoles("teacher"), handler.GetMyClassStudentsHandler)
			class.POST("/user/bind", middleware.AuthMiddlewareRequireRoles("admin", "teacher"), handler.BindUserToClassHandler)
		}

		account := api.Group("/account")
		{
			account.GET("/auth/current", middleware.AuthMiddleware(), handler.CurrentUser)
			account.POST("/auth/login", handler.Login)
			account.POST("/auth/register", handler.Register)
			account.GET("/auth/check", middleware.AuthMiddleware(), handler.CheckToken)
		}

		course := api.Group("/course")
		{
			course.POST("/create", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.CreateCourseHandler)
			course.POST("/subject/create", middleware.AuthMiddlewareRequireRoles("admin"), handler.CreateSubjectHandler)
			course.GET("/get", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.GetCoursesByPage)
			course.GET("/subject/getAll", handler.GetSubjects)
			course.POST("/delete", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.DeleteCourseHandler)
			course.GET("/courses", handler.GetCoursesHandler)
			course.GET("/detail", middleware.AuthAlwaysAllow(), handler.GetCourseDetailHandler)
			course.POST("/update", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.UpdateCourseHandler)
			course.GET("/view/increase", middleware.AuthAlwaysAllow(), handler.IncrementCourseViewHandler)
			course.POST("/favorite", middleware.AuthMiddlewareRequireRoles("teacher", "admin", "student"), handler.FavoriteCourseHandler)
			course.POST("/unfavorite", middleware.AuthMiddlewareRequireRoles("teacher", "admin", "student"), handler.UnfavoriteCourseHandler)
			course.GET("/hot", handler.FindHotNCourse)
			course.POST("/algo/detail", middleware.AuthMiddlewareRequireRoles("teacher", "admin"), handler.SmartCourseDetail)
			course.GET("/algo/hot", middleware.AuthAlwaysAllow(), handler.SmartNewCourses)
			course.GET("/byId/get", handler.GetCourseById)
		}

		user := api.Group("/user")
		{
			user.GET("/routes", handler.GetRoutesHandler)
			user.GET("/list", middleware.AuthMiddlewareRequireRoles("admin"), handler.GetUserListHandler)
			user.POST("/create", middleware.AuthMiddlewareRequireRoles("admin"), handler.CreateUserByAdmin)
			user.GET("/byType", middleware.AuthMiddlewareRequireRoles("admin", "teacher", "student"), handler.GetUserByType)
			user.POST("/update", middleware.AuthMiddlewareRequireRoles("admin"), handler.UpdateUserByAdmin)
			user.GET("/info", middleware.AuthMiddlewareRequireRoles("admin"), handler.GetUserInfoById)
			user.POST("/delete", middleware.AuthMiddlewareRequireRoles("admin"), handler.DeleteUserByAdmin)
		}

		chat := api.Group("/chat")
		{
			chat.POST("/assistant", middleware.AuthMiddlewareRequireRoles("admin", "teacher", "student"), handler.ChatAssistant)
			chat.GET("/conversations", middleware.AuthMiddlewareRequireRoles("admin", "teacher", "student"), handler.GetConversations)
			chat.GET("/messages", middleware.AuthMiddlewareRequireRoles("admin", "teacher", "student"), handler.GetMessages)
			chat.POST("/conversation/delete", middleware.AuthMiddlewareRequireRoles("admin", "teacher", "student"), handler.DeleteConversationHandler)
		}

		imageHash := api.Group("/image")
		{
			imageHash.POST("/hash/bind", hash.BindImageHash)
			imageHash.GET("/hash/similar", hash.SimilarImage)
		}

	}
	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8998")
}
