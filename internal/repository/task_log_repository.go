package repository

import (
	"context"
)

type TaskLogRepository struct {
}

func NewTaskLogRepository() *TaskLogRepository {
	return &TaskLogRepository{}
}

func (r *TaskLogRepository) Create(
	ctx context.Context,
	db DBTX,
	taskID int64,
	userID int64,
	action string,
) error {

	_, err := db.Exec(ctx, `
		INSERT INTO task_logs (
			task_id,
			user_id,
			action
		)
		VALUES ($1, $2, $3)
	`,
		taskID,
		userID,
		action,
	)

	return err
}
