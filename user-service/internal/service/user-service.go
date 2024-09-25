package service

import (
	"database/sql"
	"errors"
	"task-management/user-service/internal/model"
	"task-management/user-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateNewUser(user model.User) (string, error) {
	existUser, err := s.Repo.GetUser(user)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return "", err
	}

	if existUser.Username == user.Username {
		return "", errors.New("user already created")
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Id = uuid.New().String()
	return s.Repo.CreateNewUser(user)
}

func (s *UserService) GetUser(user model.User) (model.User, error) {
	existUser, err := s.Repo.GetUser(user)
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}

func (s *UserService) UpdateUser(user model.User) (model.User, error) {
	user.UpdatedAt = time.Now()
	id, err := s.Repo.UpdateUser(user)
	if err != nil {
		return model.User{}, err // Kembalikan model.User kosong jika ada error
	}
	user.Id = id
	return user, nil
}
