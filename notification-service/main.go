package main

import (
	"fmt"
	"log"
	rabbitmq "notification-service/internal/rabbitmq"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("== notification service ==")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rabbitmq, err := rabbitmq.Init()
	if err != nil {
		log.Fatalf("failed init rabbitmq, err = %v", err)
	}
	defer rabbitmq.Conn.Close()
	fmt.Println("🔥 Init Redis...")
	rabbitmq.ReceiveMessage()

}
