package loggrpc

import (
	"context"
	"fmt"
	"time"

	logPB "task-service/infrastructure/gRPC/logGrpc/pb"
	"task-service/internal/model"

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

func (client *ClientGrpc) Close() {
	client.conn.Close()
	client.cancel()
}

func (client *ClientGrpc) SendTaskToLogging(timeout time.Duration, task *model.Task, userId string, actionType logPB.TaskAction) error {
	ctx, cancel := context.WithTimeout(client.ctx, timeout)
	defer cancel()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var taskLog *logPB.LogTaskRequest = &logPB.LogTaskRequest{
		UserId:    userId,
		TaskId:    task.Id,
		Action:    actionType,
		Timestamp: timestamp,
		After: &logPB.TaskDetails{
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate.String(),
			Status:      task.Status,
		},
	}

	response, err := client.client.LogTaskAction(ctx, taskLog)
	if err != nil {
		fmt.Println("failed response from log service = ", response)
	}

	return err
}
