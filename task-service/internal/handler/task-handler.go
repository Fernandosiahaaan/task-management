package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	grpc "task-service/infrastructure/gRPC"
	"task-service/infrastructure/rabbitmq"
	"task-service/infrastructure/reddis"
	"task-service/internal/model"
	"task-service/internal/service"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	Service    *service.TaskService
	Ctx        context.Context
	ClientGrpc grpc.ClientGrpc
	RabbitMq   *rabbitmq.RabbitMq
	Redis      *reddis.RedisCln
}

func NewTaskHandler(service *service.TaskService, ctx context.Context, clientGrpc grpc.ClientGrpc, rabbitmq *rabbitmq.RabbitMq, redis *reddis.RedisCln) *TaskHandler {
	return &TaskHandler{
		Service:    service,
		Ctx:        ctx,
		ClientGrpc: clientGrpc,
		RabbitMq:   rabbitmq,
		Redis:      redis,
	}
}

func (s *TaskHandler) TaskCreate(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var err error
	tokenStr := r.Context().Value("jwtToken").(string)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed parse body request"})
		return
	}
	fmt.Println("task = ", task)

	if err = s.compareUser(tokenStr, task.CreatedBy); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	// Validate UUID of user Assigned & Created User
	err = s.ClientGrpc.ValidateUserUUID(task.AssignedTo, task.CreatedBy)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	taskId, err := s.Service.CreateNewTask(&task)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: fmt.Sprintf("failed create new task. err = %s", err.Error())})
		return
	}
	task.Id = taskId
	msg := fmt.Sprintf("successfully created task %d", task.Id)

	model.CreateResponseHttp(w, http.StatusCreated, model.ResponseHttp{Error: false, Message: msg, Data: task})
	s.RabbitMq.SendMessage(rabbitmq.EXCHANGE_NAME_TaskService, rabbitmq.ACTION_TASK_CREATE, msg)
}

func (s *TaskHandler) TaskUpdate(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Context().Value("jwtToken").(string)
	var task model.Task
	var err error
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "failed parse body request"})
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed task id"})
	}
	task.Id = int64(id)

	if err = s.compareUser(tokenStr, task.UpdatedBy); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	// Validate UUID of user Assigned & Created User
	err = s.ClientGrpc.ValidateUserUUID(task.AssignedTo, task.UpdatedBy)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}

	taskId, err := s.Service.UpdateTask(&task)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}
	task.Id = *taskId
	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "Task updated", Data: task})
}

func (s *TaskHandler) TaskRead(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed task id"})
	}
	taskId := int64(id)

	task, err := s.Service.GetTask(&taskId)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	} else if task == nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "not found task from db"})
		return
	}
	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: fmt.Sprintf("Read Task %d", taskId), Data: task})
}

func (s *TaskHandler) TaskReadAll(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	tasks, err := s.Service.GetAllTask()
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	} else if tasks == nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: "not found all task from db"})
		return
	}
	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: "Success Read Task", Data: tasks})
}

func (s *TaskHandler) TaskDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.ResponseHttp{Error: true, Message: "failed task id"})
	}
	taskId := int64(id)

	err = s.Service.DeleteTask(&taskId)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.ResponseHttp{Error: true, Message: err.Error()})
		return
	}
	model.CreateResponseHttp(w, http.StatusOK, model.ResponseHttp{Error: false, Message: fmt.Sprintf("Success Delete Task %d", taskId)})
}

// compare user login from jwt with user created/ user updated when update/create task
func (s *TaskHandler) compareUser(jwtToken string, userId string) error {
	loginInfo, err := s.Redis.GetLoginInfoFromRedis(jwtToken)
	if err != nil {
		return fmt.Errorf("failed validation user. err = %s", err.Error())
	}

	if loginInfo.Id != userId {
		return fmt.Errorf("updated_by / created_by not equal with user login")
	}
	return nil
}
