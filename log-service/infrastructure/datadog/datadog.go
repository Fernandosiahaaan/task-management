package datadog

import (
	"fmt"
	"os"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Init() {
	host := fmt.Sprintf("localhost:%s", os.Getenv("DD_PORT"))
	tracer.Start(
		tracer.WithDebugMode(false),
		tracer.WithAgentAddr(host),
		tracer.WithServiceName("my-task-service"),
	)
}

func Close() {
	tracer.Stop()
}
