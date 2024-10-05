package main

import (
	"fmt"
	"notification-service/internal/rabbitmq"
)

func main() {
	fmt.Println("== notification service ==")
	rabbitmq.Init()
}
