package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "task-management/user-service/internal/gRPC/user"
	"task-management/user-service/internal/model"
	"task-management/user-service/internal/reddis"
	"task-management/user-service/internal/service"

	"google.golang.org/grpc"
)

type ParamServerGrpc struct {
	Ctx     context.Context
	Port    string
	Service *service.UserService
}

type ServerGrpc struct {
	Hostname                          string
	Ctx                               context.Context
	Cancel                            context.CancelFunc
	Listener                          net.Listener
	Server                            *grpc.Server
	Service                           *service.UserService
	pb.UnimplementedUserServiceServer // Tambahkan ini untuk memastikan implementasi
}

// NewConnect initializes the gRPC server connection.
func NewConnect(param ParamServerGrpc) (client ServerGrpc, err error) {
	client.Ctx = param.Ctx
	client.Service = param.Service
	client.Listener, err = net.Listen("tcp", fmt.Sprintf(":%s", param.Port))
	if err != nil {
		return client, fmt.Errorf("Failed to listen: %v", err)
	}

	// Create gRPC server
	client.Server = grpc.NewServer()
	pb.RegisterUserServiceServer(client.Server, &client) // Ubah ini menjadi `&client`

	return client, nil
}

// StartListen starts the gRPC server to listen for incoming requests.
func (s *ServerGrpc) StartListen() {
	fmt.Printf("🌐 Server GRPC is running on port %s...\n", s.Listener.Addr().String())
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

func (s *ServerGrpc) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	fmt.Printf("Received request for User ID: %s\n", req.UserId)

	// Initialize the user model with the provided User ID
	var user model.User

	// Attempt to get user info from Redis cache
	user, err := reddis.GetUserInfoFromRedis(s.Ctx, req.UserId)
	if err == nil {
		// Return successful response if found in Redis
		return &pb.GetUserResponse{
			UserId:   req.UserId,
			Username: user.Username,
			Email:    user.Email,
			IsError:  false,
			Message:  "Successfully retrieved data from Redis cache",
		}, nil
	}

	// If Redis fails, fetch the user info from the database via service
	user.Id = req.UserId
	user2, err := s.Service.GetUserById(&user)
	if err != nil {
		return &pb.GetUserResponse{
			UserId:   "",
			Username: "",
			Email:    "",
			IsError:  true,
			Message:  fmt.Sprintf("Failed to retrieve data from user microservice: %s", err.Error()),
		}, nil
	}

	// Handle the case where the user is not found in the database
	if user2 == nil {
		return &pb.GetUserResponse{
			UserId:   "",
			Username: "",
			Email:    "",
			IsError:  true,
			Message:  "User ID not found in the database",
		}, nil
	}

	// Return successful response with data from the database
	return &pb.GetUserResponse{
		UserId:   req.UserId,
		Username: user2.Username, // Use the user2 data retrieved from the DB
		Email:    user2.Email,
		IsError:  false,
		Message:  "Successfully retrieved data from user microservice",
	}, nil
}
