package controller

import (
	"net/http"
	"strconv"
	"task-manager/internal/apperror"
	"task-manager/internal/model"
	"task-manager/internal/response"

	"github.com/gin-gonic/gin"
)

func GetConversationPage(c *gin.Context) {

	page := c.Query("page")
	pageSize := c.Query("pageSize")
	taskId, err := strconv.Atoi(c.Query("taskId"))
	if err != nil {
		_ = c.Error(apperror.New(400, http.StatusBadRequest, "任务id为必传值"))
		return
	}
	if page == "" || pageSize == "" {

		_ = c.Error(apperror.New(400, http.StatusBadRequest, "请检查分页参数准确性"))
		return
	}
	iPage, err1 := strconv.Atoi(page)
	iPageSize, err2 := strconv.Atoi(pageSize)
	if err1 != nil || err2 != nil {
		_ = c.Error(apperror.New(400, http.StatusBadRequest, "请检查分页参数准确性"))
		return
	}
	if iPage < 1 || iPageSize < 1 || iPageSize > 100 {
		_ = c.Error(apperror.New(400, http.StatusBadRequest, "请检查分页参数准确性"))
		return
	}

	conver, total, err := model.GetConversations(c.Request.Context(), taskId, iPage, iPageSize)
	if err != nil {
		_ = c.Error(apperror.New(200, http.StatusInternalServerError, "系统异常"))
		return
	}
	pageResq := model.Page{}
	pageResq.Page = iPage
	pageResq.PageSize = iPageSize
	pageResq.List = conver
	pageResq.Total = total
	response.Success(c, pageResq)

}
