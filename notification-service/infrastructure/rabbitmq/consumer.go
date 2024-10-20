package rabbitmq

import (
	"fmt"
	"log"
	"notification-service/internal/mail"
	"notification-service/internal/model"
	"os"

	"github.com/streadway/amqp"
)

type rabbitMq struct {
	URL   string
	Conn  *amqp.Connection
	Email *mail.Mail
}

// Helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Initialize RabbitMQ connection
func Init(email *mail.Mail) (*rabbitMq, error) {
	// Create a new instance of rabbitMq
	output := &rabbitMq{}

	// Retrieve credentials from environment variables
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	port := os.Getenv("RABBITMQ_PORT")

	// Format the RabbitMQ URL
	output.URL = fmt.Sprintf("amqp://%s:%s@localhost:%s/", username, password, port)
	output.Email = email

	// Establish connection to RabbitMQ
	var err error
	output.Conn, err = amqp.Dial(output.URL)
	if err != nil {
		return nil, err
	}

	// Return the RabbitMQ instance with connection
	return output, nil
}

// Method to receive messages from RabbitMQ
// Method to receive messages from RabbitMQ
func (r *rabbitMq) ReceiveMessage() {
	// Open a channel
	channel, err := r.Conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	// Declare a temporary queue with a random name
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

	// Goroutine to handle incoming messages
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			if err = r.Email.SendTaskMsgEmail(msg.RoutingKey, string(msg.Body)); err != nil {
				fmt.Println("❌ Failed send message rabbitmq to email. err = ", err)
			}

			fmt.Println("✔️ Success send message rabbitmq to email.")
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
