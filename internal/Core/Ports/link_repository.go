package Ports

import (
	"short-link/internal/Core/Domin"
)

type RepositoryInterface interface {
	FindById(idIn int) (*Domin.Link, error)
	FindByShortKey(shortKey string) (*Domin.Link, error)
	Create(link string, shortKey string) (int, error)
	GetAll() (map[int]*Domin.Link, error)
	GetChunk(start int, limit int, status string) (map[int]*Domin.Link, error)
	UpdateVisit(visit int, shortKey string) (int, error)
	UpdateStatus(status string, link string) (int, error)
	GetByStatus(status string) (map[int]*Domin.Link, error)
}
