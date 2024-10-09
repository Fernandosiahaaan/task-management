package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"user-service/infrastructure/reddis"
	"user-service/internal/model"
	"user-service/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo   *repository.UserRepository
	ctx    context.Context
	cancel context.CancelFunc
	redis  *reddis.RedisCln
}

func NewUserService(ctx context.Context, redis *reddis.RedisCln, repo *repository.UserRepository) *UserService {
	serviceCtx, serviceCancel := context.WithCancel(ctx)
	return &UserService{
		repo:   repo,
		ctx:    serviceCtx,
		cancel: serviceCancel,
		redis:  redis,
	}
}

func (s *UserService) CreateNewUser(user model.User) (string, error) {
	user.Password = strings.TrimSpace(user.Password)
	hashPassword, err := s.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("failed hash password. err = %s", err.Error())
	}

	user.Password = hashPassword
	existUser, err := s.repo.GetUser(user)
	if err != nil && !errors.Is(sql.ErrNoRows, err) {
		return "", err
	}

	if existUser.Username == user.Username {
		return "", errors.New("user already created")
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Id = uuid.New().String()
	return s.repo.CreateNewUser(user)
}

func (s *UserService) GetUser(user model.User) (model.User, error) {
	existUser, err := s.repo.GetUser(user)
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

func (s *UserService) GetUserById(userId string) (*model.User, error) {
	exitUser, err := s.redis.GetUserInfo(userId)
	if (err == nil) && (exitUser != nil) {
		return exitUser, nil
	}

	existUser, err := s.repo.GetUserById(userId)
	if err != nil && err != sql.ErrNoRows {
		return existUser, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return existUser, nil
}

func (s *UserService) GetAllUsers() ([]model.User, error) {
	existUser, err := s.repo.GetAllUsers()
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
	user.UpdatedAt = time.Now()

	id, err := s.repo.UpdateUser(user)
	if err != nil {
		return model.User{}, fmt.Errorf("failed update user %s to db. err = %s", user.Username, err.Error())
	}
	user.Id = id

	if err = s.redis.SaveUserInfo(user); err != nil {
		return model.User{}, fmt.Errorf("failed update user %s to redis. err = %s", user.Username, err.Error()) // Kembalikan model.User kosong jika ada error
	}
	return user, nil
}

func (s *UserService) Close() {
	s.cancel()
}
