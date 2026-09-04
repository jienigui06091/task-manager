package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"task-manager/internal/controller"
	"task-manager/internal/middleware"
)

func SetRouter(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	authGroup.POST("/register", controller.Register)
	authGroup.POST("/login", controller.Login)

	taskGroup := r.Group("/api/task")
	taskGroup.Use(middleware.AuthMiddleware())
	taskGroup.GET("/page", controller.GetTasks)
	taskGroup.POST("/create", controller.CreateTask)
	taskGroup.POST("/update", controller.UpdateTask)
	taskGroup.GET("/getById", controller.GetTaskByID)
	taskGroup.DELETE("/deleteByList", controller.DeleteByList)
	taskGroup.GET("/duplicate", controller.Duplicate)

	userGroup := r.Group("/api/user")
	userGroup.Use(middleware.AuthMiddleware())
	userGroup.GET("/page", controller.GetAdminUserPage)

	conversationGroup := r.Group("/api/conversation")
	conversationGroup.Use(middleware.AuthMiddleware())
	conversationGroup.GET("/page", controller.GetConversationPage)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Task Manager API"})
	})
}
