package service_interface

type ServiceInterface interface {
	GetUrl(shortKey string) string

	SetUrl(link string) string
	GetAllUrl() map[string]string
}
