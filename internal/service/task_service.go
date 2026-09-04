package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"task-manager/internal/apperror"
	"task-manager/internal/cache"
	"task-manager/internal/model"
)

var (
	ErrTaskNotFound     = errors.New("task not found")
	ErrUpdateTaskFailed = errors.New("failed to update task")
)

func GetAllTasks(
	ctx context.Context,
	userID int64,
	page, pageSize int,
) ([]model.Task, int64, error) {
	return model.GetTasks(ctx, userID, page, pageSize)
}

func CreateTask(ctx context.Context, userID int64, title string) (model.Task, error) {
	task := model.Task{
		UserId:    userID,
		Title:     title,
		Completed: false,
	}

	err := model.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := model.CreateTask(ctx, tx, &task); err != nil {
			return err
		}

		return model.CreateTaskLog(ctx, tx, task.ID, userID, "create")
	})
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func GetTaskByID(ctx context.Context, userID, taskID int64) (model.Task, error) {
	cacheKey := taskCacheKey(userID, taskID)
	if cache.RedisClient != nil {
		cachedTask, err := cache.RedisClient.Get(ctx, cacheKey).Bytes()
		switch {
		case err == nil:
			var task model.Task
			if err := json.Unmarshal(cachedTask, &task); err == nil {
				return task, nil
			}
			_ = cache.RedisClient.Del(ctx, cacheKey).Err()
		case err != redis.Nil:
			// Redis is an optimization. Read from PostgreSQL when it is unavailable.
		}
	}

	task, err := model.GetTaskByID(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Task{}, apperror.ErrNotFound
		}
		return model.Task{}, err
	}

	cacheTask(ctx, cacheKey, task)
	return task, nil
}

func UpdateTask(
	ctx context.Context,
	id, userID int64,
	title string,
	completed bool,
) (model.Task, error) {
	updatedTask, err := model.UpdateTask(ctx, model.Task{
		ID:        id,
		UserId:    userID,
		Title:     title,
		Completed: completed,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Task{}, apperror.ErrNotFound
		}
		return model.Task{}, apperror.ErrInvalid
	}

	deleteTaskCache(ctx, userID, id)
	return updatedTask, nil
}

func DeleteTasks(ctx context.Context, ids []int64, userID int64) error {
	if err := model.DeleteTasks(ctx, ids, userID); err != nil {
		return err
	}

	for _, id := range ids {
		deleteTaskCache(ctx, userID, id)
	}

	return nil
}

func DuplicateTask(ctx context.Context, userID, taskID int64) error {
	return model.DuplicateTask(ctx, userID, taskID)
}

const taskCacheTTL = 5 * time.Minute

func taskCacheKey(userID, taskID int64) string {
	return fmt.Sprintf("task:%d:%d", userID, taskID)
}

func cacheTask(ctx context.Context, cacheKey string, task model.Task) {
	if cache.RedisClient == nil {
		return
	}

	data, err := json.Marshal(task)
	if err != nil {
		return
	}

	_ = cache.RedisClient.Set(ctx, cacheKey, data, taskCacheTTL).Err()
}

func deleteTaskCache(ctx context.Context, userID, taskID int64) {
	if cache.RedisClient == nil {
		return
	}

	_ = cache.RedisClient.Del(ctx, taskCacheKey(userID, taskID)).Err()
}
