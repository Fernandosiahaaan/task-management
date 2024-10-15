package usergrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	userPB "user-service/infrastructure/gRPC/userGrpc/pb"
	"user-service/infrastructure/reddis"
	"user-service/internal/service"

	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type ParamServerGrpc struct {
	Ctx     context.Context
	Port    string
	Service *service.UserService
	Redis   *reddis.RedisCln
}

type ServerGrpc struct {
	ctx                                   context.Context
	cancel                                context.CancelFunc
	listener                              net.Listener
	server                                *grpc.Server
	service                               *service.UserService
	redis                                 *reddis.RedisCln
	userPB.UnimplementedUserServiceServer // Tambahkan ini untuk memastikan implementasi
}

// NewConnect initializes the gRPC server connection.
func NewConnect(param ParamServerGrpc) (*ServerGrpc, error) {
	var err error
	grpcCtx, grpcCancel := context.WithCancel(param.Ctx)

	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName("my-grpc-server"))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName("my-grpc-server"))
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
	client.server = grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	userPB.RegisterUserServiceServer(client.server, client) // Ubah ini menjadi `&client`

	return client, nil
}

// StartListen starts the gRPC server to listen for incoming requests.
func (s *ServerGrpc) StartListen() {
	fmt.Printf("🌐 Server GRPC is running on port %s...\n", s.listener.Addr().String())
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *ServerGrpc) GetUser(ctx context.Context, req *userPB.GetUserRequest) (*userPB.GetUserResponse, error) {
	fmt.Printf("Received request for User ID: %s\n", req.UserId)

	// Attempt to get user info from Redis cache
	user, err := s.redis.GetUserInfo(req.UserId)
	if err == nil {
		// Return successful response if found in Redis
		return &userPB.GetUserResponse{
			UserId:   req.UserId,
			Username: user.Username,
			Email:    user.Email,
			IsError:  false,
			Message:  "Successfully retrieved data from Redis cache",
		}, nil
	}

	// If Redis fails, fetch the user info from the database via service
	user, err = s.service.GetUserById(req.UserId)
	if err != nil {
		return &userPB.GetUserResponse{
			UserId:   "",
			Username: "",
			Email:    "",
			IsError:  true,
			Message:  fmt.Sprintf("Failed to retrieve data from user microservice: %s", err.Error()),
		}, nil
	}

	// Handle the case where the user is not found in the database
	if user == nil {
		return &userPB.GetUserResponse{
			UserId:   "",
			Username: "",
			Email:    "",
			IsError:  true,
			Message:  "User ID not found in the database",
		}, nil
	}

	// Return successful response with data from the database
	return &userPB.GetUserResponse{
		UserId:   req.UserId,
		Username: user.Username, // Use the user data retrieved from the DB
		Email:    user.Email,
		IsError:  false,
		Message:  "Successfully retrieved data from user microservice",
	}, nil
}

func (s *ServerGrpc) Close() {
	s.server.GracefulStop()
	s.cancel()
	fmt.Println("Server gRPC Logger stopped gracefully.")
}
