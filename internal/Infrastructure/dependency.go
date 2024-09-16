package Infrastructure

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
	"short-link/internal/Core/Ports"
	"short-link/internal/Queue"
	"short-link/pkg/logger"
)

var queueMain *Queue.Queue

type handlerDependencies struct {
	Handler     *Handler
	HandlerRest *rest.HandlerRest
	HandlerWeb  *web.HandlerWeb
}

func CreateHandlerDependencies(cfg *Config.Config) handlerDependencies {

	// connect to DB first
	var errDb error
	dbLayer := Db.CreateDb(cfg)
	_, errDb = dbLayer.ConnectDB()
	if errDb != nil {
		log.Fatalf("failed to start the server: %v", errDb)
	}

	linkRepo := Repository.CreateLinkRepository(cfg, dbLayer)
	shortKeyRepo := Repository.CreateShortKeyRepository(cfg, dbLayer)

	cache := Cache.CreateCache(cfg)

	memCache := MemCache.CreateMemCache(cfg)

	queue := setQueue(cfg)

	var linkService = Service.CreateLinkService(cfg, linkRepo, shortKeyRepo, cache, memCache, queue)

	setServiceForQueue(queue, linkService)

	HandlerRest := CreateHandlerRest(linkService)

	handlerWeb := CreateHandlerWeb(linkService)

	router := SetupRouter(HandlerRest, handlerWeb)

	// Create Router for HTTP server
	HandlerMain := CreateHandlerMain(router, cfg.HTTPPort)

	return handlerDependencies{
		HandlerMain,
		HandlerRest,
		handlerWeb,
	}

}

func createConfigDependency() *Config.Config {
	cfg := Config.LoadConfigEnvApp()
	initLogger(cfg)

	return cfg
}

func setServiceForQueue(queue *Queue.Queue, service Ports.LinkServiceInterface) {
	queue.Service = service
}

func setQueue(cfg *Config.Config) *Queue.Queue {
	queue := Queue.CreateQueue(cfg)

	queueMain = queue
	return queue
}

func getQueue() *Queue.Queue {
	return queueMain
}

func initLogger(cfg *Config.Config) {
	loggerInstance := logger.CreateLogger(cfg.Graylog)
	loggerInstance.Info("[OK] Graylog Configured")
}
