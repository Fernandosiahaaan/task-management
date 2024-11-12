package rabbitmq

import (
	"context"
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
	URL    string
	Conn   *amqp.Connection
	ctx    context.Context
	cancel context.CancelFunc
}

// Helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Initialize RabbitMQ connection
func Init(ctx context.Context) (*RabbitMq, error) {
	rabbitmqCtx, rabbitmqCancel := context.WithCancel(ctx)
	var output *RabbitMq = &RabbitMq{
		ctx:    rabbitmqCtx,
		cancel: rabbitmqCancel,
	}

	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	output.URL = fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)

	// Establish connection to RabbitMQ
	var err error
	output.Conn, err = amqp.Dial(output.URL)
	if err != nil {
		return nil, err
	}

	channel, err := output.Conn.Channel()
	if err != nil {
		return nil, err
	}
	defer channel.Close()

	// Declare the exchange with type 'topic'
	err = channel.ExchangeDeclare(
		EXCHANGE_NAME_TaskService, // exchange name
		"topic",                   // exchange type (topic)
		true,                      // durable
		false,                     // auto-delete when unused
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// Method to send a message to RabbitMQ
func (r *RabbitMq) SendMessage(exchangeName, action, message string) {
	channel, err := r.Conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

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

func (r *RabbitMq) Close() {
	r.Conn.Close()
	r.cancel()
}
