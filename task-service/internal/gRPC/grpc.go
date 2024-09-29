package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "task-management/task-service/internal/model/user"

	"google.golang.org/grpc"
)

// Server struct untuk implementasi service
type server struct {
	pb.UnimplementedUserServiceServer
}

// Implementasi method GetUser
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Simulasi mendapatkan data user
	user := pb.GetUserResponse{
		Id:       req.Id,
		Username: "JohnDoe",
		Email:    "johndoe@example.com",
	}
	fmt.Printf("Received request for user ID: %s\n", req.Id)
	return &user, nil
}

func RunGrpc() {
	// Listen pada port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Membuat gRPC server
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	fmt.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
