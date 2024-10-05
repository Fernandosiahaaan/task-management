package main

import (
	"fmt"
	"log"
	rabbitmq "notification-service/infrastructure/rabbitmq"
	"notification-service/internal/mail"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("== notification service ==")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mails, err := mail.Init(os.Getenv("MAIL_USERNAME"), os.Getenv("APP_GMAIL_KEY"))
	fmt.Println("🔥 Init SMTP...")
	rabbitmq, err := rabbitmq.Init(mails)
	if err != nil {
		log.Fatalf("failed init rabbitmq, err = %v", err)
	}
	defer rabbitmq.Conn.Close()
	fmt.Println("🔥 Init Redis...")
	rabbitmq.ReceiveMessage()

}
