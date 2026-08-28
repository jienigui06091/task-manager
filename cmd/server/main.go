// package main 表示这是可执行程序的入口包；Go 会从 main 函数开始运行。
package main

// import (...) 把本文件使用的其他包集中导入。
import (
	"context"
	"errors"
	"fmt"
	"log"      // log 用于向终端输出日志或终止程序。
	"net/http" // net/http 提供 HTTP 标准状态码等基础能力。
	"os"
	"strconv"

	"github.com/gin-gonic/gin" // Gin 是用于构建 HTTP API 的第三方框架。
	"github.com/joho/godotenv"

	"task-manager/internal/cache"
	"task-manager/internal/database" // database 包负责建立数据库连接池。
	"task-manager/internal/handler"  // handler 包负责处理 HTTP 请求。
	"task-manager/internal/middleware"
	"task-manager/internal/repository" // repository 包负责直接读写数据库。
	"task-manager/internal/service"    // service 包负责业务规则。
)

// main 是 Go 程序固定的入口函数，启动程序时会自动调用。
func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	// := 同时声明 db、err 两个变量，并接收连接池和可能出现的错误。
	db, err := database.NewPostgresPool()
	// if 在条件为 true 时执行花括号内的语句；err != nil 表示创建失败。
	if err != nil {
		// log.Fatal 输出错误后以失败状态立即结束程序。
		log.Fatal(err)
	}
	// defer 登记延后调用：main 结束前关闭连接池，释放数据库资源。
	defer db.Close()

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("REDIS_DB must be an integer: %v", err)
	}

	//连接redis
	client := cache.NewRedis(
		os.Getenv("REDIS_ADDR"),
		os.Getenv("REDIS_USERNAME"),
		os.Getenv("REDIS_PASSWORD"),
		redisDB,
	)
	defer client.Close()

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Redis connected")

	// 创建仓储层对象，并让它持有数据库连接池。
	taskRepository := repository.NewTaskRepository(db)
	taskLogRepository := repository.NewTaskLogRepository()

	// 创建业务层对象，并把仓储层作为依赖传入。
	taskService := service.NewTaskService(db,
		taskRepository,
		taskLogRepository,
		client)

	// 创建 HTTP 处理器对象，并把业务层作为依赖传入。
	taskHandler := handler.NewTaskHandler(taskService)

	userRepository := repository.NewAdminUseRepository(db)
	userService := service.NewAdminUserService(userRepository)
	userHandler := handler.NewAdminUserHandler(userService)

	// gin.Default 创建路由引擎，同时启用默认日志与异常恢复中间件。
	r := gin.Default()
	r.Use(
		middleware.ErrorHandler(),
	)
	//不需要JWT
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", userHandler.Register)
		authGroup.POST("/login", userHandler.Login)
	}
	//需要JWT
	taskGroup := r.Group("/api/task")
	taskGroup.Use(middleware.AuthMiddleware())
	{
		// 注册查询全部任务的 GET 接口；taskHandler.GetTasks 是处理函数。
		taskGroup.GET("/page", taskHandler.GetTasks)
		// 注册创建任务的 POST 接口；taskHandler.CreateTask 是处理函数。
		taskGroup.POST("/create", taskHandler.CreateTask)
		// 注册更新任务的 POST 接口；taskHandler.UpdateTask 是处理函数。
		taskGroup.POST("/update", taskHandler.UpdateTask)
		//获取单个任务
		taskGroup.GET("/getById", taskHandler.GetTaskByID)
		//根据id集合删除任务
		taskGroup.DELETE("/deleteByList", taskHandler.DeleteByList)
		taskGroup.GET("/duplicate", taskHandler.Duplicate)
	}
	userGroup := r.Group("/api/user")
	userGroup.Use(middleware.AuthMiddleware())
	userGroup.GET("/page", userHandler.Page)

	// GET 注册 HTTP GET 路由：客户端访问 / 时执行后面的匿名函数。
	r.GET("/", func(c *gin.Context) {
		// JSON 以 HTTP 200 OK 状态返回 JSON；gin.H 是 map[string]any 的简写。
		c.JSON(http.StatusOK, gin.H{
			// message 是 JSON 字段名，右侧字符串是字段值。
			"message": "Task Manager API",
		})
	})

	// Println 在终端输出普通日志，提示服务准备监听的端口。
	log.Println("server running on :8089")

	// r.Run 开始监听端口；if 初始化语句中的 err 只在本 if 语句内可见。
	if err := r.Run(":8089"); err != nil {
		// 端口监听失败时输出错误并终止程序。
		log.Fatal(err)
	}
}
