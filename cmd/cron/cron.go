package cron

import (
	"github.com/go-co-op/gocron"
	service_interface "short-link/internal/interface"
	"short-link/pkg/logger"
)

var cron *gocron.Scheduler

var wait = make(chan interface{})

func StartCron(cron *gocron.Scheduler, service service_interface.ServiceInterface) {

	// 4
	cron.Every(10).Seconds().Tag("test").Do(func() {
		logger.CreateLogInfo("Updating stat links...")
		service.UpdateStats()
	})

	// 5
	cron.StartBlocking()

	//cron.StopBlockingChan()

}
