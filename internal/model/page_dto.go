package model

type PageDTO struct {
	Page     int   `json:"page"`
	PageSize int   `json:pageSize`
	Total    int64 `json:"total"`
}