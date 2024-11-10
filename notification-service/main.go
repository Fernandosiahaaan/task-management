package main

import (
	"context"
	"fmt"
	"log"
	grpc "notification-service/infrastructure/gRPC"
	rabbitmq "notification-service/infrastructure/rabbitmq"
	"notification-service/internal/mail"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// time.Sleep(5 * time.Second)
	fmt.Println("=== Notification Microservice ===")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var grpcCommParam grpc.ParamGrpc = grpc.ParamGrpc{Ctx: ctx}
	grpcCom, err := grpc.NewGrpc(grpcCommParam)
	if err != nil {
		log.Fatalf("failed init rabbitmq, err = %s", err.Error())
	}
	defer grpcCom.Close()
	fmt.Println("🔥 Init GRPC...")

	var mailParam mail.MailParams = mail.MailParams{
		Email:    os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("APP_GMAIL_KEY"),
		GRPC:     grpcCom,
	}
	mails, err := mail.Init(mailParam)
	if err != nil {
		log.Fatalf("failed init mails. err= %s", err)
	}
	fmt.Println("🔥 Init SMTP...")

	var rabitmqParams rabbitmq.RabbitMqParam = rabbitmq.RabbitMqParam{
		Email: mails,
		GRPC:  grpcCom,
	}
	rabbitmq, err := rabbitmq.Init(rabitmqParams)
	if err != nil {
		log.Fatalf("failed init rabbitmq, err = %v", err)
	}
	defer rabbitmq.Conn.Close()
	fmt.Println("🔥 Init RabbitMQ...")
	rabbitmq.ReceiveMessage()

}
