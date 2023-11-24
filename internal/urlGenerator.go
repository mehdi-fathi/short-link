package internal

import (
	"context"
	"fmt"
	"math/rand"
	cache_interface "short-link/internal/Cache/Interface"
	Config2 "short-link/internal/Config"
	"short-link/internal/Db/Model"
	"short-link/internal/Db/Repository/interface"
	"short-link/internal/Event"
	"short-link/internal/Queue"
	service_interface "short-link/internal/interface"
	"short-link/pkg/logger"
	"short-link/pkg/url"
	"strconv"
	"sync"
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

func (service *Service) GetUrl(shortKey string) *Model.Link {

	link, _ := service.LinkRepo.FindByShortKey(shortKey)

	if link != nil && link.Link != "" {
		service.Cache.IncrBy(shortKey, 1)
	}

	return link
}

func (service *Service) UpdateStats(s *sync.WaitGroup, ctx context.Context) int {

	var all map[int]*Model.Link

	limit := 10

	start := 0

	// Setting up the main context
	//ctx1 := context.Background()

	var counter int = 1

	for (all[0] != nil && counter > 1) || counter == 1 {

		if counter > 1 {
			start = counter * limit
		}

		all, _ = service.LinkRepo.GetChunk(start, limit, "approve")

		if all[0] != nil {

			s.Add(1)

			go func(start int, all map[int]*Model.Link) {

				defer s.Done()

				// Create a new context with a timeout for the processing
				_, cancelProc := context.WithTimeout(ctx, 5*time.Second)

				defer cancelProc()

				//logger.CreateLogInfo(fmt.Sprintf("Run Go routine %d", start))

				for _, data := range all {

					hget, _ := service.Cache.Get(data.ShortKey)

					//logger.CreateLogInfo(fmt.Sprintf("Run %s ", data.ShortKey))

					visitCache, _ := strconv.Atoi(hget)

					if visitCache > data.Visit {
						service.LinkRepo.UpdateVisit(visitCache, data.ShortKey)
						logger.CreateLogInfo(fmt.Sprintf("Updated %s : visit :%v", data.ShortKey, visitCache))
					}

					status := "approve"

					if !url.CheckURL(data.Link) {
						logger.CreateLogInfo(fmt.Sprintf("Not approved ShortKey :%v", data.ShortKey))
						status = "reject"
					}

					service.LinkRepo.UpdateStatus(status, data.Link)

				}
				//logger.CreateLogInfo(fmt.Sprintf("Finish Go routine  %d", start))

			}(start, all)

		}

		counter++
	}

	// Listen for the context cancellation signal to stop the cron scheduler
	<-ctx.Done()
	logger.CreateLogInfo("[ * ] Received shutdown signal, UpdateStats...")
	//
	////todo we can make a goroutine version for this
	//
	//for _, data := range all {
	//
	//	hget, _ := service.Cache.Get(data.ShortKey)
	//
	//	visitCache, _ := strconv.Atoi(hget)
	//
	//	if visitCache > data.Visit {
	//		service.LinkRepo.UpdateVisit(visitCache, data.ShortKey)
	//		logger.CreateLogInfo(fmt.Sprintf("Updated %s : visit :%v", data.ShortKey, visitCache))
	//	}
	//
	//	//var linkTable Model.Link
	//	//
	//	//err = rows.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey)
	//	//
	//	//users[i] = &linkTable
	//
	//}
	return 1
}

func (service *Service) UpdateStatus(status string, shortKey string) {

	service.LinkRepo.UpdateStatus(status, shortKey)

}

func (service *Service) checkPendingLinks() int {

	all, _ := service.LinkRepo.GetByStatus("pending")

	var status string
	for _, data := range all {
		status = "approve"
		logger.CreateLogInfo(data.Link)

		if !url.CheckURL(data.Link) {
			logger.CreateLogInfo(fmt.Sprintf("Not approved ShortKey :%v", data.ShortKey))
			status = "reject"
		}

		service.LinkRepo.UpdateStatus(status, data.ShortKey)

	}
	return 1
}

func (service *Service) SetUrl(link string) string {

	shortKey := GenerateShortKey(service.Shortener.Config.HASHCODE)

	_, err := service.LinkRepo.Create(link, shortKey)

	var data = make(map[string]string)

	data["link"] = link

	// Emit an event
	event := Event.Event{Type: "OrderPlaced", Data: data}

	ch, err := service.Queue.Connection.Channel()

	//service.checkPendingLinks()

	service.Queue.Publish(ch, "test", event)

	if err != nil {
		return ""
	}

	return shortKey
}

//func (service *Service) GetAllUrl() map[string]string {
//	//return service.Shortener.Urls
//	return service.LinkRepo.GetAll()
//}

func (service *Service) GetAllUrlV2() (map[int]*Model.Link, error) {
	//return service.Shortener.Urls
	return service.LinkRepo.GetAll()
}
