package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"task-manager/internal/database"
	"task-manager/internal/handler"
	"task-manager/internal/repository"
	"task-manager/internal/service"
)

func main() {
	db, err := database.NewPostgresPool()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	taskRepository := repository.NewTaskRepository(db)

	taskService := service.NewTaskService(taskRepository)

	taskHandler := handler.NewTaskHandler(taskService)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Task Manager API",
		})
	})

	r.GET("/api/tasks", taskHandler.GetTasks)
	r.POST("/api/tasks", taskHandler.CreateTask)

	log.Println("server running on :8089")

	if err := r.Run(":8089"); err != nil {
		log.Fatal(err)
	}
}