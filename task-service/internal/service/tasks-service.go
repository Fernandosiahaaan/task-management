package service

import (
	"database/sql"
	"errors"
	"fmt"
	"task-management/task-service/internal/model"
	"task-management/task-service/repository"
	"time"
)

type TaskService struct {
	Repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{Repo: repo}
}

func (s *TaskService) CreateNewTask(task *model.Task) (int64, error) {
	if err := s.validateDueDate(task.DueDate); err != nil {
		return 0, err
	}

	if task.Id <= 0 {
		return s.Repo.CreateNewTask(task)
	}

	existTask, err := s.Repo.GetTask(task)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return 0, err
	}

	if existTask != nil {
		return 0, errors.New("task already created")
	}

	return s.Repo.CreateNewTask(task)
}

func (s *TaskService) validateDueDate(dueDate time.Time) error {
	now := time.Now()
	if dueDate.After(now) || dueDate.Equal(now) {
		return nil
	}
	return fmt.Errorf("due date time has passed.")
}
