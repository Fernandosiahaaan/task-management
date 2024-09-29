package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "task-management/user-service/internal/gRPC/user"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

// Implementasi GetUser RPC
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	fmt.Printf("Received request for User ID: %s\n", req.Id)
	// Dummy data
	return &pb.GetUserResponse{
		Id:       req.Id,
		Username: "john_doe",
		Email:    "john_doe@example.com",
	}, nil
}

func main() {
	// Membuat listener di port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Membuat server gRPC
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	fmt.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
