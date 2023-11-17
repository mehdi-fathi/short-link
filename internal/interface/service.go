package service_interface

import repository_interface "short-link/internal/Db/Repository/interface"

type ServiceInterface interface {
	GetUrl(shortKey string) *repository_interface.Link
	UpdateStats() int
	SetUrl(link string) string
	//GetAllUrl() map[string]string
	GetAllUrlV2() (map[int]*repository_interface.Link, error)
	UpdateStatus(status string, shortKey string)
}
