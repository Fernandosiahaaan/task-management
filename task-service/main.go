package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"task-management/task-service/internal/handler"
	"task-management/task-service/internal/middleware"
	"task-management/task-service/internal/reddis"
	"task-management/task-service/internal/repository"
	services "task-management/task-service/internal/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func router(taskHandler *handler.TaskHandler) {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}).Methods("GET")

	router.Use(middleware.AuthMiddleware)
	router.HandleFunc("/task/create", taskHandler.TaskCreate).Methods(http.MethodPost)

	fmt.Println("🌐 localhost:4001")
	// err := http.ListenAndServe("localhost:4001", router)
	err := http.ListenAndServe("localhost:4001",
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

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	reddis.RedisClient, err = reddis.NewReddisClient(ctx)
	if err != nil {
		log.Fatalf("Could not to redis server. err = %v", err)
	}
	defer reddis.RedisClient.Close()

	fmt.Println("🔥 Init Repository...")
	repo := repository.NewTaskRepository(db, ctx)

	fmt.Println("🔥 Init Service...")
	taskService := services.NewTaskService(repo)

	fmt.Println("🔥 Init Handler...")
	taskHandler := handler.NewTaskHandler(taskService, ctx)

	router(taskHandler)

}
