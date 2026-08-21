package handler

import (
	"errors"
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
