package model

import "context"

type Conversation struct {
	// ID 是任务的数据库主键；int64 是 64 位有符号整数。
	ID     int64  `json:"id"`
	Type   string `json:"type"`
	TaskId int64  `json:"taskId"`
	Base
}

func GetConversations(ctx context.Context, taskID int, page, pageSize int) ([]Conversation, int64, error) {
	var total int64
	query := DB.WithContext(ctx).Where("task_id = ?", taskID)
	if err := query.Model(&Conversation{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var conversations []Conversation
	if err := query.Order("id DESC").
		Limit(pageSize).
		Offset(Offset(page, pageSize)).
		Find(&conversations).Error; err != nil {
		return nil, 0, err
	}
	return conversations, total, nil
}
