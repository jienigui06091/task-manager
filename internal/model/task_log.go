package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type TaskLog struct {
	ID        int64     `json:"id"`
	TaskID    int64     `json:"task_id"`
	UserID    int64     `json:"user_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateTaskLog(ctx context.Context, tx *gorm.DB, taskID, userID int64, action string) error {
	return tx.WithContext(ctx).Create(&TaskLog{
		TaskID: taskID,
		UserID: userID,
		Action: action,
	}).Error
}
