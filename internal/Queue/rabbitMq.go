package Queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"short-link/internal/Config"
	"short-link/internal/Db/Model"
	"short-link/internal/Event"
	service_interface "short-link/internal/interface"
	"short-link/pkg/logger"
	"short-link/pkg/url"
	"time"
)

type Queue struct {
	Connection *amqp.Connection
	cfg        *Config.Config
	Service    service_interface.ServiceInterface
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
		logger.CreateLogError(fmt.Sprintf("failed to connect to RabbitMQ: %v", err))
		panic(err.(interface{}))
	}

	logger.CreateLogInfo("Successfully connected to RabbitMQ instance")

	return connection
}

func (qu *Queue) Publish(ch *amqp.Channel, queueName string, event Event.Event) {

	var err error
	// opening a channel over the connection established to interact with RabbitMQ
	channel := ch
	if err != nil {
		panic(err.(interface{}))
	}
	defer channel.Close()

	// declaring queue with its properties over the the channel opened
	queue, err := channel.QueueDeclare(
		qu.cfg.QueueRabbit.MainQueueName, // name
		true,                             // durable
		false,                            // auto delete
		false,                            // exclusive
		false,                            // no wait
		nil,                              // args
	)
	if err != nil {
		panic(err.(interface{}))
	}

	body, _ := json.Marshal(event)

	// publishing a message
	err = channel.Publish(
		"",                               // exchange
		qu.cfg.QueueRabbit.MainQueueName, // key
		false,                            // mandatory
		false,                            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
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
		Service:    nil,
	}

	return queue
}

// ConsumeEvents listens for messages on a RabbitMQ queue and processes them
func (qu *Queue) ConsumeEvents(ctx context.Context, ch *amqp.Channel, queueName string) {

	logger.CreateLogInfo(" [*] Waiting for events")

	ch, err := qu.Connection.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	for {
		select {
		case <-ctx.Done(): // if cancel() is called
			logger.CreateLogInfo("Consumer received shutdown signal")
			return
		case msg, ok := <-msgs:
			if !ok {
				logger.CreateLogInfo("Channel closed, consumer exiting")
				return
			}

			logger.CreateLogInfo("Received a message")

			var event Event.Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				logger.CreateLogError(fmt.Sprintf("Error decoding event: %s", err))
				msg.Nack(false, true) // negative acknowledgment, requeue the message
				continue
			}
			if len(msg.Body) == 0 {
				logger.CreateLogInfo("Received an empty message, skipping...")
				msg.Nack(false, true) // negative acknowledgment, requeue the message
				continue
			}
			// Create a new context with a timeout for the processing
			procCtx, cancelProc := context.WithTimeout(ctx, 10*time.Second)
			defer cancelProc()

			// Process the event with its own context
			// Replace `ProcessEvent` with actual event processing logic
			if err := qu.ProcessEvent(procCtx, event); err != nil {
				logger.CreateLogError(fmt.Sprintf("Failed to process event: %s, error: %v", event.Type, err))
				msg.Nack(false, true) // negative acknowledgment, requeue the message
			} else {
				logger.CreateLogInfo(fmt.Sprintf("Event processed successfully: %s", event.Type))
				msg.Ack(false) // Acknowledge the message after successful processing
			}
		}
	}

	//<-forever
}

// todo make some event listener
// ProcessEvent simulates event processing.
func (qu *Queue) ProcessEvent(ctx context.Context, event Event.Event) error {
	// Simulate work
	select {
	case <-time.After(1 * time.Second):
		//time.Sleep(5 * time.Second)

		data := event.Data.(map[string]interface{})

		status := Model.LINK_STATUS_APPROVE
		if !url.CheckURL(data["link"].(string)) {
			logger.CreateLogInfo(fmt.Sprintf("Rejected link :%s", data["link"].(string)))
			status = Model.Link_STATUS_REJECT
		}

		qu.Service.UpdateStatus(status, data["link"].(string))

		logger.CreateLogInfo(fmt.Sprintf("Event processed Done: %s with status: %s", event.Type, status))
		return nil
	case <-ctx.Done():
		logger.CreateLogInfo(" Cancel queue")

		return ctx.Err()
	}
}
