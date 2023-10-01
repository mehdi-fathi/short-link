package Repository

import (
	"log"
	"short-link/internal/Db"
)

// Db holds database connection to Postgres
type Repository struct {
	*Db.Db
}

func (db *Repository) FindById(idIn int) {

	var id int
	var link string

	q := `SELECT * FROM links WHERE id=$1;`
	row := db.Sql.QueryRow(q, idIn)
	var err error

	err = row.Scan(&id, &link)

	log.Println("err", err)

	log.Printf("we are able to fetch album with given link: %s", link)

}

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateRepository(dbLayer *Db.Db) *Repository {

	Repo := &Repository{dbLayer}

	return Repo
}
