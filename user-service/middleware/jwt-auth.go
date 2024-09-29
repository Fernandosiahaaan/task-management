package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"task-management/user-service/internal/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(username string, password string) (string, error) {
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

func VerifyToken(tokenString string) (*jwt.Token, error) {
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

func AuthMiddleware(next http.Handler) http.Handler {
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

		token, err := VerifyToken(bearerToken[1])
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

		ctx := context.WithValue(r.Context(), "jwtToken", bearerToken[1])
		ctx2 := context.WithValue(ctx, "user", claims)
		r = r.WithContext(ctx2)
		next.ServeHTTP(w, r)
	})

}
