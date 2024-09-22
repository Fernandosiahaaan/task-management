package service

import (
	"fmt"
	"task-management/user-service/internal/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	secretKey = []byte("secret -key")
)

func (s *UserService) CreateToken(username string, password string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"password": username,
			"exp":      time.Now().Add(model.UserSessionTime).Unix(),
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
