package grpc

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

type ParamServerGrpc struct {
	Ctx  context.Context
	Port string
}

type ServerGrpc struct {
	Hostname string
	Ctx      context.Context
	Cancel   context.CancelFunc
	Listener net.Listener
	Server   *grpc.Server
}

// NewConnect initializes the gRPC server connection.
func NewConnect(param ParamServerGrpc) (client ServerGrpc, err error) {
	client.Listener, err = net.Listen("tcp", fmt.Sprintf(":%s", param.Port))
	if err != nil {
		return client, fmt.Errorf("Failed to listen: %v", err)
	}

	// Create gRPC server
	client.Server = grpc.NewServer()
	pb.RegisterUserServiceServer(client.Server, &server{})

	return client, nil
}

// StartListen starts the gRPC server to listen for incoming requests.
func (s *ServerGrpc) StartListen() {
	fmt.Printf("Server GRPC is running on port %s...\n", s.Listener.Addr().String())
	if err := s.Server.Serve(s.Listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Stop gracefully stops the gRPC server.
func (s *ServerGrpc) Stop() {
	s.Server.GracefulStop()
	if s.Cancel != nil {
		s.Cancel()
	}
	fmt.Println("Server stopped gracefully.")
}

// Implementasi GetUser RPC
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	fmt.Printf("Received request for User ID: %s\n", req.UserId)
	// Dummy data
	return &pb.GetUserResponse{
		UserId:   req.UserId,
		Username: "john_doe",
		Email:    "john_doe@example.com",
		IsError:  false,
		Message:  "success get data",
	}, nil
}
