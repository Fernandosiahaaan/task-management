package service

import (
	"database/sql"
	"errors"
	"fmt"
	"task-service/infrastructure/reddis"
	"task-service/internal/model"
	"task-service/repository"
	"time"
)

type TaskService struct {
	Repo  *repository.TaskRepository
	Redis *reddis.RedisCln
}

func NewTaskService(repo *repository.TaskRepository, redis *reddis.RedisCln) *TaskService {
	return &TaskService{
		Repo:  repo,
		Redis: redis,
	}
}

func (s *TaskService) GetTask(id *int64) (*model.Task, error) {
	taskInfo, err := s.Redis.GetTaskInfoFromRedis(*id)
	if (err == nil) && (taskInfo != nil) {
		return taskInfo, nil
	}

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
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	task.Id, err = s.Repo.CreateNewTask(task)
	if err != nil {
		return 0, fmt.Errorf("failed save task info to db. err %v", err)
	}
	if err = s.Redis.SetTaskInfoToRedis(task); err != nil {
		return 0, fmt.Errorf("failed save task info to redis caching. err %v", err)
	}
	return task.Id, nil
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
	task.UpdatedAt = time.Now()

	id, err := s.Repo.UpdateTask(task)
	if err != nil {
		return nil, fmt.Errorf("failed update task info to db. err = %v", err)
	}
	task.Id = *id
	if err = s.Redis.SetTaskInfoToRedis(task); err != nil {
		return nil, fmt.Errorf("failed update task info in redis caching. err = %v", err)
	}
	return id, nil
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
