package cron

import (
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

func StartCron(cron *gocron.Scheduler, service service_interface.ServiceInterface) {

	// 4
	cron.Every(10).Seconds().Tag("test").Do(func() {
		hello("Updating stat links...")
		service.UpdateStats()
	})

	// 5
	cron.StartBlocking()

	//cron.StopBlockingChan()

}
