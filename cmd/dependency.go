package main

import (
	"log"
	"short-link/cmd/rest"
	"short-link/internal"
	"short-link/internal/Cache"
	"short-link/internal/Config"
	"short-link/internal/Db"
	"short-link/internal/Db/Repository"
	"short-link/internal/Queue"
)

var queueMain *Queue.Queue

type out struct {
	Handler *rest.Handler
}

func CreateDependencies(cfg *Config.Config) out {

	// connect to DB first
	var errDb error
	dbLayer := Db.CreateDb(cfg)
	_, errDb = dbLayer.ConnectDB()
	if errDb != nil {
		log.Fatalf("failed to start the server: %v", errDb)
	}

	repo := Repository.CreateRepository(cfg, dbLayer)

	//httpHandler := &Handler{
	//	Service: internal.CreateService(cfg),
	//}

	client := Cache.CreateCache(cfg)

	queue := Queue.CreateQueue(cfg)

	queueMain = queue

	var ser = internal.CreateService(cfg, repo, client, queue)

	handler := rest.CreateHandler(ser)

	return out{
		handler,
	}

}
