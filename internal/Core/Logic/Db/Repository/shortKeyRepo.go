package Repository

import (
	"database/sql"
	"errors"
	"short-link/internal/Core/Domin"
)

func (db *RepositoryShortKey) Create(id int, uid string, isActive bool) (int, error) {

	var idCreated int

	q := `insert into short_keys (id,uid,is_active) values($1,$2,$3) returning id;`
	row := db.Sql.QueryRow(q, id, uid, isActive)

	row.Scan(&idCreated)

	var err error

	return idCreated, err

}

func (db *RepositoryShortKey) GetLast() (*Domin.ShortKey, error) {

	q := `SELECT * 
		  FROM short_keys 
		  ORDER BY id DESC
		  LIMIT 1;`

	row := db.Sql.QueryRow(q)

	var err error

	var shortKeyTable Domin.ShortKey

	err = row.Scan(&shortKeyTable.ID, &shortKeyTable.Uid, &shortKeyTable.IsActive)
	//err = db.FillShortKeyRow(row, &shortKeyTable)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &shortKeyTable, err

}

func (db *RepositoryShortKey) FillShortKeyRow(row *sql.Row, linkTable *Domin.ShortKey) error {
	return row.Scan(&linkTable.ID, &linkTable.Uid, &linkTable.IsActive)
}
