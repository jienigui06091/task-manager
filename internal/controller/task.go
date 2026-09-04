// package handler 存放 HTTP 请求处理函数，也称控制器或处理器。
package controller

// import (...) 导入 HTTP 状态码、Gin 框架和业务层。
import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin" // Gin 提供路由上下文和 JSON 响应方法。

	"task-manager/internal/apperror"
	"task-manager/internal/middleware"
	"task-manager/internal/model"
	"task-manager/internal/response"
	"task-manager/internal/service" // service 包封装创建、查询任务的业务规则。
)

// GetTasks 处理 GET /api/tasks 请求，并返回所有任务的 JSON 数据。
func GetTasks(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	userId, ok := middleware.GetUserId(c)
	if !ok {

		_ = c.Error(apperror.New(40101, http.StatusUnauthorized, "token过期或不存在"))
		return
	}
	if page == "" || pageSize == "" {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "请检查分页参数准确性"))
		return
	}
	iPage, err1 := strconv.Atoi(page)
	iPageSize, err2 := strconv.Atoi(pageSize)
	if err1 != nil || err2 != nil {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "请检查分页参数准确性"))
		return
	}
	if iPage < 1 || iPageSize < 1 || iPageSize > 100 {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "请检查分页参数准确性"))
		return
	}

	// c.Request.Context() 取得当前 HTTP 请求上下文，并传给业务层。
	tasks, total, err := service.GetAllTasks(c.Request.Context(), userId, iPage, iPageSize)
	// 数据库查询或业务处理失败时，进入错误响应分支。
	if err != nil {
		_ = c.Error(err)
		// return 结束处理函数，避免继续发送成功响应。
		return
	}

	pageResq := model.Page{}
	pageResq.List = tasks
	pageResq.Total = total
	pageResq.PageSize = iPageSize
	pageResq.Page = iPage

	// 查询成功时，以 HTTP 200 状态返回包含任务列表的 JSON 对象。
	response.Success(c, pageResq)
}

// CreateTaskRequest 定义创建任务接口期望接收的 JSON 请求体结构。
type CreateTaskRequest struct {
	// Title 对应客户端 JSON 中的 title 字段。
	Title string `json:"title"`
}

// CreateTask 处理 POST /api/tasks 请求，验证 JSON 后创建任务。
func CreateTask(c *gin.Context) {
	// var 声明 req；零值 CreateTaskRequest 的 Title 是空字符串。
	var req CreateTaskRequest

	// ShouldBindJSON 解析请求体 JSON 并写入 req；&req 传入变量地址。
	if err := c.ShouldBindJSON(&req); err != nil {
		// JSON 格式不正确或不能映射到结构体时，由错误中间件输出响应。
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "invalid request"))
		// 结束函数，避免无效输入继续进入业务层。
		return
	}

	//去除参数中的空格
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "invalid request"))
		return
	}
	// 调用业务层创建任务，并接收新任务和可能的业务错误。
	user_id, ok := middleware.GetUserId(c)
	if !ok {
		_ = c.Error(apperror.New(40101, http.StatusUnauthorized, "token过期或不存在"))
		return
	}
	task, err := service.CreateTask(
		// 将当前请求的上下文传下去。
		c.Request.Context(),
		user_id,
		// 传入从 JSON 请求体解析出的标题。
		req.Title,
	)

	// 例如标题为空时业务层会返回错误。
	if err != nil {
		_ = c.Error(err)
		// 发送错误响应后立即结束函数。
		return
	}

	// 创建成功后，以 HTTP 201 Created 状态返回新任务的完整 JSON 数据。
	response.Created(c, task)
}

type UpdateTaskRequest struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func UpdateTask(c *gin.Context) {
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "invalid request"))
		return
	}
	if req.ID <= 0 {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "任务的id不能小于等于0"))
		return
	}
	if req.Title == "" {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "任务的标题不能为空"))
		return
	}
	user_id, ok := middleware.GetUserId(c)
	if !ok {
		_ = c.Error(apperror.New(40101, http.StatusUnauthorized, "token过期或不存在"))
		return
	}
	task, err := service.UpdateTask(c.Request.Context(), req.ID, user_id, req.Title, req.Completed)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, task)
}

func GetTaskByID(c *gin.Context) {
	id := c.Query("id")
	taskID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "无效的任务ID"))
		return
	}
	user_id, ok := middleware.GetUserId(c)
	if !ok {
		_ = c.Error(apperror.New(40101, http.StatusUnauthorized, "token过期或不存在"))
		return
	}
	task, err := service.GetTaskByID(c.Request.Context(), user_id, taskID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, task)
}
func DeleteByList(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "invalid request"))
		return
	}
	if len(req.IDs) == 0 {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "任务ID列表不能为空"))
		return
	}
	user_id, ok := middleware.GetUserId(c)
	if !ok {
		_ = c.Error(apperror.New(40101, http.StatusUnauthorized, "token过期或不存在"))
		return
	}
	err := service.DeleteTasks(c.Request.Context(), req.IDs, user_id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, "删除成功")

}

func Duplicate(c *gin.Context) {
	newVar := c.Query("taskId")
	taskId, newVar2 := strconv.Atoi(newVar)
	if newVar2 != nil {
		_ = c.Error(apperror.New(40001, http.StatusBadRequest, "请传入正确的taskId"))
		return
	}
	user_id, ok := middleware.GetUserId(c)
	if !ok {
		_ = c.Error(apperror.New(40101, http.StatusUnauthorized, "token过期或不存在"))
		return
	}
	err := service.DuplicateTask(c.Request.Context(), user_id, int64(taskId))
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, "ok")

}
