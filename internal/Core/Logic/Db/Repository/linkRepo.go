package Repository

import (
	"database/sql"
	"errors"
	"short-link/internal/Core/Domin"
)

func (db *Repository) FindById(idIn int) (*Domin.Link, error) {

	q := `SELECT * FROM links WHERE id=$1;`
	row := db.Sql.QueryRow(q, idIn)
	var err error

	var linkTable Domin.Link

	err = row.Scan(&linkTable.ID, &linkTable.Link)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &linkTable, err

}

func (db *Repository) FindByShortKey(shortKey string) (*Domin.Link, error) {

	q := `SELECT *
		  FROM links
		  WHERE short_key=$1;`

	row := db.Sql.QueryRow(q, shortKey)

	var err error

	var linkTable Domin.Link

	//err = row.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey, &linkTable.Visit, &linkTable.UpdatedAt, &linkTable.Status, &linkTable.CreatedAt)
	err = db.FillLinkRow(row, &linkTable)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &linkTable, err

}

func (db *Repository) FillLinkRow(row *sql.Row, linkTable *Domin.Link) error {
	return row.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey, &linkTable.Visit, &linkTable.UpdatedAt, &linkTable.Status, &linkTable.CreatedAt)
}

func (db *Repository) FillLinkRows(rows *sql.Rows, linkTable *Domin.Link) error {
	return rows.Scan(&linkTable.ID, &linkTable.Link, &linkTable.ShortKey, &linkTable.Visit, &linkTable.UpdatedAt, &linkTable.Status, &linkTable.CreatedAt)
}

func (db *Repository) GetAll() (map[int]*Domin.Link, error) {

	q := `SELECT *
		  FROM links
		  order by id desc ;`

	rows, _ := db.Sql.Query(q)
	var err error

	var links = make(map[int]*Domin.Link)

	if rows != nil {

		for i := 0; rows.Next(); i++ {
			var linkTable Domin.Link

			err = db.FillLinkRows(rows, &linkTable)

			links[i] = &linkTable

		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return links, err

}

func (db *Repository) GetChunk(start int, limit int, status string) (map[int]*Domin.Link, error) {

	q := `SELECT * FROM links
		  where status = $3
		  order by id desc
		  limit $2 OFFSET $1 ;`

	rows, _ := db.Sql.Query(q, start, limit, status)
	var err error

	var links = make(map[int]*Domin.Link)

	if rows != nil {

		for i := 0; rows.Next(); i++ {
			var linkTable Domin.Link

			err = db.FillLinkRows(rows, &linkTable)

			links[i] = &linkTable

		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return links, err

}

func (db *Repository) GetByStatus(status string) (map[int]*Domin.Link, error) {

	q := `SELECT * 
		  FROM links
		  where status = $1
		  limit 100;`

	rows, _ := db.Sql.Query(q, status)

	var err error

	var links = make(map[int]*Domin.Link)

	if rows != nil {

		for i := 0; rows.Next(); i++ {
			var linkTable Domin.Link

			err = db.FillLinkRows(rows, &linkTable)

			links[i] = &linkTable

		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return links, err

}

func (db *Repository) Create(link string, shortKey string) (int, error) {

	var id int

	q := `insert into links (link,short_key) values($1,$2) returning id;`
	row := db.Sql.QueryRow(q, link, shortKey)

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

func (db *Repository) UpdateStatusShortKey(status string, shortKey string, link string) (int, error) {

	q := `update links set
                 status = $1,
                 short_key = $2,
                 updated_at=now()
		  where link = $3;`

	row := db.Sql.QueryRow(q, status, shortKey, link)

	var id int

	row.Scan(&id)

	var err error

	return id, err

}
