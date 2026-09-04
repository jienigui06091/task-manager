package model

import (
	"time"
)

type ConversationMembers struct {
	// ID 是任务的数据库主键；int64 是 64 位有符号整数。
	ID                int64     `json:"id"`
	ConversationId    int64     `json:"conversationId"`
	LastReadMessageId int64     `json:"lastReadMessageId"`
	UserId            int64     `json:"userId"`
	JoinedAt          time.Time `json:"joinedAt"`
	Type              string    `json:"type"`
	TaskId            int64     `json:"taskId"`
	IsPinned          bool      `json:"isPinned"`
	IsMuted           bool      `json:"isMuted"`
}





