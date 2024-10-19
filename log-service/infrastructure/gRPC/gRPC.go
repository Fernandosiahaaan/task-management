package grpc

import (
	"context"
	"fmt"
	"log"
	logPB "log-service/infrastructure/gRPC/logging/pb"
	"log-service/repository"
	"net"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type ServerGrpc struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	server   *grpc.Server
	repo     *repository.RepoMongo
	logPB.UnimplementedLogServiceServer
}

// NewConnect initializes the gRPC server connection.
func NewConnect(ctx context.Context, repo *repository.RepoMongo) (*ServerGrpc, error) {
	var err error
	grpcCtx, grpcCancel := context.WithCancel(ctx)

	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName("my-grpc-server"))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName("my-grpc-server"))
	var client *ServerGrpc = &ServerGrpc{
		ctx:    grpcCtx,
		cancel: grpcCancel,
		repo:   repo,
	}

	portGRPC := os.Getenv("GRPC_PORT")
	if portGRPC == "" {
		portGRPC = "50052"
	}
	client.listener, err = net.Listen("tcp", fmt.Sprintf(":%s", portGRPC))
	if err != nil {
		return nil, fmt.Errorf("failed listen log grpc. err %v", err)
	}

	// Create gRPC server
	client.server = grpc.NewServer(grpc.StreamInterceptor(si), grpc.UnaryInterceptor(ui))
	logPB.RegisterLogServiceServer(client.server, client)

	return client, nil
}

// StartListen starts the gRPC server to listen for incoming requests.
func (s *ServerGrpc) StartListen() {
	fmt.Printf("🌐 Server GRPC is running on port %s...\n", s.listener.Addr().String())
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *ServerGrpc) LogTaskAction(ctx context.Context, req *logPB.LogTaskRequest) (*logPB.LogResponse, error) {
	fmt.Printf("received log task id = %d; state = %s", req.TaskId, req.Action.String())

	logEntry := bson.M{
		"timestamp": req.Timestamp,
		"action":    req.Action.String(),
		"userId":    req.UserId,
		"taskId":    req.TaskId,
		"before":    req.Before,
		"after":     req.After,
	}
	result, err := s.repo.InsertTaskLog(logEntry)
	if err != nil {
		return &logPB.LogResponse{
			Success: false,
			Message: fmt.Sprintf("failed save to log task id = '%s'; action = %s", req.UserId, req.Action.String()),
		}, err
	}

	fmt.Println("result log task to mongo db = ", result)
	return &logPB.LogResponse{
		Success: true,
		Message: fmt.Sprintf("success receive task id = %d; action = %s", req.TaskId, req.Action.String()),
	}, nil
}

func (s *ServerGrpc) LogUserAction(ctx context.Context, req *logPB.LogUserRequest) (*logPB.LogResponse, error) {
	fmt.Printf("received log user id = '%s'; state = %s\n", req.UserId, req.Action.String())

	logEntry := bson.M{
		"timestamp": req.Timestamp,
		"action":    req.Action.String(),
		"userId":    req.UserId,
		"before":    req.Before,
		"after":     req.After,
	}
	result, err := s.repo.InsertUserLog(logEntry)
	if err != nil {
		return &logPB.LogResponse{
			Success: false,
			Message: fmt.Sprintf("failed save to log user id = '%s'; action = %s", req.UserId, req.Action.String()),
		}, err
	}

	fmt.Println("result log user to mongo db = ", result)

	return &logPB.LogResponse{
		Success: true,
		Message: fmt.Sprintf("success save to log user id = '%s'; action = %s", req.UserId, req.Action.String()),
	}, nil
}

// Stop gracefully stops the gRPC server.
func (s *ServerGrpc) Close() {
	s.server.GracefulStop()
	s.cancel()
	fmt.Println("Server gRPC Logger stopped gracefully.")
}
