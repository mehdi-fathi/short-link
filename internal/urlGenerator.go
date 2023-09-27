package internal

import (
	"math/rand"
	service_interface "short-link/internal/interface"
	"time"
)

type UrlShortener struct {
	Urls   map[string]string
	Config *Config
}

type Service struct {
	Shortener *UrlShortener
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
func CreateService(cfg *Config) service_interface.Service {

	shortenerUrl := &UrlShortener{
		Urls:   make(map[string]string),
		Config: cfg,
	}

	return &Service{
		Shortener: shortenerUrl,
	}
}

func (service *Service) GetUrl(shortKey string) string {
	return service.Shortener.Urls[shortKey]
}

func (service *Service) SetUrl(link string) string {

	shortKey := GenerateShortKey(service.Shortener.Config.HASHCODE)
	service.Shortener.Urls[shortKey] = link
	return shortKey
}
