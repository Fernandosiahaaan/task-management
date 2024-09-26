package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"task-management/user-service/internal/model"
	"task-management/user-service/internal/reddis"

	"github.com/dgrijalva/jwt-go"
)

func (s *UserHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authToken := r.Header.Get("Authorization")
		if authToken == "" {
			s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Authentication header null"})
			return
		}

		bearerToken := strings.Split(authToken, " ")
		if len(bearerToken) != 2 {
			s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Invalid format token"})
			return
		}
		token, err := s.Service.VerifyToken(bearerToken[1])
		if err != nil {
			s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Failed token"})
			return
		}

		ctxBg := context.Background()
		_, err = reddis.RedisClient.Get(ctxBg, bearerToken[1]).Result()
		if err != nil {
			s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Not found token in reddis"})
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
