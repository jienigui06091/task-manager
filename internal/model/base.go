package model

import (
	"time"
)

type Base struct {
	// CreatedAt 是创建时间；time.Time 是 time 包提供的时间类型。
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 是最后更新时间；JSON 字段名为 updated_at。
	UpdatedAt time.Time `json:"updated_at"`
}

func Offset (page int,pageSize int) int{
	return (page-1)*pageSize
}
