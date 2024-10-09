package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	grpc "user-service/infrastructure/gRPC"
	"user-service/infrastructure/reddis"
	"user-service/internal/model"
	"user-service/internal/service"
	"user-service/middleware"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	MessageFailedToken     string = "Failed Token"
	MessageFailedRedis     string = "Failed Login Session"
	MessageFailedJWT       string = "Invalid JWT Token"
	MessageFailedReqUserId string = "Invalid User ID uri"
)

type ParamHandler struct {
	Service    *service.UserService
	Ctx        context.Context
	GrpcServer *grpc.ServerGrpc
	Redis      *reddis.RedisCln
	Midleware  *middleware.Midleware
}

type UserHandler struct {
	service    *service.UserService
	Ctx        context.Context
	cancel     context.CancelFunc
	grpcServer *grpc.ServerGrpc
	Redis      *reddis.RedisCln
	Midleware  *middleware.Midleware
}

func NewUserHandler(param ParamHandler) *UserHandler {
	handlerCtx, handlerCancel := context.WithCancel(param.Ctx)
	return &UserHandler{
		service:    param.Service,
		Ctx:        handlerCtx,
		cancel:     handlerCancel,
		grpcServer: param.GrpcServer,
		Redis:      param.Redis,
		Midleware:  param.Midleware,
	}
}

func (handler *UserHandler) UserCreate(w http.ResponseWriter, r *http.Request) {
	var user model.User
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed parse body request"})
		return
	}
	if user.Role == "" {
		user.Role = "user"
	}

	userId, err := handler.service.CreateNewUser(user)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}
	user.Id = userId
	model.CreateResponseHttp(w, http.StatusCreated, model.ResponseHttp{Error: false, Message: "User created", Data: user})
}

func (handler *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed body request"})
		return
	}

	user, err := handler.service.GetUser(user)
	if err == nil {
		tokenString, err := handler.Midleware.CreateToken(user.Username, user.Password)
		if err != nil {
			model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed created token"})
			return
		}

		// send login info and user info to redis
		ctx := context.Background()
		err = handler.Redis.SetLoginInfo(ctx, tokenString, model.LoginCacheData{Id: user.Id, Username: user.Username})
		if err != nil {
			model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
			return
		}

		err = handler.Redis.SaveUserInfo(user)
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

func (handler *UserHandler) UserLogout(w http.ResponseWriter, r *http.Request) {
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

	jwtToken := bearerToken[1]
	_, err := handler.Midleware.VerifyToken(jwtToken)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s; err = %s", MessageFailedToken, err)})
		return
	}

	err = handler.Redis.DeleteLoginInfo(jwtToken)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "Failed logout session"})
		return
	}

	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success logout session"})

}

func (handler *UserHandler) UserGet(w http.ResponseWriter, r *http.Request) {
	// receive value form context midleware
	tokenStr := r.Context().Value("jwtToken").(string)

	// get login info from redis
	loginInfo, err := handler.Redis.GetLoginInfo(tokenStr)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s. Err = %s", MessageFailedRedis, err)})
		return
	}

	// validate uri
	vars := mux.Vars(r)
	userId := vars["user_id"]
	if loginInfo.Id != vars["user_id"] {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: MessageFailedReqUserId})
		return
	}

	// get user info from db
	user, err := handler.service.GetUserById(userId)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Invalid username and password. err = %v", err)})
		return
	} else if user == nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "not found user"})
		return
	}

	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "success get info me", Data: user})
}

func (handler *UserHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var userRequest model.User

	// receive value form context midleware
	tokenStr := r.Context().Value("jwtToken").(string)

	// get login info from redis
	loginInfo, err := handler.Redis.GetLoginInfo(tokenStr)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s. Err = %s", MessageFailedRedis, err)})
		return
	}

	// validate uri
	vars := mux.Vars(r)
	userId := vars["user_id"]
	if loginInfo.Id != userId {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: MessageFailedReqUserId})
		return
	}

	// get user info from service
	user, err := handler.service.GetUserById(userId)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("Invalid username and password. err = %v", err)})
		return
	}

	// Decode JSON body
	userRequest.Id = user.Id
	if err = json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("failed to decode body request, err = %s", err.Error())})
		return
	}

	// Update user
	userUpdated, err := handler.service.UpdateUser(userRequest)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "success update user", Data: userUpdated})
}

func (handler *UserHandler) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: "Missing authorization header"})
		return
	}
	tokenString = tokenString[len("Bearer "):]

	_, err := handler.Midleware.VerifyToken(tokenString)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusUnauthorized, model.ResponseHttp{Error: true, Message: fmt.Sprintf("%s. err = %s", MessageFailedToken, err)})
		return
	}
	fmt.Fprint(w, "Welcome to the the protected area")

}

func (handler *UserHandler) UsersGetAll(w http.ResponseWriter, r *http.Request) {
	// receive value form context midleware
	tokenStr := r.Context().Value("jwtToken").(string)
	userClaims := r.Context().Value("user").(jwt.MapClaims)

	// get login info from redis
	_, err := handler.Redis.GetLoginInfo(tokenStr)
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
	users, err := handler.service.GetAllUsers()
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "Invalid username and password"})
		return
	}

	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "success get info me", Data: users})
}

func (handler *UserHandler) Close() {
	handler.cancel()
}
