package cron

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	service_interface "short-link/internal/interface"
	"short-link/pkg/logger"
	"time"
)

var cron *gocron.Scheduler

var wait = make(chan interface{})

func StartCron(ctx context.Context, cron *gocron.Scheduler, service service_interface.ServiceInterface) {

	// 4
	cron.Every(10).Seconds().Tag("test").Do(func() {

		// Create a new context with a timeout for the processing
		procCtx, cancelProc := context.WithTimeout(ctx, 10*time.Second)

		defer cancelProc()

		logger.CreateLogInfo("Updating stat links...")

		// Process the event with its own context
		// Replace `ProcessEvent` with actual event processing logic
		if err := ProcessEvent(service, procCtx); err != nil {
			logger.CreateLogError(fmt.Sprintf("Failed to process cron:"))
		} else {
			logger.CreateLogInfo(fmt.Sprintf("Cron processed successfully"))
		}

	})

	// 5
	//cron.StartBlocking()

	// Start the scheduler
	cron.StartAsync()

	// Listen for the context cancellation signal to stop the cron scheduler
	<-ctx.Done()
	logger.CreateLogInfo("Received shutdown signal, stopping cron scheduler...")
	cron.Stop() // Stop the scheduler

	//cron.StopBlockingChan()

}

// ProcessEvent simulates event processing.
func ProcessEvent(service service_interface.ServiceInterface, ctx context.Context) error {
	// Simulate work
	select {
	case <-time.After(1 * time.Second):
		service.UpdateStats()
		time.Sleep(5 * time.Second)
		logger.CreateLogInfo(fmt.Sprintf("Cron done"))
		return nil
	case <-ctx.Done():
		logger.CreateLogInfo(" Cancel Cron")

		return ctx.Err()
	}
}
