package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"user-service/infrastructure/reddis"
	"user-service/internal/model"

	"github.com/dgrijalva/jwt-go"
)

type Midleware struct {
	ctx    context.Context
	cancel context.CancelFunc
	redis  *reddis.RedisCln
}

func NewMidleware(ctx context.Context, redis *reddis.RedisCln) *Midleware {
	midlewareCtx, midlewareCancel := context.WithCancel(ctx)
	var middleware *Midleware = &Midleware{
		ctx:    midlewareCtx,
		cancel: midlewareCancel,
		redis:  redis,
	}
	return middleware
}

func (m *Midleware) CreateToken(username string, password string) (string, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"password": password,
			"exp":      time.Now().Add(model.UserSessionTime).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func (m *Midleware) VerifyToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
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

func (m *Midleware) AuthMiddleware(next http.Handler) http.Handler {
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
		token, err := m.VerifyToken(jwtToken)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Failed token. err = %s", err)})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ResponseHttp{
				Error:   true,
				Message: "Fail claims token",
			})
			return
		}

		_, err = m.redis.GetLoginInfo(jwtToken)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Token session expired. err = %s", err)})
			return
		}

		ctx := context.WithValue(r.Context(), "jwtToken", jwtToken)
		ctx2 := context.WithValue(ctx, "user", claims)
		r = r.WithContext(ctx2)
		next.ServeHTTP(w, r)
	})

}

func (m *Midleware) Close() {
	m.cancel()
}
