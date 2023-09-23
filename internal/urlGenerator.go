package internal

import (
	"math/rand"
	service_interface "short-link/internal/interface"
	"time"
)

type UrlShortener struct {
	Urls map[string]string
}

type Service struct {
	Shortener *UrlShortener
}

func GenerateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateService() service_interface.Service {

	shortenerUrl := &UrlShortener{
		Urls: make(map[string]string),
	}

	return &Service{
		Shortener: shortenerUrl,
	}
}

func (service *Service) GetUrl(shortKey string) string {
	return service.Shortener.Urls[shortKey]
}

func (service *Service) SetUrl(link string) string {

	shortKey := GenerateShortKey()
	service.Shortener.Urls[shortKey] = link
	return shortKey
}
