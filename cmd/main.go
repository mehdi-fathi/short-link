package main

import (
	"github.com/pkg/errors"
	"log"
	"short-link/internal/Infrastructure"
	"time"
)

/*
go mod tidy ensures that the go.mod file matches the source code in the module.
It adds any missing module requirements necessary to build the current module’s packages and dependencies,
if there are some not used dependencies go mod tidy will remove those from go.mod accordingly
*/
var startTime time.Time

func init() {

	startTime = time.Now()
	log.Println("[OK] Running App...")
}

func main() {

	//loggerInstance := logrus.Graylog{}
	//loggerInstance.Info("[OK] Graylog Configured")

	// Create New server
	server := Infrastructure.NewServer(startTime)
	// StartApp the server Dependencies
	err := server.StartApp()

	if err != nil {
		log.Fatal(errors.Wrap(err, "server error"))
	}

}
