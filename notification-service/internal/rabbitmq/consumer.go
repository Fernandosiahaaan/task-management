package rabbitmq

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type rabbitMq struct {
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
func Init() (*rabbitMq, error) {
	// Create a new instance of rabbitMq
	output := &rabbitMq{}

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

// Method to receive messages from RabbitMQ
func (r *rabbitMq) ReceiveMessage() {
	// Open a channel
	channel, err := r.Conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	// Set up a consumer for the queue
	msgs, err := channel.Consume(
		"testing", // Queue name
		"",        // Consumer
		true,      // Auto acknowledge
		false,     // Exclusive
		false,     // No local
		false,     // No wait
		nil,       // Arguments
	)
	failOnError(err, "Failed to register a consumer")

	// Goroutine to handle incoming messages
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Printf("Received Message: %s\n", msg.Body)
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
