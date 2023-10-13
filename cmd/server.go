package main

import (
	"context"
	"os"
	"short-link/cmd/rest"
	"short-link/internal"
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
func (s *server) Initialize(cfg *internal.Config, ctx context.Context) error {

	dependencies := CreateDependencies(cfg)

	s.RESTHandler = dependencies.Handler

	return nil
}

// Start starts the application in blocking mode
func (s *server) Start(ctx context.Context) {
	const op = "app.start"

	// Create Router for HTTP Server
	router := SetupRouter(s.RESTHandler)
	// Start GRPC Server in go-routine
	//go s.GRPCHandler.Start(ctx, s.Config.GRPCPort)
	// Start REST Server in Blocking mode
	s.RESTHandler.Start(ctx, router, 8080)
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the app
func (s *server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	const op = "app.gacefulshutdown"
	// Wait for OS signals
	<-quitSignal

	// Kill the API Endpoints first
	s.RESTHandler.Stop()
	//s.GRPCHandler.Stop()

	close(done)
}
