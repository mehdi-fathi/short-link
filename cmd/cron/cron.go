package cron

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	service_interface "short-link/internal/interface"
)

var cron *gocron.Scheduler

var wait = make(chan interface{})

// 2
func hello(name string) {
	message := fmt.Sprintf("Cron runned : %v", name)
	fmt.Println(message)
}

func StartCron(cron *gocron.Scheduler, service service_interface.ServiceInterface, ctx context.Context) {

	// 4
	cron.Every(2).Seconds().Tag("test").Do(func() {
		hello("Saving stat links")
		service.UpdateStats()
	})

	// 5
	cron.StartBlocking()

	//cron.StopBlockingChan()

}
