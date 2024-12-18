package usergrpc

import (
	"context"
	"fmt"
	"time"

	userPB "notification-service/infrastructure/gRPC/userGrpc/pb"

	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type ParamClientGrpc struct {
	Ctx  context.Context
	Host string
	Port string
}

type ClientGrpc struct {
	hostname string
	ctx      context.Context
	cancel   context.CancelFunc
	client   userPB.UserServiceClient
	conn     *grpc.ClientConn
}

func ConnectToServerGrpc(param ParamClientGrpc) (*ClientGrpc, error) {
	var err error
	var client *ClientGrpc = &ClientGrpc{}
	client.ctx, client.cancel = context.WithCancel(param.Ctx)

	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName("my-grpc-client"))
	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName("my-grpc-client"))

	client.hostname = fmt.Sprintf("%s:%s", param.Host, param.Port)
	client.conn, err = grpc.Dial(
		client.hostname,
		grpc.WithInsecure(),
		grpc.WithIdleTimeout(10*time.Second),
		grpc.WithStreamInterceptor(si),
		grpc.WithUnaryInterceptor(ui),
	)
	if err != nil {
		return nil, err
	}

	client.client = userPB.NewUserServiceClient(client.conn)
	return client, nil
}

func (client *ClientGrpc) RequestUserInfo(userId string, timeout time.Duration) (*userPB.GetUserResponse, error) {
	// Membuat request ke server
	ctx, cancel := context.WithTimeout(client.ctx, timeout)
	defer cancel()

	// Panggil RPC GetUser dengan ID user
	res, err := client.client.GetUser(ctx, &userPB.GetUserRequest{UserId: userId})
	if err != nil {
		return nil, fmt.Errorf("Could not get user response. err = %s", err.Error())
	} else if res.IsError {
		return nil, fmt.Errorf("%s", res.Message)
	}
	return res, nil
}

func (client *ClientGrpc) ValidateUserUUID(assignedTo string, createdBy string) error {
	// Validate UUID Assigned & Created User from user microservice
	_, err := client.RequestUserInfo(createdBy, 1*time.Second)
	if err != nil {
		return fmt.Errorf("failed uuid created_by/updated_by of task. err %s", err.Error())
	}
	_, err = client.RequestUserInfo(assignedTo, 1*time.Second)
	if err != nil {
		return fmt.Errorf("failed uuid assigned_to of task. err %s", err.Error())
	}
	return nil
}

func (client *ClientGrpc) Close() {
	client.conn.Close()
	client.cancel()
}
