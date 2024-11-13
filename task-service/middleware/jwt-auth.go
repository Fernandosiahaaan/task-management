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

type Middleware struct {
	redis  *reddis.Redis
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMidleware(ctx context.Context, redis *reddis.Redis) *Middleware {
	midlewareCtx, midlewareCancel := context.WithCancel(ctx)
	return &Middleware{
		ctx:    midlewareCtx,
		cancel: midlewareCancel,
		redis:  redis,
	}

}

func (m *Middleware) verifyToken(tokenString string) (*jwt.Token, error) {
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

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			model.CreateResponseHttp(w, r, http.StatusUnauthorized, model.Response{Error: true, Message: "Authentication header null"})
			return
		}

		bearerToken := strings.Split(authToken, " ")
		if len(bearerToken) != 2 {
			model.CreateResponseHttp(w, r, http.StatusUnauthorized, model.Response{Error: true, Message: "Invalid format token"})
			return
		}

		var jwtToken string = bearerToken[1]
		token, err := m.verifyToken(jwtToken)
		if err != nil {
			model.CreateResponseHttp(w, r, http.StatusUnauthorized, model.Response{Error: true, Message: fmt.Sprintf("Failed token. err = %s", err)})
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			model.CreateResponseHttp(w, r, http.StatusUnauthorized, model.Response{Error: true, Message: "Failed claims token"})
			return
		}

		_, err = m.redis.GetLoginInfoFromRedis(jwtToken)
		if err != nil {
			model.CreateResponseHttp(w, r, http.StatusUnauthorized, model.Response{Error: true, Message: fmt.Sprintf("Failed Login Session. Err = %s", err.Error())})
			return
		}

		ctx := context.WithValue(r.Context(), "jwtToken", jwtToken)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}

func (m *Middleware) Close() {
	m.cancel()
}
