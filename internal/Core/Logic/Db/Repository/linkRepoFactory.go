package Repository

import (
	"short-link/internal/Config"
	"short-link/internal/Core/Logic/Db"
	"short-link/internal/Core/Ports"
)

// Db holds database connection to Postgres
type Repository struct {
	*Db.Db
	Config *Config.Config
}

func CreateLinkRepository(cfg *Config.Config, dbLayer *Db.Db) Ports.LinkRepositoryInterface {

	Repo := &Repository{dbLayer, cfg}

	return Repo
}
