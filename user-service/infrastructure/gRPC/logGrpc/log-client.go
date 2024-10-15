package loggrpc

import (
	"context"
	"fmt"
	"time"
	logPB "user-service/infrastructure/gRPC/logGrpc/pb"

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
