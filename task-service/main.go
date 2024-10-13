package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"task-service/infrastructure/datadog"
	grpc "task-service/infrastructure/gRPC"
	"task-service/infrastructure/rabbitmq"
	"task-service/infrastructure/reddis"
	"task-service/internal/handler"
	services "task-service/internal/service"
	"task-service/middleware"
	"task-service/repository"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

func router(handler *handler.TaskHandler) {
	router := muxtrace.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}).Methods("GET")

	// Grup rute yang membutuhkan middleware
	api := router.PathPrefix("/").Subrouter()
	api.Use(handler.Midleware.AuthMiddleware)

	// Rute dengan middleware otentikasi
	api.HandleFunc("/task", handler.TaskCreate).Methods(http.MethodPost)
	api.HandleFunc("/task/{taskId}", handler.TaskRead).Methods(http.MethodGet)
	api.HandleFunc("/task/{taskId}", handler.TaskUpdate).Methods(http.MethodPut)
	api.HandleFunc("/task/{taskId}", handler.TaskDelete).Methods(http.MethodDelete)

	api.HandleFunc("/tasks", handler.TaskReadAll).Methods(http.MethodGet)
	api.HandleFunc("/tasks/{userId}", handler.TasksUserRead).Methods(http.MethodGet)

	portHttp := os.Getenv("PORT_HTTP")
	if portHttp == "" {
		portHttp = "4001"
	}
	localHost := fmt.Sprintf("localhost:%s", portHttp)
	fmt.Printf("🌐 %s\n", localHost)
	// err := http.ListenAndServe("localhost:4001", router)
	err := http.ListenAndServe(localHost,
		handlers.CORS(
			handlers.AllowedOrigins([]string{"*"}),                             // Allow all origins
			handlers.AllowedMethods([]string{"POST", "GET", "PUT", "OPTIONS"}), // Allow specific methods
			handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allow specific headers
		)(router))
	if err != nil {
		log.Fatalf("Could not start the server: %v", err)
	}
	fmt.Println("Server started. Listening on port 4000")
	if err != nil {
		fmt.Println("Could not start the server", err)
	}
	fmt.Println("Server started. Listenning on port 4000")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	datadog.Init()
	defer datadog.Close()

	repo, err := repository.NewTaskRepository(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	fmt.Println("🔥 Init Repository...")

	redis, err := reddis.NewReddisClient(ctx)
	if err != nil {
		log.Fatalf("Could not to redis server. err = %v", err)
	}
	defer redis.Close()
	fmt.Println("🔥 Init Redis...")

	rabbitmq, err := rabbitmq.Init(ctx)
	if err != nil {
		log.Fatalf("failed init rabbitmq, err = %v", err)
	}
	defer rabbitmq.Close()
	fmt.Println("🔥 Init Redis...")

	taskService := services.NewTaskService(ctx, repo, redis)
	defer taskService.Close()
	fmt.Println("🔥 Init Service...")

	var paramGrpc grpc.ParamClientGrpc = grpc.ParamClientGrpc{
		Ctx:  ctx,
		Port: os.Getenv("GRPC_PORT"),
	}
	clientGrpc, err := grpc.ConnectToServerGrpc(paramGrpc)
	if err != nil {
		log.Fatalf("Could not connect gRPC client. err = %s", err.Error())
	}
	defer clientGrpc.Close()
	fmt.Println("🔥 Init gRPC Client...")

	midleware := middleware.NewMidleware(ctx, redis)
	defer midleware.Close()
	fmt.Println("🔥 Init Midleware...")

	var paramHandler *handler.ParamHandler = &handler.ParamHandler{
		Service:    taskService,
		Ctx:        ctx,
		ClientGrpc: clientGrpc,
		RabbitMq:   rabbitmq,
		Redis:      redis,
		Midleware:  midleware,
	}
	taskHandler := handler.NewTaskHandler(paramHandler)
	defer taskHandler.Close()
	fmt.Println("🔥 Init Handler...")

	router(taskHandler)
}
