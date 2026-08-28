// package model 存放项目中使用的数据结构（也称模型）。
package model

// type Task struct 定义名为 Task 的结构体；结构体把多个相关字段组合成一个类型。
type Task struct {
	// ID 是任务的数据库主键；int64 是 64 位有符号整数。
	ID int64 `json:"id"`
	UserId int64 `json:"user_id"`
	// Title 是任务标题；string 是文本类型。
	// 反引号中的 json 标签规定 JSON 输出时的字段名为 title。
	Title string `json:"title"`
	// Completed 表示任务是否完成；bool 只能是 true 或 false。
	Completed bool `json:"completed"`
	Base
}
