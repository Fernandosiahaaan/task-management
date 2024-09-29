package main

import (
	"context"
	"fmt"
	"log"
	pb "task-management/task-service/internal/gRPC/user"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// Membuat koneksi ke server gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	// Membuat request ke server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Panggil RPC GetUser dengan ID user
	res, err := c.GetUser(ctx, &pb.GetUserRequest{Id: "1"})
	if err != nil {
		log.Fatalf("Could not get user: %v", err)
	}

	// Cetak hasil response
	fmt.Printf("User ID: %s, Username: %s, Email: %s\n", res.Id, res.Username, res.Email)
}
