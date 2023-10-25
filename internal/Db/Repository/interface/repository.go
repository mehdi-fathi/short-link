package repository_interface

// Character is one character from the database.
type Link struct {
	ID       int64
	Link     string
	ShortKey string
	Visit    int
}

type RepositoryInterface interface {
	FindById(idIn int) (*Link, error)
	FindByShortKey(shortKey string) (*Link, error)
	Create(link string, shortKey string) (int, error)
	GetAll() (map[int]*Link, error)
	UpdateVisit(visit int, shortKey string) (int, error)
}
