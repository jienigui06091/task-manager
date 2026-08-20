// package handler 存放 HTTP 请求处理函数，也称控制器或处理器。
package handler

// import (...) 导入 HTTP 状态码、Gin 框架和业务层。
import (
	"net/http" // net/http 提供 StatusOK、StatusBadRequest 等标准状态码常量。

	"github.com/gin-gonic/gin" // Gin 提供路由上下文和 JSON 响应方法。

	"task-manager/internal/service" // service 包封装创建、查询任务的业务规则。
)

// TaskHandler 是 HTTP 处理器；它持有业务服务，而不直接访问数据库。
type TaskHandler struct {
	// service 是业务服务指针，处理器通过它完成实际业务操作。
	service *service.TaskService
}

// NewTaskHandler 创建处理器，并保存传入的业务服务。
func NewTaskHandler(service *service.TaskService) *TaskHandler {
	// 创建 TaskHandler 结构体并返回它的地址。
	return &TaskHandler{
		// 左侧 service 是字段，右侧 service 是函数参数。
		service: service,
	}
}

// GetTasks 处理 GET /api/tasks 请求，并返回所有任务的 JSON 数据。
func (h *TaskHandler) GetTasks(c *gin.Context) {
	// c.Request.Context() 取得当前 HTTP 请求上下文，并传给业务层。
	tasks, err := h.service.GetAllTasks(c.Request.Context())
	// 数据库查询或业务处理失败时，进入错误响应分支。
	if err != nil {
		// 以 HTTP 500 状态返回 JSON；gin.H 是 JSON 对象的简写类型。
		c.JSON(http.StatusInternalServerError, gin.H{
			// error 字段向客户端说明请求失败。
			"error": "failed to get tasks",
		})
		// return 结束处理函数，避免继续发送成功响应。
		return
	}

	// 查询成功时，以 HTTP 200 状态返回包含任务列表的 JSON 对象。
	c.JSON(http.StatusOK, gin.H{
		// data 字段的值是从业务层得到的 tasks 切片。
		"data": tasks,
	})
}

// CreateTaskRequest 定义创建任务接口期望接收的 JSON 请求体结构。
type CreateTaskRequest struct {
	// Title 对应客户端 JSON 中的 title 字段。
	Title string `json:"title"`
}

// CreateTask 处理 POST /api/tasks 请求，验证 JSON 后创建任务。
func (h *TaskHandler) CreateTask(c *gin.Context) {
	// var 声明 req；零值 CreateTaskRequest 的 Title 是空字符串。
	var req CreateTaskRequest

	// ShouldBindJSON 解析请求体 JSON 并写入 req；&req 传入变量地址。
	if err := c.ShouldBindJSON(&req); err != nil {
		// JSON 格式不正确或不能映射到结构体时，以 HTTP 400 状态返回错误。
		c.JSON(http.StatusBadRequest, gin.H{
			// 向客户端说明请求体无效。
			"error": "invalid request",
		})
		// 结束函数，避免无效输入继续进入业务层。
		return
	}

	// 调用业务层创建任务，并接收新任务和可能的业务错误。
	task, err := h.service.CreateTask(
		// 将当前请求的上下文传下去。
		c.Request.Context(),
		// 传入从 JSON 请求体解析出的标题。
		req.Title,
	)

	// 例如标题为空时业务层会返回错误。
	if err != nil {
		// 以 HTTP 400 状态返回错误文本；err.Error() 将 error 转为字符串。
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		// 发送错误响应后立即结束函数。
		return
	}

	// 创建成功后，以 HTTP 201 Created 状态返回新任务的完整 JSON 数据。
	c.JSON(http.StatusCreated, gin.H{
		"data": task,
	})
}
