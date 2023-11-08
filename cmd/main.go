package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/signal"
	"short-link/internal/Config"
	"short-link/pkg/logger"
	"syscall"
	"time"
)

/*
go mod tidy ensures that the go.mod file matches the source code in the module.
It adds any missing module requirements necessary to build the current module’s packages and dependencies,
if there are some not used dependencies go mod tidy will remove those from go.mod accordingly
*/

func main() {

	startTime := time.Now()

	// Default Config file based on the environment variable
	defaultConfigFile := "config/config-local.yaml"
	if env := os.Getenv("APP_MODE"); env != "" {
		defaultConfigFile = fmt.Sprintf("config/config-%s.yaml", env)
	}

	// Load Master Config File
	var configFile string
	flag.StringVar(&configFile, "config", defaultConfigFile, "The environment configuration file of application")
	flag.Usage = usage
	flag.Parse()

	// Loading the config file
	cfg, err := Config.LoadConfig(configFile)
	if err != nil {
		log.Println(errors.Wrapf(err, "failed to load config: %s", "CreateService"))
	}

	loggerInstance := logger.CreateLogger(cfg.Logger)
	loggerInstance.Info("[OK] Logger Configured")

	//loggerInstance := logrus.Logger{}
	//loggerInstance.Info("[OK] Logger Configured")

	// Setting up the main context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Create New Server
	server := NewServer(startTime)

	// Initialize the Server Dependencies
	err = server.Initialize(cfg)

	done := make(chan bool, 1)
	quiteSignal := make(chan os.Signal, 1)
	signal.Notify(quiteSignal, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown goroutine
	go server.GracefulShutdown(quiteSignal, done)

	// Start server in blocking mode
	server.Start()

	if err != nil {
		log.Fatal(errors.Wrap(err, "server error"))
	}

	// Wait for graceful shutdown signal
	<-done

	// Kill other background jobs
	cancel()
	loggerInstance.Info("Waiting for background jobs to finish their works...")

	// Wait for all other background jobs to finish their works
	server.Wait()

	loggerInstance.Info("Master App Shutdown successfully, see you next time ;-)")

}

func usage() {
	usageStr := `
Usage: server [options]
Options:
	-c,  --config   <config file name>   Path of yaml configuration file
`
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
