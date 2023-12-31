package Http

import (
	"log"
	"short-link/internal/Cache"
	"short-link/internal/Cache/MemCache"
	"short-link/internal/Config"
	"short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Handlers/Http/web"
	"short-link/internal/Core/Logic/Db"
	"short-link/internal/Core/Logic/Db/Repository"
	"short-link/internal/Core/Logic/Service"
	"short-link/internal/Queue"
)

var queueMain *Queue.Queue

type out struct {
	Handler     *Handler
	HandlerRest *rest.HandlerRest
	HandlerWeb  *web.HandlerWeb
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

	var ser = Service.CreateService(cfg, repo, cache, memCache, queue)

	queue.Service = ser

	HandlerRest := CreateHandler(ser)
	HandlerMain := CreateHandlerMain()
	handlerWeb := CreateHandlerWeb(ser)

	log.Println("run")

	return out{
		HandlerMain,
		HandlerRest,
		handlerWeb,
	}

}
