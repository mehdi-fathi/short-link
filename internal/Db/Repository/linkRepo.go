package Repository

import (
	"database/sql"
	"errors"
	"short-link/internal/Config"
	"short-link/internal/Db"
	repository_interface "short-link/internal/Db/Repository/interface"
)

// Db holds database connection to Postgres
type Repository struct {
	*Db.Db
	Config *Config.Config
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

	err = row.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey, &linkTable.Visit, &linkTable.UpdatedAt, &linkTable.Status)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &linkTable, err

}

func (db *Repository) GetAll() (map[int]*repository_interface.Link, error) {

	q := `SELECT * FROM links order by id desc ;`
	rows, _ := db.Sql.Query(q)
	var err error

	var links = make(map[int]*repository_interface.Link)

	for i := 0; rows.Next(); i++ {
		var linkTable repository_interface.Link

		err = rows.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey, &linkTable.Visit, &linkTable.UpdatedAt, &linkTable.Status)

		links[i] = &linkTable

	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return links, err

}

func (db *Repository) GetByStatus(status string) (map[int]*repository_interface.Link, error) {

	q := `SELECT * FROM links where status = $1 limit 100;`
	rows, _ := db.Sql.Query(q, status)

	var err error

	var links = make(map[int]*repository_interface.Link)

	for i := 0; rows.Next(); i++ {
		var linkTable repository_interface.Link

		err = rows.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey, &linkTable.Visit, &linkTable.UpdatedAt, &linkTable.Status)

		links[i] = &linkTable

	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return links, err

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

func (db *Repository) UpdateStatus(status string, link string) (int, error) {

	q := `update links set status = $1,updated_at=now()
		  where link = $2;`

	row := db.Sql.QueryRow(q, status, link)

	var id int

	row.Scan(&id)

	var err error

	return id, err

}

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateRepository(cfg *Config.Config, dbLayer *Db.Db) repository_interface.RepositoryInterface {

	Repo := &Repository{dbLayer, cfg}

	return Repo
}
