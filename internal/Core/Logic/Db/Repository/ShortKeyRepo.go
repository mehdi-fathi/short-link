package Repository

import (
	"database/sql"
	"errors"
	"short-link/internal/Config"
	"short-link/internal/Core/Domin"
	"short-link/internal/Core/Logic/Db"
	"short-link/internal/Core/Ports"
	"short-link/pkg/logger"
)

// Db holds database connection to Postgres
type RepositoryShortKey struct {
	Repository
}

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

	logger.CreateLogInfo("shortKeyTable")
	logger.CreateLogInfo(row)
	logger.CreateLogInfo(&shortKeyTable)

	return &shortKeyTable, err

}

func (db *RepositoryShortKey) FillShortKeyRow(row *sql.Row, linkTable *Domin.ShortKey) error {
	return row.Scan(&linkTable.ID, &linkTable.Uid, &linkTable.IsActive)
}

func CreateShortKeyRepository(cfg *Config.Config, dbLayer *Db.Db) Ports.ShortkeyRepositoryInterface {

	Repo := &RepositoryShortKey{Repository{dbLayer, cfg}}

	return Repo
}
