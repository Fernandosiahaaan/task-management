package midleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"task-management/user-service/internal/model"
	"task-management/user-service/service"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ResponseHttp{
				Error:   true,
				Message: "Authentication header null",
			})
			return
		}

		bearerToken := strings.Split(authToken, " ")
		if len(bearerToken) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ResponseHttp{
				Error:   true,
				Message: "Invalid format token",
			})
			return
		}
		token, err := service.VerifyToken(bearerToken[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ResponseHttp{
				Error:   true,
				Message: "Failed token",
			})
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

		ctx := context.WithValue(r.Context(), "user", claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}
