package Http

import (
	"context"
	"github.com/go-co-op/gocron"
	"os"
	"short-link/internal/Config"
	"short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Handlers/Http/web"
	"short-link/internal/Cron"
	"sync"
	"time"
)

type server struct {
	sync.WaitGroup
	StartTime   time.Time
	Handler     *Handler
	RESTHandler *rest.HandlerRest
	WebHandler  *web.HandlerWeb
}

var cronjob *gocron.Scheduler

// NewServer Create a new instance of server application
func NewServer(startTime time.Time) *server {
	return &server{
		StartTime: startTime,
	}
}

// Initialize is responsible for app initialization and wrapping required dependencies
func (s *server) Initialize(cfg *Config.Config) error {

	dependencies := CreateDependencies(cfg)

	s.Handler = dependencies.Handler
	s.RESTHandler = dependencies.HandlerRest
	s.WebHandler = dependencies.HandlerWeb

	cronjob = gocron.NewScheduler(time.UTC)

	return nil
}

// Start starts the application in blocking mode
func (s *server) Start(ctx context.Context, cfg *Config.Config) {
	const op = "app.start"

	// Use a WaitGroup to wait for goroutines to finish
	s.Add(1)

	go func() {
		defer s.Done()
		Cron.StartCron(ctx, cronjob, s.RESTHandler.LinkService)
	}()

	// Use a WaitGroup to wait for goroutines to finish
	s.Add(1)

	go func() {
		defer s.Done()
		ch, _ := queueMain.Connection.Channel()
		// Start the consumer
		queueMain.ConsumeEvents(ctx, ch, cfg.QueueRabbit.MainQueueName)
	}()

	// Create Router for HTTP server
	router := SetupRouter(s.RESTHandler, s.WebHandler)

	// Start GRPC server in go-routine
	//go s.GRPCHandler.Start(ctx, s.Config.GRPCPort)
	// Start REST server in Blocking mode
	s.Handler.Start(router, cfg.HTTPPort)
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the app
func (s *server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	defer s.Done()

	const op = "app.gacefulshutdown"
	// Wait for OS signals
	<-quitSignal

	// Kill the API Endpoints first
	s.Handler.Stop()
	//s.GRPCHandler.Stop()

	cronjob.StopBlockingChan()

	close(done)

}
