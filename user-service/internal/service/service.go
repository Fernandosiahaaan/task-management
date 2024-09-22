package service

import (
	"database/sql"
	"errors"
	"fmt"
	"task-management/user-service/internal/model"
	"task-management/user-service/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey = []byte("secret -key")
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
	fmt.Println("existUser.Username = ", existUser.Username)
	fmt.Println("user.Username = ", user.Username)

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

func (s *UserService) CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func (s *UserService) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	return token, nil
}
