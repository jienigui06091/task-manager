package model

// type Task struct 定义名为 Task 的结构体；结构体把多个相关字段组合成一个类型。
type AdminUser struct {
	// ID 是任务的数据库主键；int64 是 64 位有符号整数。
	ID int64 `json:"id"`
	// 反引号中的 json 标签规定 JSON 输出时的字段名为 title。
	PasswordHash string `json:"password_hash"`
	Username string `json:"username"`
	Base
}
