package grpc

import (
	"context"
	"fmt"
	"time"

	pb "task-management/task-service/internal/gRPC/user"

	"google.golang.org/grpc"
)

type ParamClientGrpc struct {
	Ctx  context.Context
	Port string
}

type ClientGrpc struct {
	Hostname string
	Ctx      context.Context
	Cancel   context.CancelFunc
	Client   pb.UserServiceClient
	Conn     *grpc.ClientConn
}

func ConnectToServerGrpc(param ParamClientGrpc) (client ClientGrpc, err error) {
	client.Ctx, client.Cancel = context.WithCancel(param.Ctx)

	client.Hostname = fmt.Sprintf("localhost:%s", param.Port)
	client.Conn, err = grpc.Dial(client.Hostname, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return client, fmt.Errorf("Failed connected client grpc. err = %v", err)
	}

	client.Client = pb.NewUserServiceClient(client.Conn)
	return client, nil
}

func (client *ClientGrpc) RequestUserInfo(userId string, timeout time.Duration) (*pb.GetUserResponse, error) {
	// Membuat request ke server
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Panggil RPC GetUser dengan ID user
	res, err := client.Client.GetUser(ctx, &pb.GetUserRequest{UserId: userId})
	if err != nil {
		return nil, fmt.Errorf("Could not get user response. err = %s", err.Error())
	}
	return res, nil
	// fmt.Printf("User ID: %s, Username: %s, Email: %s\n; error = %d; message = %s", res.UserId, res.Username, res.Email, res.IsError, res.Message)
}

func (client *ClientGrpc) Stop() {
	client.Conn.Close()
	client.Cancel()
}
