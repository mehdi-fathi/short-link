package main

import (
	"log"
	"short-link/cmd/rest"
	"short-link/internal"
	"short-link/internal/Cache"
	"short-link/internal/Db"
	"short-link/internal/Db/Repository"
)

type out struct {
	Handler *rest.Handler
}

func CreateDependencies(cfg *internal.Config) out {

	// connect to DB first
	var errDb error
	dbLayer := Db.CreateDb(cfg)
	_, errDb = dbLayer.ConnectDB()
	if errDb != nil {
		log.Fatalf("failed to start the server: %v", errDb)
	}

	Repo := Repository.CreateRepository(cfg, dbLayer)

	//httpHandler := &Handler{
	//	Service: internal.CreateService(cfg),
	//}

	client := Cache.CreateCache()

	var ser = internal.CreateService(cfg, Repo, client)

	handler := rest.CreateHandler(ser)

	return out{
		handler,
	}

}
