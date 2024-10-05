package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	grpc "user-service/internal/gRPC"
	"user-service/internal/model"
	"user-service/internal/reddis"
	"user-service/internal/service"
	"user-service/middleware"

	"github.com/dgrijalva/jwt-go"
)

type UserHandler struct {
	Service    *service.UserService
	Ctx        context.Context
	gRPCServer grpc.ServerGrpc
}

const (
	MessageFailedToken string = "Failed Token"
	MessageFailedRedis string = "Failed Login Session"
	MessageFailedJWT   string = "Invalid JWT Token"
)

func NewUserHandler(service *service.UserService, ctx context.Context, serverGrpc grpc.ServerGrpc) *UserHandler {
	go serverGrpc.StartListen()
	return &UserHandler{Service: service, Ctx: ctx, gRPCServer: serverGrpc}
}

func (s *UserHandler) UserCreate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed parse body request"})
		return
	}
	if user.Role == "" {
		user.Role = "user"
	}

	userId, err := s.Service.CreateNewUser(user)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}
	user.Id = userId
	model.CreateResponseHttp(w, http.StatusCreated, model.ResponseHttp{Error: false, Message: "User created", Data: user})
}

func (s *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed body request"})
		return
	}

	user, err := s.Service.GetUser(user)
	if err == nil {
		tokenString, err := middleware.CreateToken(user.Username, user.Password)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed created token"})
			return
		}

		// send login info and user info to redis
		ctx := context.Background()
		err = reddis.SetLoginInfoToRedis(ctx, tokenString, model.LoginCacheData{Id: user.Id, Username: user.Username})
		if err != nil {
			model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
			return
		}

		err = reddis.SetUserInfoToRedis(ctx, user)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
			return
		}

		dataResponse := model.LoginData{Token: tokenString}
		model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success login", Data: dataResponse})
		return
	}

	model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
}

func (s *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
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

	_, err := middleware.VerifyToken(bearerToken[1])
	if err != nil {
		model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s; err = %s", MessageFailedToken, err)})
		return
	}

	ctx, cancel := context.WithCancel(s.Ctx)
	defer cancel()
	err = reddis.RedisClient.Del(ctx, bearerToken[1]).Err()
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "Failed logout session"})
		return
	}

	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success logout session"})

}

func (s *UserHandler) UserGet(w http.ResponseWriter, r *http.Request) {
	// receive value form context midleware
	tokenStr := r.Context().Value("jwtToken").(string)
	userClaims := r.Context().Value("user").(jwt.MapClaims)

	// get login info from redis
	loginInfo, err := reddis.GetLoginInfoFromRedis(r.Context(), tokenStr)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s. Err = %s", MessageFailedRedis, err)})
		return
	}

	// get user info from db
	user, err := reddis.GetUserInfoFromRedis(r.Context(), loginInfo.Id)
	if err != nil {
		// get user from db
		userName, ok := userClaims["username"].(string)
		userPass, ok2 := userClaims["password"].(string)
		if !ok || !ok2 {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: MessageFailedJWT})
			return
		}
		user = model.User{Username: userName, Password: userPass}
		user, err = s.Service.GetUser(user)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "Invalid username and password"})
			return
		}
	}

	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "success get info me", Data: user})
}

func (s *UserHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var userRequest model.User

	// receive value form context midleware
	tokenStr := r.Context().Value("jwtToken").(string)
	userClaims := r.Context().Value("user").(jwt.MapClaims)

	// get login info from redis
	loginInfo, err := reddis.GetLoginInfoFromRedis(r.Context(), tokenStr)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s. Err = %s", MessageFailedRedis, err)})
		return
	}

	// get user info from db
	user, err := reddis.GetUserInfoFromRedis(r.Context(), loginInfo.Id)
	if err != nil {
		// Validasi token
		userName, ok := userClaims["username"].(string)
		userPass, ok2 := userClaims["password"].(string)
		if !ok || !ok2 {
			model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: MessageFailedJWT})
			return
		}
		user = model.User{Username: userName, Password: userPass}
		user, err = s.Service.GetUser(user)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "Invalid username or password"})
			return
		}
	}
	// Decode JSON body
	userRequest.Id = user.Id
	if err = json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("failed to decode body request, err = %s", err.Error())})
		return
	}

	// Update user
	userUpdated, err := s.Service.UpdateUser(userRequest)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: fmt.Sprintf("failed to update user %s, err: %s", userRequest.Username, err.Error())})
		return
	}

	if err = reddis.SetUserInfoToRedis(r.Context(), userUpdated); err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "success update user", Data: userUpdated})
}

func (s *UserHandler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Missing authorization header"})
		return
	}
	tokenString = tokenString[len("Bearer "):]

	_, err := middleware.VerifyToken(tokenString)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s. err = %s", MessageFailedToken, err)})
		return
	}
	fmt.Fprint(w, "Welcome to the the protected area")

}

func (s *UserHandler) UsersGetAll(w http.ResponseWriter, r *http.Request) {
	// receive value form context midleware
	tokenStr := r.Context().Value("jwtToken").(string)
	userClaims := r.Context().Value("user").(jwt.MapClaims)

	// get login info from redis
	_, err := reddis.GetLoginInfoFromRedis(r.Context(), tokenStr)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s. Err = %s", MessageFailedRedis, err)})
		return
	}
	// get user from db
	_, ok := userClaims["username"].(string)
	_, ok2 := userClaims["password"].(string)
	if !ok || !ok2 {
		model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Invalid jwt token"})
		return
	}
	users, err := s.Service.GetAllUsers()
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "Invalid username and password"})
		return
	}

	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "success get info me", Data: users})
}
