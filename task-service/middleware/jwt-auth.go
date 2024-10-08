package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"task-service/infrastructure/reddis"
	"task-service/internal/model"

	"github.com/dgrijalva/jwt-go"
)

type middleware struct {
	Redis *reddis.RedisCln
	Ctx   context.Context
}

func NewMidleware(ctx context.Context, redis *reddis.RedisCln) *middleware {
	return &middleware{
		Ctx:   ctx,
		Redis: redis,
	}

}

func (m *middleware) verifyToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Authentication header null"})
			return
		}

		bearerToken := strings.Split(authToken, " ")
		if len(bearerToken) != 2 {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Invalid format token"})
			return
		}

		var jwtToken string = bearerToken[1]
		token, err := m.verifyToken(jwtToken)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Failed token. err = %s", err)})
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Failed claims token"})
			return
		}

		_, err = m.Redis.GetLoginInfoFromRedis(jwtToken)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Failed Login Session. Err = %s", err.Error())})
			return
		}

		ctx := context.WithValue(r.Context(), "jwtToken", jwtToken)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}
