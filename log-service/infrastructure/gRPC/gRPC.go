package grpc

import (
	"context"
	"fmt"
	"log"
	logPB "log-service/infrastructure/gRPC/logging"
	"net"
	"os"

	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type ServerGrpc struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	server   *grpc.Server
	logPB.UnimplementedLogServiceServer
}

// NewConnect initializes the gRPC server connection.
func NewConnect(ctx context.Context) (*ServerGrpc, error) {
	var err error
	grpcCtx, grpcCancel := context.WithCancel(ctx)

	si := grpctrace.StreamServerInterceptor(grpctrace.WithServiceName("my-grpc-server"))
	ui := grpctrace.UnaryServerInterceptor(grpctrace.WithServiceName("my-grpc-server"))
	var client *ServerGrpc = &ServerGrpc{
		ctx:    grpcCtx,
		cancel: grpcCancel,
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

// Stop gracefully stops the gRPC server.
func (s *ServerGrpc) Close() {
	s.server.GracefulStop()
	s.cancel()
	fmt.Println("Server gRPC Logger stopped gracefully.")
}
