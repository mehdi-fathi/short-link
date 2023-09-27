package service_interface

type Service interface {
	GetUrl(shortKey string) string

	SetUrl(link string) string
	GetAllUrl() map[string]string
}
