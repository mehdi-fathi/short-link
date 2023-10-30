package Queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type Queue struct {
	Connection *amqp.Connection
}

func CreateConnection() *Queue {

	fmt.Println("RabbitMQ in Golang: Getting started tutorial")

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println(err)
		panic(err)
	}

	queue := &Queue{
		Connection: connection,
	}

	return queue
}

func (qu *Queue) Publish() {

	defer qu.Connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	// opening a channel over the connection established to interact with RabbitMQ
	channel, err := qu.Connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// declaring queue with its properties over the the channel opened
	queue, err := channel.QueueDeclare(
		"testing", // name
		false,     // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		panic(err)
	}

	// publishing a message
	err = channel.Publish(
		"",        // exchange
		"testing", // key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Test Message"),
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Queue status:", queue)
	fmt.Println("Successfully published message")
}
