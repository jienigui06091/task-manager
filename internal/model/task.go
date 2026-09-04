// package model 存放项目中使用的数据结构（也称模型）。
package model

import (
	"context"

	"gorm.io/gorm"
)

// type Task struct 定义名为 Task 的结构体；结构体把多个相关字段组合成一个类型。
type Task struct {
	// ID 是任务的数据库主键；int64 是 64 位有符号整数。
	ID     int64 `json:"id"`
	UserId int64 `json:"user_id"`
	// Title 是任务标题；string 是文本类型。
	// 反引号中的 json 标签规定 JSON 输出时的字段名为 title。
	Title string `json:"title"`
	// Completed 表示任务是否完成；bool 只能是 true 或 false。
	Completed bool `json:"completed"`
	Base
}

func GetTasks(ctx context.Context, userID int64, page, pageSize int) ([]Task, int64, error) {
	var total int64
	query := DB.WithContext(ctx).Where("user_id = ?", userID)
	if err := query.Model(&Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var tasks []Task
	if err := query.Order("updated_at DESC").
		Limit(pageSize).
		Offset(Offset(page, pageSize)).
		Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func GetTaskByID(ctx context.Context, userID, taskID int64) (Task, error) {
	var task Task
	if err := DB.WithContext(ctx).
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&task).Error; err != nil {
		return Task{}, err
	}
	return task, nil
}

func CreateTask(ctx context.Context, tx *gorm.DB, task *Task) error {
	return tx.WithContext(ctx).Create(task).Error
}

func UpdateTask(ctx context.Context, task Task) (Task, error) {
	result := DB.WithContext(ctx).
		Model(&Task{}).
		Where("id = ? AND user_id = ?", task.ID, task.UserId).
		Updates(map[string]any{
			"title":     task.Title,
			"completed": task.Completed,
		})
	if result.Error != nil {
		return Task{}, result.Error
	}
	if result.RowsAffected == 0 {
		return Task{}, gorm.ErrRecordNotFound
	}
	return GetTaskByID(ctx, task.UserId, task.ID)
}

func DeleteTasks(ctx context.Context, ids []int64, userID int64) error {
	return DB.WithContext(ctx).
		Where("id IN ? AND user_id = ?", ids, userID).
		Delete(&Task{}).Error
}

func DuplicateTask(ctx context.Context, userID, taskID int64) error {
	task, err := GetTaskByID(ctx, userID, taskID)
	if err != nil {
		return err
	}
	return DB.WithContext(ctx).Create(&Task{
		UserId:    userID,
		Title:     task.Title,
		Completed: task.Completed,
	}).Error
}
