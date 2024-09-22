package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task-management/user-service/internal/model"
	"task-management/user-service/service"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (s *UserHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Error:   true,
			Message: "failed parse body request",
		})
		return
	}

	userId, err := s.Service.CreateNewUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ResponseHttp{
			Error:   true,
			Message: err.Error(),
		})
		return
	}

	user.Id = userId

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.ResponseHttp{
		Message: "Product created",
		Data:    user,
	})
}

func (s *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// fmt.Printf("The request body is %v\n", r.Body)

	var user model.User
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println("[LoginHandler] The user request body = ", user)

	user, err := s.Service.GetUser(user)
	if err == nil {
		tokenString, err := s.Service.CreateToken(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Errorf("No username found")
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, tokenString)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(model.ResponseHttp{
		Message: "Invalid credentials",
		Error:   true,
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
