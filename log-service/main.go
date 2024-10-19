package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	grpc "log-service/infrastructure/gRPC"
	"log-service/repository"
)

func main() {
	fmt.Println("===== LOG SERVICE =====")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	repo, err := repository.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	fmt.Println("🔥 Init Repository...")

	logServer, err := grpc.NewConnect(ctx, repo)
	if err != nil {
		log.Fatal(err)
	}
	logServer.StartListen()

	defer logServer.Close()
	fmt.Println("🔥 Init Server Log GRPC...")

}
