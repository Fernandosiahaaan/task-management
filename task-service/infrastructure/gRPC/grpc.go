package grpc

import (
	"context"
	"fmt"
	"os"
	loggrpc "task-service/infrastructure/gRPC/logGrpc"
	usergrpc "task-service/infrastructure/gRPC/userGrpc"
)

type ParamGrpc struct {
	Ctx context.Context
}

type GrpcComm struct {
	ctx            context.Context
	cancel         context.CancelFunc
	UserGrpcClient *usergrpc.ClientGrpc
	LogGrpcClient  *loggrpc.ClientGrpc
}

func NewGrpc(param ParamGrpc) (*GrpcComm, error) {
	grpcCtx, grpcCancel := context.WithCancel(param.Ctx)
	portUserGRPC := os.Getenv("GRPC_USER_PORT")
	if portUserGRPC == "" {
		portUserGRPC = "50052"
	}
	portLogGRPC := os.Getenv("GRPC_LOG_PORT")
	if portLogGRPC == "" {
		portLogGRPC = "50052"
	}

	userServer, err := usergrpc.ConnectToServerGrpc(usergrpc.ParamClientGrpc{
		Ctx:  grpcCtx,
		Port: portUserGRPC,
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect to user-gRPC-server. err = %s", err.Error())
	}

	logClient, err := loggrpc.ConnectToServerGrpc(loggrpc.ParamClientGrpc{
		Ctx:  param.Ctx,
		Port: portLogGRPC,
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect log-gRPC-client. err = %s", err.Error())
	}

	var grpc *GrpcComm = &GrpcComm{
		ctx:            grpcCtx,
		cancel:         grpcCancel,
		UserGrpcClient: userServer,
		LogGrpcClient:  logClient,
	}
	return grpc, nil
}

func (g *GrpcComm) Close() {
	g.UserGrpcClient.Close()
	g.LogGrpcClient.Close()
	g.cancel()
}
