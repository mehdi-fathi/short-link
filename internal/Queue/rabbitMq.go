package Queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"log"
	"short-link/internal/Config"
	"short-link/internal/Event"
	"short-link/pkg/logger"
	"time"
)

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

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	queueName := "sync"
	q, err := ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return q, errors.Wrap(err, "failed to declare queue")
	}
	return q, nil
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

// ConsumeEvents listens for messages on a RabbitMQ queue and processes them
func (qu *Queue) ConsumeEvents(ctx context.Context, queueName string) {

	logger.CreateLogInfo(" [*] Queue is Waiting for  events")

	ch, err := qu.openChannel()
	defer ch.Close()

	err = qu.declareQueue(err, ch)

	msgs := qu.consumeChannel(err, ch, queueName)

	qu.listenEvents(ctx, msgs)

	//<-forever
}

func (qu *Queue) listenEvents(ctx context.Context, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done(): // if cancel() is called
			logger.CreateLogInfo(" [*] Queue Consumer received shutdown signal")
			return
		case msg, ok := <-msgs:
			if !ok {
				logger.CreateLogInfo(" [*] Queue Channel closed, consumer exiting")
				return
			}

			logger.CreateLogInfo(" [*] Queue Received a message")

			event, err := unmarshalMsgEvent(msg)
			if err != nil {
				continue
			}

			// Create a new context with a timeout for the processing
			procCtx, cancelProc := context.WithTimeout(ctx, 10*time.Second)
			defer cancelProc()

			qu.processEventWrapper(procCtx, event, msg)
		}
	}
}

func unmarshalMsgEvent(msg amqp.Delivery) (Event.Event, error) {
	var event Event.Event
	err := json.Unmarshal(msg.Body, &event)

	if len(msg.Body) == 0 {
		err = errors.New("Event Msg is empty")
	}

	if err != nil {
		logger.CreateLogError(fmt.Sprintf("Error decoding event: %s", err))
		msg.Nack(false, true) // negative acknowledgment, requeue the message
	}

	return event, err
}

func (qu *Queue) processEventWrapper(procCtx context.Context, event Event.Event, msg amqp.Delivery) {
	// Process the event with its own context
	// Replace `ProcessEvent` with actual event processing Logic
	if err := qu.ProcessEvent(procCtx, event); err != nil {
		logger.CreateLogError(fmt.Sprintf("[*] Queue Failed to process event: %s, error: %v", event.Type, err))
		msg.Nack(false, true) // negative acknowledgment, requeue the message
	} else {
		logger.CreateLogInfo(fmt.Sprintf("[*] Queue Event processed successfully: %s", event.Type))
		msg.Ack(false) // Acknowledge the message after successful processing
	}
}

func (qu *Queue) consumeChannel(err error, ch *amqp.Channel, queueName string) <-chan amqp.Delivery {
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
	return msgs
}

func (qu *Queue) declareQueue(err error, ch *amqp.Channel) error {
	// Declare the queue
	_, err = declareQueue(ch)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
	return err
}

func (qu *Queue) openChannel() (*amqp.Channel, error) {
	ch, err := qu.Connection.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	return ch, err
}

func (qu *Queue) consume(queueName string) <-chan amqp.Delivery {

	ch, err := qu.Connection.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare the queue
	_, err = declareQueue(ch)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}
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
	return msgs
}

// ProcessEvent simulates event processing.
func (qu *Queue) ProcessEvent(ctx context.Context, event Event.Event) error {
	// Simulate work
	select {
	case <-time.After(1 * time.Second):
		//time.Sleep(5 * time.Second)
		data := event.Data.(map[string]interface{})
		eventType := event.Type

		if eventType == Event.CreateLink {
			status := qu.Service.VerifyLinkIsValid(data["link"].(string))
			logger.CreateLogInfo(fmt.Sprintf("[*] Queue Event processed Done: %s with status: %s", event.Type, status))
		}
		return nil
	case <-ctx.Done():
		logger.CreateLogInfo("[*] Queue shutdown")

		return ctx.Err()
	}
}
