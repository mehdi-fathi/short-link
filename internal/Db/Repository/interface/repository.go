package repository_interface

import "short-link/internal/Db/Model"

type RepositoryInterface interface {
	FindById(idIn int) (*Model.Link, error)
	FindByShortKey(shortKey string) (*Model.Link, error)
	Create(link string, shortKey string) (int, error)
	GetAll() (map[int]*Model.Link, error)
	GetChunk(start int, limit int, status string) (map[int]*Model.Link, error)
	UpdateVisit(visit int, shortKey string) (int, error)
	UpdateStatus(status string, link string) (int, error)
	GetByStatus(status string) (map[int]*Model.Link, error)
}
