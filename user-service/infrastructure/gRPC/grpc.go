package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "user-service/infrastructure/gRPC/user"
	"user-service/infrastructure/reddis"
	"user-service/internal/model"
	"user-service/internal/service"

	"google.golang.org/grpc"
)

type ParamServerGrpc struct {
	Ctx     context.Context
	Port    string
	Service *service.UserService
	Redis   *reddis.RedisCln
}

type ServerGrpc struct {
	ctx                               context.Context
	cancel                            context.CancelFunc
	listener                          net.Listener
	server                            *grpc.Server
	service                           *service.UserService
	redis                             *reddis.RedisCln
	pb.UnimplementedUserServiceServer // Tambahkan ini untuk memastikan implementasi
}

// NewConnect initializes the gRPC server connection.
func NewConnect(param ParamServerGrpc) (*ServerGrpc, error) {
	var err error
	grpcCtx, grpcCancel := context.WithCancel(param.Ctx)
	var client *ServerGrpc = &ServerGrpc{
		ctx:     grpcCtx,
		cancel:  grpcCancel,
		redis:   param.Redis,
		service: param.Service,
	}
	client.listener, err = net.Listen("tcp", fmt.Sprintf(":%s", param.Port))
	if err != nil {
		return nil, err
	}

	// Create gRPC server
	client.server = grpc.NewServer()
	pb.RegisterUserServiceServer(client.server, client) // Ubah ini menjadi `&client`

	return client, nil
}

// StartListen starts the gRPC server to listen for incoming requests.
func (s *ServerGrpc) StartListen() {
	fmt.Printf("🌐 Server GRPC is running on port %s...\n", s.listener.Addr().String())
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// Stop gracefully stops the gRPC server.
func (s *ServerGrpc) Stop() {
	s.server.GracefulStop()
	if s.cancel != nil {
		s.cancel()
	}
	fmt.Println("Server stopped gracefully.")
}

func (s *ServerGrpc) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	fmt.Printf("Received request for User ID: %s\n", req.UserId)

	// Initialize the user model with the provided User ID
	var user *model.User

	// Attempt to get user info from Redis cache
	user, err := s.redis.GetUserInfo(req.UserId)
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
	user2, err := s.service.GetUserById(user.Id)
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

func (s *ServerGrpc) Close() {
	s.server.Stop()
	s.cancel()
}
