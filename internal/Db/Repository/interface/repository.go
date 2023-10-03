package repository_interface


// Character is one character from the database.
type Link struct {
	ID   int64
	Link string
}

type RepositoryInterface interface {
	FindById(idIn int) (*Link, error)
}
