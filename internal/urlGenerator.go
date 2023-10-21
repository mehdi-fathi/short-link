package internal

import (
	"log"
	"math/rand"
	cache_interface "short-link/internal/Cache/Interface"
	"short-link/internal/Db/Repository/interface"
	service_interface "short-link/internal/interface"
	"time"
)

type UrlShortener struct {
	Urls   map[string]string
	Config *Config
}

type Service struct {
	Shortener *UrlShortener
	LinkRepo  repository_interface.RepositoryInterface
	Cache     cache_interface.CacheInterface
}

func GenerateShortKey(hashCode string) string {
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = hashCode[rand.Intn(len(hashCode))]
	}
	return string(shortKey)
}

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateService(cfg *Config, linkRepo repository_interface.RepositoryInterface, cache cache_interface.CacheInterface) service_interface.ServiceInterface {

	shortenerUrl := &UrlShortener{
		Urls:   make(map[string]string),
		Config: cfg,
	}

	return &Service{
		Shortener: shortenerUrl,
		LinkRepo:  linkRepo,
		Cache:     cache,
	}
}

func (service *Service) GetUrl(shortKey string) *repository_interface.Link {

	c, _ := service.LinkRepo.FindByShortKey(shortKey)
	log.Println("links : ", c)

	service.Cache.IncrBy(shortKey, 1)

	hget,_ := service.Cache.Get(shortKey)

	log.Println("run",hget)

	return c
}

func (service *Service) SetUrl(link string) string {

	shortKey := GenerateShortKey(service.Shortener.Config.HASHCODE)
	service.Shortener.Urls[shortKey] = link

	id, err := service.LinkRepo.Create(link, shortKey)

	log.Println(id)
	if err != nil {
		return ""
	}

	return shortKey
}

//func (service *Service) GetAllUrl() map[string]string {
//	//return service.Shortener.Urls
//	return service.LinkRepo.GetAll()
//}

func (service *Service) GetAllUrlV2() (map[int]*repository_interface.Link, error) {
	//return service.Shortener.Urls
	return service.LinkRepo.GetAll()
}
