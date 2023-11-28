package Ports

import (
	"context"
	"short-link/internal/Core/Domin"
	"sync"
)

type ServiceInterface interface {
	GetUrl(shortKey string) *Domin.Link
	UpdateStats(s *sync.WaitGroup, ctx context.Context) int
	SetUrl(link string) string
	//GetAllUrl() map[string]string
	GetAllUrlV2() (map[int]*Domin.Link, error)
	UpdateStatus(status string, shortKey string)
	GetAllLinkApi() ([]interface{}, error)
}
