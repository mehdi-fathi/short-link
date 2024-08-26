package Service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"short-link/internal/Core/Domin"
	"short-link/internal/Event"
	"short-link/pkg/logger"
	"short-link/pkg/url"
	"strconv"
	"sync"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func (linkService *LinkService) IntToBase62(num int) string {
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

func (linkService *LinkService) GetUrl(shortKey string) *Domin.Link {

	link, _ := linkService.LinkRepo.FindByShortKey(shortKey)

	if link != nil && link.Link != "" {
		linkService.Cache.IncrBy(shortKey, 1)
	}

	return link
}

func (linkService *LinkService) FindValidUrlByShortKey(shortKey string) *Domin.Link {

	link := linkService.GetUrl(shortKey)

	if link != nil && link.Status == Domin.LINK_STATUS_APPROVE {
		// Redirect the user to the original URL
		return link
	}

	return link
}

func (linkService *LinkService) UpdateStats(wg *sync.WaitGroup, ctx context.Context) int {

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

		bunchOfLinks, _ = linkService.LinkRepo.GetChunk(start, limit, "approve")

		if bunchOfLinks[0] != nil {

			wg.Add(1)

			linkCh := make(chan *Domin.Link) // we consider as unbuffered for synchronize purpose.
			//time.Sleep(20 * time.Second)

			// Goroutine 1
			go linkService.updateStatusWorker(wg, ctx, bunchOfLinks, linkCh)

			wg.Add(1)
			// Goroutine 2
			go linkService.updateStatWorker(wg, linkCh)

		}

		counter++
	}

	logger.CreateLogInfo(fmt.Sprintf("[*] Cron end"))

	// Listen for the context cancellation signal to stop the cron scheduler
	<-ctx.Done()
	logger.CreateLogInfo("[*] Cron Received shutdown signal")

	return 1
}

func (linkService *LinkService) updateStatusWorker(wg *sync.WaitGroup, ctx context.Context, bunchOfLinks map[int]*Domin.Link, ch chan *Domin.Link) {

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
			linkService.LinkRepo.UpdateStatus(status, data.Link)
		}

		if data.Status == Domin.LINK_STATUS_APPROVE {
			fmt.Println("Goroutine updateStatusWorker send data...")
			ch <- data // Send data to the channel
		}

	}
	//logger.CreateLogInfo(fmt.Sprintf("Finish Go routine  %d", start))

}

func (linkService *LinkService) updateStatWorker(wg *sync.WaitGroup, ch chan *Domin.Link) {
	defer wg.Done()
	fmt.Println("Goroutine updateStatWorker receiving data...")
	data := <-ch // Receive data from the channel
	fmt.Println("Goroutine updateStatWorker received data:", data)

	hget, _ := linkService.Cache.Get(data.ShortKey)

	fmt.Println("Goroutine updateStatWorker hget:", hget)

	visitCache, _ := strconv.Atoi(hget)

	if visitCache > data.Visit {
		linkService.LinkRepo.UpdateVisit(visitCache, data.ShortKey)
		logger.CreateLogInfo(fmt.Sprintf("Updated %wg : visit :%v", data.ShortKey, visitCache))
	}
}

func (linkService *LinkService) UpdateStatusByLink(status string, link string) {

	linkService.LinkRepo.UpdateStatus(status, link)

}

func (linkService *LinkService) UpdateStatusShortKey(status string, shortKey string, link string) {

	linkService.LinkRepo.UpdateStatusShortKey(status, shortKey, link)

}

func (linkService *LinkService) checkPendingLinks() int {

	all, _ := linkService.LinkRepo.GetByStatus("pending")

	var status string
	for _, data := range all {
		status = "approve"
		logger.CreateLogInfo(data.Link)

		if !url.CheckURL(data.Link) {
			logger.CreateLogInfo(fmt.Sprintf("Not approved ShortKey :%v", data.ShortKey))
			status = "reject"
		}

		linkService.LinkRepo.UpdateStatus(status, data.ShortKey)

	}
	return 1
}

func (linkService *LinkService) SetUrl(link string) bool {

	linkService.createLink(link)

	linkService.publishQueue(link)

	return true
}

func (linkService *LinkService) createLink(link string) string {

	_, err := linkService.LinkRepo.Create(link, "")

	if err != nil {
		log.Fatal(errors.Wrap(err, "DB has an errorMsg."))
	}
	return ""
}

func (linkService *LinkService) publishQueue(link string) {
	var data = make(map[string]string)

	data["link"] = link

	event := Event.Event{Type: Event.CreateLink, Data: data}

	ch, err := linkService.Queue.Connection.Channel()

	if err != nil {
		log.Fatal(errors.Wrap(err, "Queue has an errorMsg."))
	}

	linkService.Queue.Publish(ch, linkService.Config.QueueRabbit.MainQueueName, event)
}

//func (LinkService *LinkService) GetAllUrl() map[string]string {
//	//return LinkService.Shortener.Urls
//	return LinkService.LinkRepo.GetAll()
//}

func (linkService *LinkService) GetAllUrlV2() (map[int]*Domin.Link, error) {
	//return LinkService.Shortener.Urls
	data, err := linkService.LinkRepo.GetAll()

	// Convert to a slice of interfaces
	var myInterfaceSlice []interface{}
	for _, item := range data {
		myInterfaceSlice = append(myInterfaceSlice, item)
	}

	//service.GenerateShortLink()

	linkService.MemCache.SetSlice("list", myInterfaceSlice, 5*time.Minute)

	return data, err
}

func (linkService *LinkService) GenerateShortLink(count int, isActive bool) string {
	ShortKey, _ := linkService.ShortKeyRepo.GetLast()

	lastId := 1
	if ShortKey != nil {
		lastId = int(ShortKey.ID) + 1
	}

	var shortLink string

	for i := lastId; i < lastId+count; i++ { // Generate first 100 unique IDs as an example
		//log.Println(h.LinkService.IntToBase62(i))

		shortLink = linkService.IntToBase62(int(i))

		linkService.ShortKeyRepo.Create(int(i), shortLink, isActive)

		//logger.CreateLogInfo(h.LinkService.IntToBase62(i))
	}

	return shortLink
}

func (linkService *LinkService) GetAllLinkApi() ([]interface{}, error) {
	//return LinkService.Shortener.Urls
	data, err := linkService.LinkRepo.GetAll()

	// Convert to a slice of interfaces
	var myInterfaceSlice []interface{}
	for _, item := range data {
		myInterfaceSlice = append(myInterfaceSlice, item)
	}

	var dataMem []interface{}

	//todo use serializer here and make full url in seprated field in serilizer

	// Try to get data from cache
	if dataMem, found := linkService.MemCache.GetSlice("list"); found {
		return dataMem, nil
	}

	return dataMem, err
}

func (linkService *LinkService) VerifyLinkIsValid(link string) string {
	status := Domin.LINK_STATUS_APPROVE
	if !url.CheckURL(link) {
		logger.CreateLogInfo(fmt.Sprintf("[*] Queue Rejected link :%s", link))
		status = Domin.Link_STATUS_REJECT
	}

	if status == Domin.LINK_STATUS_APPROVE {
		short_key := linkService.GenerateShortLink(1, true)

		linkService.UpdateStatusShortKey(status, short_key, link)
	} else {
		linkService.UpdateStatusByLink(status, link)
	}
	return status
}
