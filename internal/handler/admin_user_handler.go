package handler

import (
	"errors"
	"net/http"
	"strconv"
	"task-manager/internal/middleware"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type AdminUserHandler struct {
	service *service.AdminUserService
}

func NewAdminUserHandler(service *service.AdminUserService) *AdminUserHandler {
	return &AdminUserHandler{
		service: service,
	}
}

type RegisterUser struct {
	Usernam  string `json:"username"`
	Password string `json:"password"`
}

func (h *AdminUserHandler) Register(c *gin.Context) {
	var req RegisterUser
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "请检查入参",
		})
		return
	}
	username := req.Usernam
	password := req.Password

	if username == "" {
		c.JSON(400, gin.H{
			"errror": "用户名不能为空",
		})
		return
	}
	if password == "" {
		c.JSON(400, gin.H{
			"errror": "密码不能为空",
		})
		return
	}
	if len(password) < 6 || len(password) > 18 {
		c.JSON(400, gin.H{
			"errror": "密码长度请控制在6~18之间",
		})
		return
	}
	reqAdminuser, err := h.service.Register(c.Request.Context(), username, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			c.JSON(500, gin.H{
				"error": "注册的用户名已存在，请修改用户名",
			})
			return
		}
		c.JSON(500, gin.H{
			"error": "注册失败，请联系管理员",
		})
		return
	}
	c.JSON(200, gin.H{
		"data": reqAdminuser,
	})
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AdminUserHandler) Login(c *gin.Context) {
	//获取用户名和密码
	var req LoginReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "请检查入参",
		})
		return
	}
	token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "用户名或密码错误",
		})
		return
	}
	c.JSON(200, gin.H{
		"data": gin.H{
			"token": token,
		},
	})
}

func (h *AdminUserHandler) Page(c *gin.Context) {
	_, ok := middleware.GetUserId(c)
	if !ok {
		c.JSON(401, gin.H{
			"error": "token过期或不存在",
		})
		return
	}
	page := c.Query("page")
	pageSize := c.Query("pageSize")

	if page == "" || pageSize == "" {
		c.JSON(400, gin.H{
			"error": "请检查分页参数准确性",
		})
		return
	}
	iPage, err1 := strconv.Atoi(page)
	iPageSize, err2 := strconv.Atoi(pageSize)
	if err1 != nil || err2 != nil {
		c.JSON(400, gin.H{
			"error": "请检查分页参数准确性",
		})
		return
	}
	if iPage < 1 || iPageSize < 1 || iPageSize > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请检查分页参数准确性",
		})
		return
	}

	users, total, err := h.service.Page(c.Request.Context(), iPage, iPageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  users,
		"total": total,
	})
}
