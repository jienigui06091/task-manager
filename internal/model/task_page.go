package model

type TaskPage struct {
	Data     []Task `json:"data"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Total    int64  `json:"total"`
}
