package Repository

import (
	"database/sql"
	"errors"
	"short-link/internal/Db"
	repository_interface "short-link/internal/Db/Repository/interface"
)

// Db holds database connection to Postgres
type Repository struct {
	*Db.Db
}


func (db *Repository) FindById(idIn int) (*repository_interface.Link, error) {

	q := `SELECT * FROM links WHERE id=$1;`
	row := db.Sql.QueryRow(q, idIn)
	var err error

	var linkTable repository_interface.Link

	err = row.Scan(&linkTable.ID, &linkTable.Link)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &linkTable, err

}

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateRepository(dbLayer *Db.Db) repository_interface.RepositoryInterface {

	Repo := &Repository{dbLayer}

	return Repo
}
