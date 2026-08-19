package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"task-manager/internal/model"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, completed, created_at, updated_at
		FROM tasks
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var task model.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Create(
	ctx context.Context,
	task model.Task,
) (model.Task, error) {

	err := r.db.QueryRow(ctx, `
		INSERT INTO tasks (title, completed)
		VALUES ($1, $2)
		RETURNING id, title, completed, created_at, updated_at
	`,
		task.Title,
		task.Completed,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}
