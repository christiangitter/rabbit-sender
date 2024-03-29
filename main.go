package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a queue, which will create it if it doesn't exist
	queue_name := "test_queue"
	_, err = ch.QueueDeclare(
		queue_name, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")
	log.Printf("Queue '%s' is ensured to exist.", queue_name)

	// Create a message in JSON format
	message := map[string]interface{}{
		"hello": "world",
		"foo":   "bar",
		"baz":   123,
		"qux":   []string{"a", "b", "c"},
	}
	body, err := json.Marshal(message)
	failOnError(err, "Failed to marshal JSON")

	err = ch.Publish(
		"",         // exchange
		queue_name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
