package internal

import (
	"log"
	"math/rand"
	cache_interface "short-link/internal/Cache/Interface"
	Config2 "short-link/internal/Config"
	"short-link/internal/Db/Repository/interface"
	"short-link/internal/Queue"
	service_interface "short-link/internal/interface"
	"strconv"
	"time"
)

type UrlShortener struct {
	Config *Config2.Config
}

type Service struct {
	Shortener *UrlShortener
	LinkRepo  repository_interface.RepositoryInterface
	Cache     cache_interface.CacheInterface
	Queue     *Queue.Queue
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
func CreateService(
	cfg *Config2.Config, linkRepo repository_interface.RepositoryInterface,
	cache cache_interface.CacheInterface, queue *Queue.Queue,
) service_interface.ServiceInterface {

	shortenerUrl := &UrlShortener{
		Config: cfg,
	}

	return &Service{
		Shortener: shortenerUrl,
		LinkRepo:  linkRepo,
		Cache:     cache,
		Queue:     queue,
	}
}

func (service *Service) GetUrl(shortKey string) *repository_interface.Link {

	link, _ := service.LinkRepo.FindByShortKey(shortKey)

	if link.Link != "" {
		service.Cache.IncrBy(shortKey, 1)
	}

	return link
}

func (service *Service) UpdateStats() int {

	all, _ := service.LinkRepo.GetAll()

	for _, data := range all {

		hget, _ := service.Cache.Get(data.ShortKey)

		visitCache, _ := strconv.Atoi(hget)

		if visitCache > data.Visit {
			service.LinkRepo.UpdateVisit(visitCache, data.ShortKey)
			log.Printf("Updated %s : visit :%v", data.ShortKey, visitCache)
		}

		//var linkTable repository_interface.Link
		//
		//err = rows.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey)
		//
		//users[i] = &linkTable

	}
	return 1
}

func (service *Service) SetUrl(link string) string {

	shortKey := GenerateShortKey(service.Shortener.Config.HASHCODE)

	id, err := service.LinkRepo.Create(link, shortKey)

	service.Queue.Publish()

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
