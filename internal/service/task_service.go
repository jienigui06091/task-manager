package service

import (
	"context"
	"errors"

	"task-manager/internal/model"
	"task-manager/internal/repository"
)

type TaskService struct {
	repository *repository.TaskRepository
}

func NewTaskService(repository *repository.TaskRepository) *TaskService {
	return &TaskService{
		repository: repository,
	}
}

func (s *TaskService) GetAllTasks(
	ctx context.Context,
) ([]model.Task, error) {
	return s.repository.GetAll(ctx)
}

func (s *TaskService) CreateTask(
	ctx context.Context,
	title string,
) (model.Task, error) {

	if title == "" {
		return model.Task{}, errors.New("title is required")
	}

	task := model.Task{
		Title:     title,
		Completed: false,
	}

	return s.repository.Create(ctx, task)
}
