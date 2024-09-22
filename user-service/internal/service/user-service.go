package service

import (
	"database/sql"
	"errors"
	"task-management/user-service/internal/model"
	"task-management/user-service/internal/repository"
	"time"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateNewUser(user model.User) (int64, error) {
	existUser, err := s.Repo.GetUser(user)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return 0, err
	}

	if existUser.Username == user.Username {
		return 0, errors.New("user already created")
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return s.Repo.CreateNewUser(user)

}

func (s *UserService) GetUser(user model.User) (model.User, error) {
	existUser, err := s.Repo.GetUser(user)
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}
