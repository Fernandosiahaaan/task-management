package main

import (
	"context"
	"fmt"
	"log"
	client "task-service/infrastructure/gRPC"
	"time"
)

func main() {

	var clientParam client.ParamClientGrpc = client.ParamClientGrpc{
		Ctx:  context.Background(),
		Port: "50051",
	}
	grpcClient, err := client.ConnectToServerGrpc(clientParam)
	if err != nil {
		log.Fatalf("failed connect. err = %s", err.Error())
	}
	res, err := grpcClient.RequestUserInfo("1", 5*time.Second)
	if err != nil {
		log.Fatal("failed request data to server grpc. err = ", err)
	}
	// Cetak hasil response
	fmt.Printf("User ID: %s, Username: %s, Email: %s\n; error = %d; message = %s", res.UserId, res.Username, res.Email, res.IsError, res.Message)
	grpcClient.Stop()
}
