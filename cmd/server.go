package main

import (
	"context"
	"github.com/go-co-op/gocron"
	"os"
	"short-link/cmd/cron"
	"short-link/cmd/rest"
	"short-link/internal/Config"
	"sync"
	"time"
)

type server struct {
	sync.WaitGroup
	StartTime   time.Time
	RESTHandler *rest.Handler
}

// NewServer Create a new instance of server application
func NewServer(startTime time.Time) *server {
	return &server{
		StartTime: startTime,
	}
}

// Initialize is responsible for app initialization and wrapping required dependencies
func (s *server) Initialize(cfg *Config.Config) error {

	dependencies := CreateDependencies(cfg)

	s.RESTHandler = dependencies.Handler

	cronjob = gocron.NewScheduler(time.UTC)

	return nil
}

var cronjob *gocron.Scheduler

// Start starts the application in blocking mode
func (s *server) Start(ctx context.Context, cfg *Config.Config) {
	const op = "app.start"

	// Use a WaitGroup to wait for goroutines to finish
	s.Add(1)

	go func() {
		defer s.Done()
		cron.StartCron(ctx, cronjob, s.RESTHandler.LinkService)
	}()

	// Use a WaitGroup to wait for goroutines to finish
	s.Add(1)

	go func() {
		defer s.Done()
		ch, _ := queueMain.Connection.Channel()
		// Start the consumer
		queueMain.ConsumeEvents(ctx, ch, cfg.QueueRabbit.MainQueueName)
	}()

	// Emit an event
	//event := Event.Event{Type: "OrderPlaced", Data: "Order123"}
	//if err := Event.EmitEvent(ch, q.Name, event); err != nil {
	//	log.Fatalf("Failed to emit event: %s", err)
	//}

	//forever := make(chan bool)
	//<-forever

	// Create Router for HTTP server
	router := SetupRouter(s.RESTHandler)

	// Start GRPC server in go-routine
	//go s.GRPCHandler.Start(ctx, s.Config.GRPCPort)
	// Start REST server in Blocking mode
	s.RESTHandler.Start(router, 8080)
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the app
func (s *server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	defer s.Done()

	const op = "app.gacefulshutdown"
	// Wait for OS signals
	<-quitSignal

	// Kill the API Endpoints first
	s.RESTHandler.Stop()
	//s.GRPCHandler.Stop()

	cronjob.StopBlockingChan()

	close(done)

}
