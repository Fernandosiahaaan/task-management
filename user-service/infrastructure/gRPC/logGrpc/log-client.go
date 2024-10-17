package loggrpc

import (
	"context"
	"fmt"
	"time"
	logPB "user-service/infrastructure/gRPC/logGrpc/pb"
	"user-service/internal/model"

	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type ParamClientGrpc struct {
	Ctx  context.Context
	Port string
}

type ClientGrpc struct {
	hostname string
	ctx      context.Context
	cancel   context.CancelFunc
	client   logPB.LogServiceClient
	conn     *grpc.ClientConn
}

func ConnectToServerGrpc(param ParamClientGrpc) (*ClientGrpc, error) {
	var err error
	var client *ClientGrpc = &ClientGrpc{}
	client.ctx, client.cancel = context.WithCancel(param.Ctx)

	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName("my-grpc-client"))
	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName("my-grpc-client"))

	client.hostname = fmt.Sprintf("localhost:%s", param.Port)
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

	client.client = logPB.NewLogServiceClient(client.conn)
	return client, nil
}

func (client *ClientGrpc) SendUserToLogging(timeout time.Duration, user *model.User, actionType logPB.UserAction) error {
	ctx, cancel := context.WithTimeout(client.ctx, timeout)
	defer cancel()

	timestamp := time.Now().String()
	var userLog *logPB.LogUserRequest = &logPB.LogUserRequest{
		UserId:    user.Id,
		Action:    actionType,
		Timestamp: timestamp,
		After: &logPB.UserDetails{
			UserId:   user.Id,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
		},
	}
	_, err := client.client.LogUserAction(ctx, userLog)

	return err
}

func (client *ClientGrpc) Close() {
	client.conn.Close()
	client.cancel()
}
