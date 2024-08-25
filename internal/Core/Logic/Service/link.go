package Service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"short-link/internal/Config"
	"short-link/internal/Core/Domin"
	"short-link/internal/Core/Ports"
	"short-link/internal/Event"
	"short-link/internal/Queue"
	"short-link/pkg/logger"
	"short-link/pkg/url"
	"strconv"
	"sync"
	"time"
)

type UrlShortener struct {
	Config *Config.Config
}

type Service struct {
	Shortener    *UrlShortener
	LinkRepo     Ports.LinkRepositoryInterface
	ShortKeyRepo Ports.ShortKeyRepositoryInterface
	Cache        Ports.CacheInterface
	MemCache     Ports.MemCacheInterface
	Queue        *Queue.Queue
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func (service *Service) IntToBase62(num int) string {
	if num == 0 {
		return string(charset[0])
	}

	result := ""
	base := len(charset)

	for num > 0 {
		result = string(charset[num%base]) + result
		num = num / base

	}

	// Pad with 'A' to ensure a fixed length of 6 characters
	for len(result) < 6 {
		result = string(charset[0]) + result

	}

	return result
}

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateService(
	cfg *Config.Config,
	linkRepo Ports.LinkRepositoryInterface,
	shortKeyRepo Ports.ShortKeyRepositoryInterface,
	cache Ports.CacheInterface,
	memCache Ports.MemCacheInterface,
	queue *Queue.Queue,
) Ports.LinkServiceInterface {

	shortenerUrl := &UrlShortener{
		Config: cfg,
	}

	return &Service{
		Shortener:    shortenerUrl,
		LinkRepo:     linkRepo,
		ShortKeyRepo: shortKeyRepo,
		Cache:        cache,
		MemCache:     memCache,
		Queue:        queue,
	}
}

func (service *Service) GetUrl(shortKey string) *Domin.Link {

	link, _ := service.LinkRepo.FindByShortKey(shortKey)

	if link != nil && link.Link != "" {
		service.Cache.IncrBy(shortKey, 1)
	}

	return link
}

func (service *Service) FindValidUrlByShortKey(shortKey string) *Domin.Link {

	link := service.GetUrl(shortKey)

	if link != nil && link.Status == Domin.LINK_STATUS_APPROVE {
		// Redirect the user to the original URL
		return link
	}

	return link
}

func (service *Service) UpdateStats(wg *sync.WaitGroup, ctx context.Context) int {

	var bunchOfLinks map[int]*Domin.Link

	limit := 10

	start := 0

	// Setting up the main context
	//ctx1 := context.Background()

	var counter int = 1

	for (bunchOfLinks[0] != nil && counter > 1) || counter == 1 {

		if counter > 1 {
			start = counter * limit
		}

		bunchOfLinks, _ = service.LinkRepo.GetChunk(start, limit, "approve")

		if bunchOfLinks[0] != nil {

			wg.Add(1)

			linkCh := make(chan *Domin.Link) // we consider as unbuffered for synchronize purpose.
			//time.Sleep(20 * time.Second)

			// Goroutine 1
			go service.updateStatusWorker(wg, ctx, bunchOfLinks, linkCh)

			wg.Add(1)
			// Goroutine 2
			go service.updateStatWorker(wg, linkCh)

		}

		counter++
	}

	logger.CreateLogInfo(fmt.Sprintf("[*] Cron end"))

	// Listen for the context cancellation signal to stop the cron scheduler
	<-ctx.Done()
	logger.CreateLogInfo("[*] Cron Received shutdown signal")

	return 1
}

func (service *Service) updateStatusWorker(wg *sync.WaitGroup, ctx context.Context, bunchOfLinks map[int]*Domin.Link, ch chan *Domin.Link) {

	defer wg.Done()

	// Create a new context with a timeout for the processing
	_, cancelProc := context.WithTimeout(ctx, 5*time.Second)

	defer cancelProc()

	//logger.CreateLogInfo(fmt.Sprintf("Run Go routine %d", start))

	for _, data := range bunchOfLinks {

		status := Domin.LINK_STATUS_APPROVE

		if !url.CheckURL(data.Link) {
			logger.CreateLogInfo(fmt.Sprintf("Not approved ShortKey :%v", data.ShortKey))
			status = Domin.Link_STATUS_REJECT
		}

		if data.Status != status {
			service.LinkRepo.UpdateStatus(status, data.Link)
		}

		if data.Status == Domin.LINK_STATUS_APPROVE {
			fmt.Println("Goroutine updateStatusWorker send data...")
			ch <- data // Send data to the channel
		}

	}
	//logger.CreateLogInfo(fmt.Sprintf("Finish Go routine  %d", start))

}

func (service *Service) updateStatWorker(wg *sync.WaitGroup, ch chan *Domin.Link) {
	defer wg.Done()
	fmt.Println("Goroutine updateStatWorker receiving data...")
	data := <-ch // Receive data from the channel
	fmt.Println("Goroutine updateStatWorker received data:", data)

	hget, _ := service.Cache.Get(data.ShortKey)

	fmt.Println("Goroutine updateStatWorker hget:", hget)

	visitCache, _ := strconv.Atoi(hget)

	if visitCache > data.Visit {
		service.LinkRepo.UpdateVisit(visitCache, data.ShortKey)
		logger.CreateLogInfo(fmt.Sprintf("Updated %wg : visit :%v", data.ShortKey, visitCache))
	}
}

func (service *Service) UpdateStatusByLink(status string, link string) {

	service.LinkRepo.UpdateStatus(status, link)

}

func (service *Service) UpdateStatusShortKey(status string, shortKey string, link string) {

	service.LinkRepo.UpdateStatusShortKey(status, shortKey, link)

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

func (service *Service) SetUrl(link string) bool {

	service.createLink(link)

	service.publishQueue(link)

	return true
}

func (service *Service) createLink(link string) string {

	_, err := service.LinkRepo.Create(link, "")

	if err != nil {
		log.Fatal(errors.Wrap(err, "DB has an errorMsg."))
	}
	return ""
}

func (service *Service) publishQueue(link string) {
	var data = make(map[string]string)

	data["link"] = link

	event := Event.Event{Type: Event.CreateLink, Data: data}

	ch, err := service.Queue.Connection.Channel()

	if err != nil {
		log.Fatal(errors.Wrap(err, "Queue has an errorMsg."))
	}

	service.Queue.Publish(ch, service.Shortener.Config.QueueRabbit.MainQueueName, event)
}

//func (Service *Service) GetAllUrl() map[string]string {
//	//return Service.Shortener.Urls
//	return Service.LinkRepo.GetAll()
//}

func (service *Service) GetAllUrlV2() (map[int]*Domin.Link, error) {
	//return Service.Shortener.Urls
	data, err := service.LinkRepo.GetAll()

	// Convert to a slice of interfaces
	var myInterfaceSlice []interface{}
	for _, item := range data {
		myInterfaceSlice = append(myInterfaceSlice, item)
	}

	//service.GenerateShortLink()

	service.MemCache.SetSlice("list", myInterfaceSlice, 5*time.Minute)

	return data, err
}

func (service *Service) GenerateShortLink(count int, isActive bool) string {
	ShortKey, _ := service.ShortKeyRepo.GetLast()

	lastId := 1
	if ShortKey != nil {
		lastId = int(ShortKey.ID) + 1
	}

	var shortLink string

	for i := lastId; i < lastId+count; i++ { // Generate first 100 unique IDs as an example
		//log.Println(h.LinkService.IntToBase62(i))

		shortLink = service.IntToBase62(int(i))

		service.ShortKeyRepo.Create(int(i), shortLink, isActive)

		//logger.CreateLogInfo(h.LinkService.IntToBase62(i))
	}

	return shortLink
}

func (service *Service) GetAllLinkApi() ([]interface{}, error) {
	//return Service.Shortener.Urls
	data, err := service.LinkRepo.GetAll()

	// Convert to a slice of interfaces
	var myInterfaceSlice []interface{}
	for _, item := range data {
		myInterfaceSlice = append(myInterfaceSlice, item)
	}

	var dataMem []interface{}

	//todo use serializer here and make full url in seprated field in serilizer

	// Try to get data from cache
	if dataMem, found := service.MemCache.GetSlice("list"); found {
		return dataMem, nil
	}

	return dataMem, err
}

func (service *Service) VerifyLinkIsValid(link string) string {
	status := Domin.LINK_STATUS_APPROVE
	if !url.CheckURL(link) {
		logger.CreateLogInfo(fmt.Sprintf("[*] Queue Rejected link :%s", link))
		status = Domin.Link_STATUS_REJECT
	}

	if status == Domin.LINK_STATUS_APPROVE {
		short_key := service.GenerateShortLink(1, true)

		service.UpdateStatusShortKey(status, short_key, link)
	} else {
		service.UpdateStatusByLink(status, link)
	}
	return status
}
