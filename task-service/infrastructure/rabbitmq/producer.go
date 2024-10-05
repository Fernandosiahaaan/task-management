package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

const (
	EXCHANGE_NAME_TaskService string = "task-service"
)

const (
	ACTION_TASK_CREATE = "task.create"
	ACTION_TASK_READ   = "task.read"
	ACTION_TASK_UPDATE = "task.update"
	ACTION_TASK_DELETE = "task.delete"
)

type RabbitMq struct {
	URL  string
	Conn *amqp.Connection
}

// Helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Initialize RabbitMQ connection
func Init() (*RabbitMq, error) {
	// Create a new instance of rabbitMq
	output := &RabbitMq{}

	// Retrieve credentials from environment variables
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	port := os.Getenv("RABBITMQ_PORT")

	// Format the RabbitMQ URL
	output.URL = fmt.Sprintf("amqp://%s:%s@localhost:%s/", username, password, port)

	// Establish connection to RabbitMQ
	var err error
	output.Conn, err = amqp.Dial(output.URL)
	if err != nil {
		return nil, err
	}

	// Return the RabbitMQ instance with connection
	return output, nil
}

// Method to send a message to RabbitMQ
func (r *RabbitMq) SendMessage(exchangeName, action, message string) {
	// Open a channel
	channel, err := r.Conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	// Declare the exchange with type 'topic'
	err = channel.ExchangeDeclare(
		exchangeName, // exchange name
		"topic",      // exchange type (topic)
		true,         // durable
		false,        // auto-delete when unused
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Publish the message
	err = channel.Publish(
		exchangeName, // exchange name
		action,       // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	failOnError(err, "Failed to publish message")

	fmt.Printf(" [x] Sent %s: %s\n", action, message)
}
