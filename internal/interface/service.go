package service_interface

import (
	"context"
	"short-link/internal/Db/Model"
	"sync"
)

type ServiceInterface interface {
	GetUrl(shortKey string) *Model.Link
	UpdateStats(s *sync.WaitGroup, ctx context.Context) int
	SetUrl(link string) string
	//GetAllUrl() map[string]string
	GetAllUrlV2() (map[int]*Model.Link, error)
	UpdateStatus(status string, shortKey string)
	GetAllLinkApi() ([]interface{}, error)
}
