package Repository

import (
	"short-link/internal/Config"
	"short-link/internal/Core/Logic/Db"
	"short-link/internal/Core/Ports"
)

// Db holds database connection to Postgres
type RepositoryShortKey struct {
	Repository
}

func CreateShortKeyRepository(cfg *Config.Config, dbLayer *Db.Db) Ports.ShortKeyRepositoryInterface {

	Repo := &RepositoryShortKey{Repository{dbLayer, cfg}}

	return Repo
}
