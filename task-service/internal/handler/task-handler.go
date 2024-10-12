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
	"task-service/middleware"

	"github.com/gorilla/mux"
)

type ParamHandler struct {
	Service    *service.TaskService
	Ctx        context.Context
	ClientGrpc *grpc.ClientGrpc
	RabbitMq   *rabbitmq.RabbitMq
	Redis      *reddis.Redis
	Midleware  *middleware.Middleware
}

type TaskHandler struct {
	service    *service.TaskService
	Ctx        context.Context
	cancel     context.CancelFunc
	clientGrpc *grpc.ClientGrpc
	rabbitMq   *rabbitmq.RabbitMq
	Redis      *reddis.Redis
	Midleware  *middleware.Middleware
}

func NewTaskHandler(param *ParamHandler) *TaskHandler {
	handlerCtx, handlerCancel := context.WithCancel(param.Ctx)
	return &TaskHandler{
		service:    param.Service,
		Ctx:        handlerCtx,
		cancel:     handlerCancel,
		clientGrpc: param.ClientGrpc,
		rabbitMq:   param.RabbitMq,
		Redis:      param.Redis,
		Midleware:  param.Midleware,
	}
}

func (s *TaskHandler) TaskCreate(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var err error
	tokenStr := r.Context().Value("jwtToken").(string)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: "failed parse body request"})
		return
	}

	if err = s.compareUser(tokenStr, task.CreatedBy); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: err.Error()})
		return
	}

	// Validate UUID of user Assigned & Created User
	err = s.clientGrpc.ValidateUserUUID(task.AssignedTo, task.CreatedBy)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.Response{Error: true, Message: err.Error()})
		return
	}

	taskId, err := s.service.CreateNewTask(&task)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: fmt.Sprintf("failed create new task. err = %s", err.Error())})
		return
	}
	task.Id = taskId
	msg := fmt.Sprintf("successfully created task %d", task.Id)

	var response model.Response = model.Response{Error: false, Message: msg, Data: task}
	s.sendDoubleResponse(w, http.StatusCreated, rabbitmq.ACTION_TASK_CREATE, response)
}

func (s *TaskHandler) TaskUpdate(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Context().Value("jwtToken").(string)
	var task model.Task
	var err error
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: "failed parse body request"})
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.Response{Error: true, Message: "failed task id"})
	}
	task.Id = int64(id)

	if err = s.compareUser(tokenStr, task.UpdatedBy); err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: err.Error()})
		return
	}

	// Validate UUID of user Assigned & Created User
	err = s.clientGrpc.ValidateUserUUID(task.AssignedTo, task.UpdatedBy)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.Response{Error: true, Message: err.Error()})
		return
	}

	taskId, err := s.service.UpdateTask(&task)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: err.Error()})
		return
	}
	task.Id = *taskId
	var response model.Response = model.Response{Error: false, Message: "Task updated", Data: task}
	s.sendDoubleResponse(w, http.StatusOK, rabbitmq.ACTION_TASK_UPDATE, response)
}

func (s *TaskHandler) TaskRead(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.Response{Error: true, Message: "failed task id"})
	}
	taskId := int64(id)

	task, err := s.service.GetTask(&taskId)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: err.Error()})
		return
	} else if task == nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: "not found task from db"})
		return
	}
	model.CreateResponseHttp(w, http.StatusOK, model.Response{Error: false, Message: fmt.Sprintf("Read Task %d", taskId), Data: task})
}

func (s *TaskHandler) TaskReadAll(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	tasks, err := s.service.GetAllTask()
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: err.Error()})
		return
	} else if tasks == nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: "not found all task from db"})
		return
	}
	model.CreateResponseHttp(w, http.StatusOK, model.Response{Error: false, Message: "Success Read Task", Data: tasks})
}

func (s *TaskHandler) TaskDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.Response{Error: true, Message: "failed task id"})
	}
	taskId := int64(id)

	err = s.service.DeleteTask(&taskId)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusBadRequest, model.Response{Error: true, Message: err.Error()})
		return
	}
	var response model.Response = model.Response{Error: false, Message: fmt.Sprintf("Success Delete Task %d", taskId)}
	s.sendDoubleResponse(w, http.StatusOK, rabbitmq.ACTION_TASK_DELETE, response)
}

// compare user login from jwt with user created/ user updated when update/create task
func (s *TaskHandler) compareUser(jwtToken string, userId string) error {
	loginInfo, err := s.Redis.GetLoginInfoFromRedis(jwtToken)
	if err != nil {
		return fmt.Errorf("failed validation user. err = %s", err.Error())
	}

	fmt.Println("login info id = ", loginInfo.Id)
	fmt.Println("user id = ", userId)

	if loginInfo.Id != userId {
		return fmt.Errorf("updated_by / created_by not equal with user login")
	}
	return nil
}
func (s *TaskHandler) sendDoubleResponse(w http.ResponseWriter, httpStatusCode int, actionRabitMq string, response model.Response) {
	messageQueue, err := model.ConvertResponseToStr(response)
	if err != nil {
		model.CreateResponseHttp(w, http.StatusInternalServerError, model.Response{Error: true, Message: fmt.Sprintf("failed send message to notification service. err = %s", err.Error())})
	}
	s.rabbitMq.SendMessage(rabbitmq.EXCHANGE_NAME_TaskService, actionRabitMq, messageQueue)
	model.CreateResponseHttp(w, httpStatusCode, response)
}

func (s *TaskHandler) Close() {
	s.cancel()
}
