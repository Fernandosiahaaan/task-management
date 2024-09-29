package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	grpc "task-management/task-service/internal/gRPC"
	"task-management/task-service/internal/model"
	"task-management/task-service/internal/service"
)

type TaskHandler struct {
	Service    *service.TaskService
	Ctx        context.Context
	ClientGrpc grpc.ClientGrpc
}

func NewTaskHandler(service *service.TaskService, ctx context.Context, clientGrpc grpc.ClientGrpc) *TaskHandler {
	return &TaskHandler{Service: service, Ctx: ctx, ClientGrpc: clientGrpc}
}

func (s *TaskHandler) TaskCreate(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed parse body request"})
		return
	}
	fmt.Println("task = ", task)

	// Validate UUID of user Assigned & Created User
	err := s.ClientGrpc.ValidateCreatedAndAssignedUUID(task.AssignedTo, task.CreatedBy)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	taskId, err := s.Service.CreateNewTask(&task)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}
	task.Id = taskId
	model.CreateResponseHttp(w, http.StatusCreated, model.ResponseHttp{Error: false, Message: "Task created", Data: task})
}
