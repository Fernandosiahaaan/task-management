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

func (s *TaskService) GetTask(id *int64) (*model.Task, error) {
	existTask, err := s.Repo.GetTaskById(id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed get task from database. err = %s", err.Error())
	}
	return existTask, nil
}

func (s *TaskService) GetAllTask() ([]*model.Task, error) {
	existTask, err := s.Repo.GetAllTask()
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, nil
		}
		return nil, err
	}
	return existTask, nil
}

func (s *TaskService) CreateNewTask(task *model.Task) (int64, error) {
	if err := s.validateDueDate(task.DueDate); err != nil {
		return 0, err
	}

	existTask, err := s.Repo.GetTaskById(&task.Id)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return 0, err
	}

	if existTask != nil {
		return 0, errors.New("task already created")
	}
	task.UpdatedBy = task.CreatedBy

	return s.Repo.CreateNewTask(task)
}

func (s *TaskService) UpdateTask(task *model.Task) (*int64, error) {
	if err := s.validateDueDate(task.DueDate); err != nil {
		return nil, err
	}

	existTask, err := s.Repo.GetTaskById(&task.Id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("error from database. err = %s", err.Error())
	}

	if existTask == nil {
		return nil, errors.New("task not found")
	}
	return s.Repo.UpdateTask(task)
}

func (s *TaskService) DeleteTask(id *int64) error {
	err := s.Repo.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) validateDueDate(dueDate time.Time) error {
	now := time.Now()
	if dueDate.After(now) || dueDate.Equal(now) {
		return nil
	}
	return fmt.Errorf("due date time has passed.")
}
