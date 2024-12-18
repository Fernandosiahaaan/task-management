package grpc

import (
	"context"
	"fmt"
	"os"
	loggrpc "user-service/infrastructure/gRPC/logGrpc"
	usergrpc "user-service/infrastructure/gRPC/userGrpc"
	"user-service/infrastructure/reddis"
	"user-service/internal/service"
)

type ParamGrpc struct {
	Ctx     context.Context
	Service *service.UserService
	Redis   *reddis.RedisCln
}

type GrpcComm struct {
	ctx            context.Context
	cancel         context.CancelFunc
	Service        *service.UserService
	Redis          *reddis.RedisCln
	UserGrpcServer *usergrpc.ServerGrpc
	LogGrpcClient  *loggrpc.ClientGrpc
}

func NewGrpc(param ParamGrpc) (*GrpcComm, error) {
	grpcCtx, grpcCancel := context.WithCancel(param.Ctx)
	userPortGRPC := os.Getenv("GRPC_PORT")
	if userPortGRPC == "" {
		userPortGRPC = "50052"
	}

	userServer, err := usergrpc.NewConnect(usergrpc.ParamServerGrpc{
		Ctx:     grpcCtx,
		Port:    userPortGRPC,
		Service: param.Service,
		Redis:   param.Redis,
	})
	if err != nil {
		return nil, fmt.Errorf("Could not connect to user-gRPC-server. err = %s", err.Error())
	}

	logPortGRPC := os.Getenv("GRPC_LOG_PORT")
	if logPortGRPC == "" {
		logPortGRPC = "50053"
	}
	logClient, err := loggrpc.ConnectToServerGrpc(loggrpc.ParamClientGrpc{
		Ctx:  param.Ctx,
		Port: logPortGRPC,
	})
	if err != nil {
		return nil, fmt.Errorf("Could not connect gRPC client. err = %s", err.Error())
	}

	var grpc *GrpcComm = &GrpcComm{
		ctx:            grpcCtx,
		cancel:         grpcCancel,
		Service:        param.Service,
		Redis:          param.Redis,
		UserGrpcServer: userServer,
		LogGrpcClient:  logClient,
	}
	return grpc, nil
}

func (g *GrpcComm) Start() {
	go g.UserGrpcServer.StartListen()
}

func (g *GrpcComm) Close() {
	g.UserGrpcServer.Close()
	g.LogGrpcClient.Close()
	g.cancel()
}
