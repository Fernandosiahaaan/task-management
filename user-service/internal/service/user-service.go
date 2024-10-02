package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"task-management/user-service/internal/model"
	"task-management/user-service/repository"
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
	user.Password = strings.TrimSpace(user.Password)
	hashPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed hash password. err = %s", err.Error())
	}

	user.Password = hashPassword
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
		if errors.Is(err, sql.ErrNoRows) {
			return existUser, fmt.Errorf("username not found.")
		}
		return existUser, fmt.Errorf("error sql. err = %s", err.Error())
	}

	// Verifikasi apakah password cocok dengan hash
	match := s.VerifyPassword(user.Password, existUser.Password)
	if !match {
		return existUser, fmt.Errorf("password not equal")
	}

	return existUser, nil
}

func (s *UserService) GetUserById(user *model.User) (*model.User, error) {
	fmt.Println("user = ", user)
	existUser, err := s.Repo.GetUserById(*user)
	if err != nil && err != sql.ErrNoRows {
		return &existUser, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &existUser, nil
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	existUser, err := s.Repo.GetAllUsers()
	if err != nil {
		return existUser, err
	}
	return existUser, nil
}

func (s *UserService) UpdateUser(user model.User) (model.User, error) {
	user.Password = strings.TrimSpace(user.Password)
	hashPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("failed hash password. err = %s", err.Error())
	}
	user.Password = hashPassword

	id, err := s.Repo.UpdateUser(user)
	if err != nil {
		return model.User{}, err // Kembalikan model.User kosong jika ada error
	}
	user.Id = id
	return user, nil
}
