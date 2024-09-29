package service

import (
	"database/sql"
	"errors"
	"fmt"
	"task-management/task-service/internal/model"
	"task-management/task-service/repository"
)

type TaskService struct {
	Repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{Repo: repo}
}

func (s *TaskService) CreateNewTask(task *model.Task) (int64, error) {
	fmt.Println("task = ", task)
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

func (s *TaskService) validateUserUUID(uuid string) {
}
