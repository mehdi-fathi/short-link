package Service

import (
	"short-link/internal/Config"
	"short-link/internal/Core/Ports"
	"short-link/internal/Queue"
)

type LinkService struct {
	Config       *Config.Config
	LinkRepo     Ports.LinkRepositoryInterface
	ShortKeyRepo Ports.ShortKeyRepositoryInterface
	Cache        Ports.CacheInterface
	MemCache     Ports.MemCacheInterface
	Queue        *Queue.Queue
}

// CreateLinkService creates an instance of membership interface with the necessary dependencies
func CreateLinkService(
	cfg *Config.Config,
	linkRepo Ports.LinkRepositoryInterface,
	shortKeyRepo Ports.ShortKeyRepositoryInterface,
	cache Ports.CacheInterface,
	memCache Ports.MemCacheInterface,
	queue *Queue.Queue,
) Ports.LinkServiceInterface {

	return &LinkService{
		Config:       cfg,
		LinkRepo:     linkRepo,
		ShortKeyRepo: shortKeyRepo,
		Cache:        cache,
		MemCache:     memCache,
		Queue:        queue,
	}
}
