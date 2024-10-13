package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"task-service/infrastructure/reddis"
	"task-service/internal/model"
	"task-service/repository"
	"time"
)

type TaskService struct {
	repo   *repository.TaskRepository
	redis  *reddis.Redis
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTaskService(ctx context.Context, repo *repository.TaskRepository, redis *reddis.Redis) *TaskService {
	serviceCtx, serviceCancel := context.WithCancel(ctx)
	return &TaskService{
		ctx:    serviceCtx,
		cancel: serviceCancel,
		repo:   repo,
		redis:  redis,
	}
}

func (s *TaskService) GetTask(id *int64) (*model.Task, error) {
	taskInfo, err := s.redis.GetTaskInfoFromRedis(*id)
	if (err == nil) && (taskInfo != nil) {
		return taskInfo, nil
	}

	existTask, err := s.repo.GetTaskById(id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed get task from database. err = %s", err.Error())
	}
	return existTask, nil
}

func (s *TaskService) GetTasksByUSerID(userID string) ([]*model.Task, error) {
	existTask, err := s.repo.GetAllTaskByUserId(userID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed get task user %s from database. err = %s", userID, err.Error())
	}
	return existTask, nil
}

func (s *TaskService) GetAllTask() ([]*model.Task, error) {
	existTask, err := s.repo.GetAllTask()
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

	if err := s.validateStatus(task.Status); err != nil {
		return 0, err
	}

	existTask, err := s.repo.GetTaskById(&task.Id)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return 0, err
	}

	if existTask != nil {
		return 0, errors.New("task already created")
	}
	task.UpdatedBy = task.CreatedBy
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	task.Id, err = s.repo.CreateNewTask(task)
	if err != nil {
		return 0, fmt.Errorf("failed save task info to db. err %v", err)
	}
	if err = s.redis.SetTaskInfoToRedis(task); err != nil {
		return 0, fmt.Errorf("failed save task info to redis caching. err %v", err)
	}
	return task.Id, nil
}

func (s *TaskService) UpdateTask(task *model.Task) (*int64, error) {
	if err := s.validateDueDate(task.DueDate); err != nil {
		return nil, err
	}

	if err := s.validateStatus(task.Status); err != nil {
		return nil, err
	}

	existTask, err := s.repo.GetTaskById(&task.Id)
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

	id, err := s.repo.UpdateTask(task)
	if err != nil {
		return nil, fmt.Errorf("failed update task info to db. err = %v", err)
	}
	task.Id = *id
	if err = s.redis.SetTaskInfoToRedis(task); err != nil {
		return nil, fmt.Errorf("failed update task info in redis caching. err = %v", err)
	}
	return id, nil
}

func (s *TaskService) DeleteTask(id *int64) error {
	err := s.repo.DeleteTask(id)
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

func (s *TaskService) validateStatus(status string) error {
	switch status {
	case model.StatusDone:
		return nil
	case model.StatusInProgress:
		return nil
	case model.StatusHold:
		return nil
	default:
		return fmt.Errorf("unknown status %s task", status)
	}
}

func (s *TaskService) Close() {
	s.cancel()
}
