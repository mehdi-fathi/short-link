package Infrastructure

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"short-link/internal/Config"
	"short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Handlers/Http/web"
	"short-link/internal/Core/Logic/Service"
	"short-link/internal/Cron"
	"short-link/pkg/logger"
	"sync"
	"syscall"
	"time"
)

type server struct {
	sync.WaitGroup
	StartTime   time.Time
	Handler     *Handler
	Config      *Config.Config
	RESTHandler *rest.HandlerRest
	WebHandler  *web.HandlerWeb
}

var cronjob *gocron.Scheduler

var done = make(chan bool, 1)

// NewServer Create a new instance of server application
func NewServer(startTime time.Time) *server {
	return &server{
		StartTime: startTime,
	}
}

func init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(Service.GoroutinesCount)
}

// StartApp is responsible for app initialization and wrapping required handlerDependencies
func (s *server) StartApp() error {

	s.injectServerDependencies()

	s.shutdownListener()

	// Start server in blocking mode
	s.Start()

	return nil
}

func (s *server) injectServerDependencies() {

	//cfg := Graylog.LoadConfigApp()

	cfg := createConfigDependency()

	dependencies := CreateHandlerDependencies(cfg)

	s.Handler = dependencies.Handler
	s.Config = cfg
	s.RESTHandler = dependencies.HandlerRest
	s.WebHandler = dependencies.HandlerWeb
}

func (s *server) shutdownListener() {
	quiteSignal := make(chan os.Signal, 1)
	signal.Notify(quiteSignal, syscall.SIGINT, syscall.SIGTERM)

	// Use a WaitGroup to wait for goroutines to finish
	s.Add(1)
	// Graceful shutdown goroutine
	go s.GracefulShutdown(quiteSignal)
}

// Start starts the application in blocking mode
func (s *server) Start() {

	const op = "app.start"

	ctx, cancel := s.buildContext()

	s.startCronJob(ctx)

	s.startListenEvents(ctx)

	s.startGrafanaDashboard()

	s.mainStartHttp(cancel)

}

func (s *server) startGrafanaDashboard() {
	go func() {

		http.Handle("/metrics", promhttp.Handler())
		path := fmt.Sprintf("0.0.0.0:%d", s.Config.GRAFANAPort)

		http.ListenAndServe(path, nil)

	}()
}

func (s *server) buildContext() (context.Context, context.CancelFunc) {
	// Setting up the main context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return ctx, cancel
}

func (s *server) mainStartHttp(cancel context.CancelFunc) {
	s.startHttp()
	s.listenShutdownHttp(cancel)
}

func (s *server) startHttp() {

	// Start GRPC server in go-routine
	//go s.GRPCHandler.Start(ctx, s.Graylog.GRPCPort)
	// Start REST server in Blocking mode
	s.Handler.Start()

}

func (s *server) listenShutdownHttp(cancel context.CancelFunc) {
	// Wait for graceful shutdown signal
	<-done

	// Kill other background jobs
	cancel()
	logger.CreateLogInfo("Waiting for background jobs to finish their works...")

	// Wait for all other background jobs to finish their works
	s.Wait()

	logger.CreateLogInfo("Master App Shutdown successfully, see you next time ;-)")
}

func (s *server) startListenEvents(ctx context.Context) {
	// Use a WaitGroup to wait for goroutines to finish
	s.Add(1)

	Service.GoroutinesCount.WithLabelValues("startListenEvents").Inc()
	go func() {
		defer s.Done()
		// Start the consumer
		getQueue().ConsumeEvents(ctx, s.Config.QueueRabbit.MainQueueName)
	}()
}

func (s *server) startCronJob(ctx context.Context) {
	cronjob = gocron.NewScheduler(time.UTC)

	// Use a WaitGroup to wait for goroutines to finish
	s.Add(1)
	logger.CreateLogInfo("startCronJob ...")

	Service.GoroutinesCount.WithLabelValues("startCronJob").Inc()
	go func() {
		defer s.Done()
		Cron.StartCron(ctx, cronjob, s.RESTHandler.LinkService)
	}()
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the app
func (s *server) GracefulShutdown(quitSignal <-chan os.Signal) {
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
