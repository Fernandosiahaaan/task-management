package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"user-service/infrastructure/datadog"
	grpc "user-service/infrastructure/gRPC"
	"user-service/infrastructure/reddis"
	"user-service/internal/handler"
	"user-service/internal/service"
	"user-service/middleware"
	"user-service/repository"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

func router(userHandler *handler.UserHandler) {
	router := muxtrace.NewRouter()
	// router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}).Methods("GET")

	router.HandleFunc("/user", userHandler.UserCreate).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandler.UserLogin).Methods(http.MethodPost)
	router.Handle("/user/logout", userHandler.Midleware.AuthMiddleware(http.HandlerFunc(userHandler.UserLogout))).Methods(http.MethodPost)
	router.Handle("/user", userHandler.Midleware.AuthMiddleware(http.HandlerFunc(userHandler.UsersGetAll))).Methods(http.MethodGet)
	router.Handle("/user/{user_id}", userHandler.Midleware.AuthMiddleware(http.HandlerFunc(userHandler.UserGet))).Methods(http.MethodGet)
	router.Handle("/user/{user_id}", userHandler.Midleware.AuthMiddleware(http.HandlerFunc(userHandler.UserUpdate))).Methods(http.MethodPut)
	router.HandleFunc("/user/protected", userHandler.ProtectedHandler).Methods(http.MethodGet)

	portHttp := os.Getenv("PORT_HTTP")
	if portHttp == "" {
		portHttp = "4000"
	}
	localHost := fmt.Sprintf("localhost:%s", portHttp)
	fmt.Printf("🌐 HTTP Api %s\n", localHost)
	// err := http.ListenAndServe("localhost:4000", router)
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

	repo, err := repository.NewuserRepository(ctx)
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

	userService := service.NewUserService(ctx, redis, repo)
	defer userService.Close()
	fmt.Println("🔥 Init Service...")

	var paramGrpc grpc.ParamGrpc = grpc.ParamGrpc{
		Ctx:     ctx,
		Service: userService,
		Redis:   redis,
	}
	serverGrpc, err := grpc.NewGrpc(paramGrpc)
	if err != nil {
		log.Fatalf("Could not connect to gRPC server. err = %s", err.Error())
	}
	serverGrpc.Start()
	defer serverGrpc.Close()
	fmt.Println("🔥 Init gRPC Server...")

	mw := middleware.NewMidleware(ctx, redis)
	defer mw.Close()
	fmt.Println("🔥 Init midleware...")

	fmt.Println("🔥 Init Handler...")
	var paramHandler handler.ParamHandler = handler.ParamHandler{
		Service:    userService,
		Ctx:        ctx,
		GrpcServer: serverGrpc,
		Redis:      redis,
		Midleware:  mw,
	}
	userHandler := handler.NewUserHandler(paramHandler)
	defer userHandler.Close()

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
