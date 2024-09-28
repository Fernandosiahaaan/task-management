package handler

import (
	"context"
	"task-management/task-service/service"
)

type TaskHandler struct {
	Service *service.TaskService
	Ctx     context.Context
}

func NewTaskHandler(service *service.TaskService, ctx context.Context) *TaskHandler {
	return &TaskHandler{Service: service, Ctx: ctx}
}
