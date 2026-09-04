package model

type Messages struct {
	// ID 是任务的数据库主键；int64 是 64 位有符号整数。
	ID               int64  `json:"id"`
	ConversationId   int64  `json:"conversationId"`
	SenderId         int64  `json:"senderId"`
	Content          string `json:"content"`
	ReplyToMessageId int64  `json:"replyToMessageId"`
	MessageType      string `json:"messageType"`
	IsRecalled       bool   `json:"isRecalled"`
	Base
}
