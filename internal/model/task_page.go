package model

type TaskPage struct {
	List     []Task `json:"list"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Total    int64  `json:"total"`
}
