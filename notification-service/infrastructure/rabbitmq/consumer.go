package rabbitmq

import (
	"fmt"
	"log"
	grpc "notification-service/infrastructure/gRPC"
	"notification-service/internal/mail"
	"notification-service/internal/model"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMqParam struct {
	Email *mail.Mail
	GRPC  *grpc.GrpcComm
}

type RabbitMq struct {
	url   string
	Conn  *amqp.Connection
	email *mail.Mail
	grpc  *grpc.GrpcComm
}

// Helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Initialize RabbitMQ connection
func Init(params RabbitMqParam) (*RabbitMq, error) {
	output := &RabbitMq{}
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	output.url = fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
	output.email = params.Email
	output.grpc = params.GRPC

	// Establish connection to RabbitMQ
	var err error
	output.Conn, err = amqp.Dial(output.url)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// Method to receive messages from RabbitMQ
func (r *RabbitMq) ReceiveMessage() {
	channel, err := r.Conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"",    // leave the name empty to let RabbitMQ generate a random queue name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Bind the queue to the topic exchange and use routing key pattern "tas.*"
	err = channel.QueueBind(
		queue.Name,                      // queue name
		"task.*",                        // routing key pattern (task.create, task.update, etc.)
		model.EXCHANGE_NAME_TaskService, // exchange name
		false,
		nil,
	)
	failOnError(err, "Failed to bind the queue to the exchange")

	// Start consuming messages from the queue
	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer tag
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			if err = r.email.SendTaskMsgEmail(msg.RoutingKey, string(msg.Body)); err != nil {
				fmt.Println("❌ Failed send message rabbitmq to email. err = ", err)
			}

			fmt.Println("✔️ Success send message rabbitmq to email.")
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
