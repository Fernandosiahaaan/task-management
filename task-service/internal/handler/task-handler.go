package handler

import (
	"context"
	"task-management/task-service/service"
)

type UserHandler struct {
	Service *service.UserService
	Ctx     context.Context
}

func NewUserHandler(service *service.UserService, ctx context.Context) *UserHandler {
	return &UserHandler{Service: service, Ctx: ctx}
}
