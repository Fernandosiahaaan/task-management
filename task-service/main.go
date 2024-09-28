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
	"task-management/task-service/internal/repository"
	services "task-management/task-service/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func router(userHandler *handler.UserHandler) {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}).Methods("GET")

	fmt.Println("🌐 localhost:4001")
	err := http.ListenAndServe("localhost:4001", router)
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

	fmt.Println("🔥 Init Repository...")
	repo := repository.NewuserRepository(db, ctx)

	fmt.Println("🔥 Init Service...")
	taskService := services.NewUserService(repo)

	fmt.Println("🔥 Init Handler...")
	taskHandler := handler.NewUserHandler(taskService, ctx)

	router(taskHandler)

}
