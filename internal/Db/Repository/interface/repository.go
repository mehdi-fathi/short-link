package repository_interface

// Character is one character from the database.
type Link struct {
	ID        int64
	Link      string
	ShortKey  string
	Visit     int
	UpdatedAt string
	Status    string
}

const  (
	LINK_STATUS_APPROVE = "approve"
	LINK_STATUS_REJECT = "reject"
	LINK_STATUS_PENDING = "pending"

)

type RepositoryInterface interface {
	FindById(idIn int) (*Link, error)
	FindByShortKey(shortKey string) (*Link, error)
	Create(link string, shortKey string) (int, error)
	GetAll() (map[int]*Link, error)
	UpdateVisit(visit int, shortKey string) (int, error)
	UpdateStatus(status string, link string) (int, error)
	GetByStatus(status string) (map[int]*Link, error)
}
