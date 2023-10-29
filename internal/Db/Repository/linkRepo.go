package Repository

import (
	"database/sql"
	"errors"
	"short-link/internal"
	"short-link/internal/Db"
	repository_interface "short-link/internal/Db/Repository/interface"
)

// Db holds database connection to Postgres
type Repository struct {
	*Db.Db
	Config *internal.Config
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

func (db *Repository) FindByShortKey(shortKey string) (*repository_interface.Link, error) {

	q := `SELECT * FROM links WHERE short_key=$1;`
	row := db.Sql.QueryRow(q, shortKey)

	var err error

	var linkTable repository_interface.Link

	err = row.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey, &linkTable.Visit, &linkTable.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &linkTable, err

}

func (db *Repository) GetAll() (map[int]*repository_interface.Link, error) {

	q := `SELECT * FROM links ;`
	rows, _ := db.Sql.Query(q)
	var err error

	var users = make(map[int]*repository_interface.Link)

	for i := 0; rows.Next(); i++ {
		var linkTable repository_interface.Link

		err = rows.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey, &linkTable.Visit, &linkTable.UpdatedAt)

		users[i] = &linkTable

	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return users, err

}

func (db *Repository) Create(link string, shortKey string) (int, error) {

	q := `insert into links (link,short_key) values($1,$2) returning id;`
	row := db.Sql.QueryRow(q, link, shortKey)

	var id int

	row.Scan(&id)

	var err error

	return id, err

}

func (db *Repository) UpdateVisit(visit int, shortKey string) (int, error) {

	q := `update links set visit = $1,updated_at=now() where short_key = $2;`
	row := db.Sql.QueryRow(q, visit, shortKey)

	var id int

	row.Scan(&id)

	var err error

	return id, err

}

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateRepository(cfg *internal.Config, dbLayer *Db.Db) repository_interface.RepositoryInterface {

	Repo := &Repository{dbLayer, cfg}

	return Repo
}
