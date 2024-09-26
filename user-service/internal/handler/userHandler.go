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

func (s *UserHandler) responseHttp(w http.ResponseWriter, statusCode int, response model.ResponseHttp) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	if response.Error {
		fmt.Printf("❌ status code = %d; message = %s\n", statusCode, response.Message)
		return
	}
	fmt.Printf("✔️  status code = %d; message = %s\n", statusCode, response.Message)
}

func (s *UserHandler) UserCreate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.responseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed parse body request"})
		return
	}

	userId, err := s.Service.CreateNewUser(user)
	if err != nil {
		s.responseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}
	user.Id = userId
	s.responseHttp(w, http.StatusCreated, model.ResponseHttp{Error: false, Message: "Product created", Data: user})
}

func (s *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.responseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed body request"})
	}
	fmt.Println("user = ", user)

	user, err := s.Service.GetUser(user)
	if err == nil {
		tokenString, err := s.Service.CreateToken(user.Username, user.Password)
		if err != nil {
			s.responseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed created token"})
		}

		// send login info and user info to redis
		ctx := context.Background()
		err = reddis.SetLoginInfoToRedis(ctx, tokenString, model.LoginCacheData{Id: user.Id, Username: user.Username})
		if err != nil {
			s.responseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		}

		err = reddis.SetUserInfoToRedis(ctx, user)
		if err != nil {
			s.responseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		}

		dataResponse := model.LoginData{Token: tokenString}
		s.responseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success login", Data: dataResponse})
		return
	}

	s.responseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "no username found"})
}

func (s *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
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

	_, err := s.Service.VerifyToken(bearerToken[1])
	if err != nil {
		s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Failed Token"})
		return
	}

	ctx, cancel := context.WithCancel(s.Ctx)
	defer cancel()
	err = reddis.RedisClient.Del(ctx, bearerToken[1]).Err()
	if err != nil {
		s.responseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "Failed logout session"})
		return
	}

	s.responseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success logout session"})

}

func (s *UserHandler) UserGet(w http.ResponseWriter, r *http.Request) {
	// receive value form context midleware
	tokenStr := r.Context().Value("jwtToken").(string)
	userClaims := r.Context().Value("user").(jwt.MapClaims)

	// get login info from redis
	loginInfo, err := reddis.GetLoginInfoFromRedis(r.Context(), tokenStr)
	if err != nil {
		s.responseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	// get user info from db
	user, err := reddis.GetUserInfoFromRedis(r.Context(), loginInfo.Id)
	if err != nil {
		// get user from db
		userName, ok := userClaims["username"].(string)
		userPass, ok2 := userClaims["password"].(string)
		if !ok || !ok2 {
			s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Invalid jwt token"})
			return
		}
		user = model.User{Username: userName, Password: userPass}
		user, err = s.Service.GetUser(user)
		if err != nil {
			s.responseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "Invalid username and password"})
			return
		}
	}

	s.responseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "success get info me", Data: user})
}

func (s *UserHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var userRequest model.User

	// receive value form context midleware
	tokenStr := r.Context().Value("jwtToken").(string)
	userClaims := r.Context().Value("user").(jwt.MapClaims)

	// get login info from redis
	loginInfo, err := reddis.GetLoginInfoFromRedis(r.Context(), tokenStr)
	if err != nil {
		s.responseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	// get user info from db
	user, err := reddis.GetUserInfoFromRedis(r.Context(), loginInfo.Id)
	if err != nil {
		// Validasi token
		userName, ok := userClaims["username"].(string)
		userPass, ok2 := userClaims["password"].(string)
		if !ok || !ok2 {
			s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Invalid token"})
			return
		}
		user := model.User{Username: userName, Password: userPass}
		user, err = s.Service.GetUser(user)
		if err != nil {
			s.responseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "Invalid username or password"})
			return
		}
	}
	// Decode JSON body
	userRequest.Id = user.Id
	if err = json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		s.responseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("failed to decode body request, err = %s", err.Error())})
		return
	}

	// Update user
	userUpdated, err := s.Service.UpdateUser(userRequest)
	if err != nil {
		s.responseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: fmt.Sprintf("failed to update user %s, err: %s", userRequest.Username, err.Error())})
		return
	}

	if err = reddis.SetUserInfoToRedis(r.Context(), userUpdated); err != nil {
		s.responseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	s.responseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "success update user", Data: userUpdated})
}

func (s *UserHandler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Missing authorization header"})
		return
	}
	tokenString = tokenString[len("Bearer "):]

	_, err := s.Service.VerifyToken(tokenString)
	if err != nil {
		s.responseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Invalid token"})
		return
	}
	fmt.Fprint(w, "Welcome to the the protected area")

}
