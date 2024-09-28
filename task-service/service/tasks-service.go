package service

import (
	"database/sql"
	"errors"
	"task-management/task-service/internal/model"
	"task-management/task-service/internal/repository"
	"time"
)

type TaskService struct {
	Repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{Repo: repo}
}

func (s *TaskService) CreateNewUser(task *model.Task) (*int64, error) {
	existTask, err := s.Repo.GetTask(task)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return nil, err
	}

	if existTask != nil {
		return nil, errors.New("task already created")
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	return s.Repo.CreateNewTask(task)
}
