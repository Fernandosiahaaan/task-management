package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"task-management/task-service/internal/model"

	"github.com/dgrijalva/jwt-go"
)

func verifyToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	fmt.Println("secret key = ", secretKey)
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

		token, err := verifyToken(bearerToken[1])
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
