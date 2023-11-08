package main

import (
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
func (s *server) Start() {
	const op = "app.start"

	go cron.StartCron(cronjob, s.RESTHandler.LinkService)

	// Create Router for HTTP Server
	router := SetupRouter(s.RESTHandler)

	// Start GRPC Server in go-routine
	//go s.GRPCHandler.Start(ctx, s.Config.GRPCPort)
	// Start REST Server in Blocking mode
	s.RESTHandler.Start(router, 8080)
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the app
func (s *server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	const op = "app.gacefulshutdown"
	// Wait for OS signals
	<-quitSignal

	// Kill the API Endpoints first
	s.RESTHandler.Stop()
	//s.GRPCHandler.Stop()

	cronjob.StopBlockingChan()

	close(done)

}
