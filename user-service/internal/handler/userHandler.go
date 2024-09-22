package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"task-management/user-service/internal/model"
	"task-management/user-service/internal/reddis"
	"task-management/user-service/internal/service"

	"github.com/dgrijalva/jwt-go"
)

type UserHandler struct {
	Service *service.UserService
	Ctx     context.Context
}

func NewUserHandler(service *service.UserService, ctx context.Context) *UserHandler {
	return &UserHandler{Service: service, Ctx: ctx}
}

func (s *UserHandler) UserCreate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var statusCode int
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		statusCode = http.StatusBadRequest
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Error:   true,
			Message: "failed parse body request",
		})
		return
	}

	userId, err := s.Service.CreateNewUser(user)
	if err != nil {
		statusCode = http.StatusBadRequest
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	user.Id = userId

	statusCode = http.StatusCreated
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.ResponseHttp{
		Message: "Product created",
		Data:    user,
	})
}

func (s *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	var statusCode int
	json.NewDecoder(r.Body).Decode(&user)

	user, err := s.Service.GetUser(user)
	if err == nil {
		tokenString, err := s.Service.CreateToken(user.Username, user.Password)
		if err != nil {
			statusCode = http.StatusInternalServerError
			w.WriteHeader(statusCode)
			fmt.Errorf("No username found")
		}

		// send session login to reddis data
		ctx := context.Background()
		err = reddis.RedisClient.Set(ctx, tokenString, user.Username, model.UserSessionTime).Err() // Set waktu kadaluarsa 30 menit
		if err != nil {
			fmt.Println("error saving login in redis. err = ", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.ResponseHttp{
				Message: "error saving token",
				Error:   true,
			})
			return
		}

		statusCode = http.StatusOK
		w.WriteHeader(statusCode)
		dataResponse := model.LoginData{Token: tokenString}
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Message: "Success login",
			Error:   false,
			Data:    dataResponse,
		})
		return
	}

	statusCode = http.StatusUnauthorized
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.ResponseHttp{
		Message: "Invalid login",
		Error:   true,
	})
}

func (s *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
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

	_, err := s.Service.VerifyToken(bearerToken[1])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Error:   true,
			Message: "Failed token",
		})
		return
	}

	ctx, cancel := context.WithCancel(s.Ctx)
	defer cancel()
	err = reddis.RedisClient.Del(ctx, bearerToken[1]).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Error:   true,
			Message: "Failed logout session",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.ResponseHttp{
		Error:   false,
		Message: "Success logout session",
	})
}

func (s *UserHandler) UserGet(w http.ResponseWriter, r *http.Request) {
	userClaims := r.Context().Value("user").(jwt.MapClaims)
	userName, ok := userClaims["username"].(string)
	userPass, ok2 := userClaims["password"].(string)
	if !ok || !ok2 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Message: "Invalid token",
			Error:   true,
		})
		return
	}

	user := model.User{Username: userName, Password: userPass}
	user, err := s.Service.GetUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Message: "Invalid username and password",
			Error:   true,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.ResponseHttp{
		Message: "success get info me",
		Data:    user,
	})
}

func (s *UserHandler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	tokenString = tokenString[len("Bearer "):]

	_, err := s.Service.VerifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}
	fmt.Fprint(w, "Welcome to the the protected area")

}
