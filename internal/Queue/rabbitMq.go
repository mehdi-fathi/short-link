package Queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"short-link/internal/Config"
	"short-link/pkg/logger"
)

type Queue struct {
	Connection *amqp.Connection
	cfg        *Config.Config
}

func CreateConnection(cfg *Config.Config) *amqp.Connection {

	connString := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		cfg.QueueRabbit.User,
		cfg.QueueRabbit.Password,
		cfg.QueueRabbit.Host,
		cfg.QueueRabbit.Port,
	)

	connection, err := amqp.Dial(connString)

	if err != nil {
		logger.CreateLogError(fmt.Sprintf("failed to connect to database: %v", err))
		panic(err.(interface{}))
	}

	logger.CreateLogInfo("Successfully connected to RabbitMQ instance")

	return connection
}

func (qu *Queue) Publish() {

	defer qu.Connection.Close()

	var err error
	// opening a channel over the connection established to interact with RabbitMQ
	channel, err := qu.Connection.Channel()
	if err != nil {
		panic(err.(interface{}))
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
		panic(err.(interface{}))
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
		panic(err.(interface{}))
	}

	logger.CreateLogInfo(fmt.Sprintf("Successfully published message - Queue status: %v", queue))
}

func CreateQueue(cfg *Config.Config) *Queue {

	queue := &Queue{
		Connection: CreateConnection(cfg),
		cfg:        cfg,
	}

	return queue
}
