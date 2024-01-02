package Cron

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	service_interface "short-link/internal/Core/Ports"
	"short-link/pkg/logger"
	"sync"
	"time"
)

var cron *gocron.Scheduler

var wait = make(chan interface{})

func StartCron(ctx context.Context, cron *gocron.Scheduler, service service_interface.ServiceInterface) {

	cron.Every(10).Seconds().Tag("cron").Do(func() {

		logger.CreateLogInfo("[*] Cron start...")

		// Process the event with its own context
		// Replace `ProcessCron` with actual event processing Logic
		if err := ProcessCron(service, ctx); err != nil {
			logger.CreateLogError(fmt.Sprintf("[*] Cron Failed to process cron"))
		} else {
			logger.CreateLogInfo(fmt.Sprintf("[*] Cron processed successfully"))
		}

	})

	// 5
	//cron.StartBlocking()

	// Start the scheduler
	cron.StartAsync()

	// Listen for the context cancellation signal to stop the cron scheduler
	<-ctx.Done()
	logger.CreateLogInfo("[*] Cron -- Received shutdown signal, stopping cron scheduler...")
	cron.Stop() // Stop the scheduler

	//cron.StopBlockingChan()

}

func ProcessCron(service service_interface.ServiceInterface, ctx context.Context) error {

	var cronWaitGroup sync.WaitGroup

	// Simulate work
	select {
	case <-time.After(1 * time.Second):

		service.UpdateStats(&cronWaitGroup, ctx)
		//time.Sleep(15 * time.Second)
		cronWaitGroup.Wait()
		return nil
	case <-ctx.Done():
		logger.CreateLogInfo(fmt.Sprintf("[*] Cron Shutdown"))
		return ctx.Err()
	}
}
