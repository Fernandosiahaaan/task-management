package grpc

import (
	"context"
	"fmt"
	"os"

	usergrpc "notification-service/infrastructure/gRPC/userGrpc"
)

type ParamGrpc struct {
	Ctx context.Context
}

type GrpcComm struct {
	ctx            context.Context
	cancel         context.CancelFunc
	UserGrpcClient *usergrpc.ClientGrpc
}

func NewGrpc(param ParamGrpc) (*GrpcComm, error) {
	grpcCtx, grpcCancel := context.WithCancel(param.Ctx)

	hostUserGRPC := os.Getenv("GRPC_USER_HOST")
	if hostUserGRPC == "" {
		hostUserGRPC = "localhost"
	}
	portUserGRPC := os.Getenv("GRPC_USER_PORT")
	if portUserGRPC == "" {
		portUserGRPC = "50052"
	}

	userServer, err := usergrpc.ConnectToServerGrpc(usergrpc.ParamClientGrpc{
		Ctx:  grpcCtx,
		Host: hostUserGRPC,
		Port: portUserGRPC,
	})
	if err != nil {
		return nil, fmt.Errorf("could not connect to user-gRPC-server. err = %s", err.Error())
	}

	var grpc *GrpcComm = &GrpcComm{
		ctx:            grpcCtx,
		cancel:         grpcCancel,
		UserGrpcClient: userServer,
	}
	return grpc, nil
}

func (g *GrpcComm) Close() {
	g.UserGrpcClient.Close()
	g.cancel()
}
