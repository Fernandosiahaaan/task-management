package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-management/user-service/internal/handler"
	"task-management/user-service/internal/reddis"
	"task-management/user-service/internal/service"
	"task-management/user-service/middleware"
	"task-management/user-service/repository"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func router(userHandler *handler.UserHandler) {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}).Methods("GET")

	router.HandleFunc("/user/register", userHandler.UserCreate).Methods(http.MethodPost)
	router.HandleFunc("/user/login", userHandler.UserLogin).Methods(http.MethodPost)
	router.HandleFunc("/user/logout", userHandler.UserLogout).Methods(http.MethodPost)
	router.Handle("/user/aboutme", middleware.AuthMiddleware(http.HandlerFunc(userHandler.UserGet))).Methods(http.MethodGet)
	router.Handle("/users", middleware.AuthMiddleware(http.HandlerFunc(userHandler.UsersGetAll))).Methods(http.MethodGet)
	router.Handle("/user/update", middleware.AuthMiddleware(http.HandlerFunc(userHandler.UserUpdate))).Methods(http.MethodPut)
	router.HandleFunc("/user/protected", userHandler.ProtectedHandler).Methods(http.MethodGet)

	fmt.Println("🌐 localhost:4000")
	// err := http.ListenAndServe("localhost:4000", router)
	err := http.ListenAndServe("localhost:4000",
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
	repo := repository.NewuserRepository(db, ctx)

	fmt.Println("🔥 Init Service...")
	userService := service.NewUserService(repo)

	fmt.Println("🔥 Init Handler...")
	userHandler := handler.NewUserHandler(userService, ctx)

	router(userHandler)

	// Handle system interrupts (e.g., Ctrl+C) to gracefully shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		// Block until a signal is received
		<-sigChan
		fmt.Println("\nReceived shutdown signal")
		cancel() // Cancel the context to trigger cleanup
	}()
}
