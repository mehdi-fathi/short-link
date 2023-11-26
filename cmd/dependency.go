package main

import (
	"log"
	"short-link/cmd/rest"
	"short-link/internal"
	"short-link/internal/Cache"
	"short-link/internal/Cache/MemCache"
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

	cache := Cache.CreateCache(cfg)

	memCache := MemCache.CreateMemCache(cfg)

	queue := Queue.CreateQueue(cfg)

	queueMain = queue

	var ser = internal.CreateService(cfg, repo, cache, memCache, queue)

	queue.Service = ser

	handler := rest.CreateHandler(ser)

	return out{
		handler,
	}

}
