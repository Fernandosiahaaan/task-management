package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"log-service/infrastructure/datadog"
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

	datadog.Init()
	defer datadog.Close()
	fmt.Println("🔥 Init Datadog...")

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
