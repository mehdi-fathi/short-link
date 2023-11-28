package Db

import (
	"database/sql"
	"fmt"
	"short-link/internal/Config"
	"short-link/pkg/logger"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Db holds database connection to Postgres
type Db struct {
	Sql    *sql.DB
	Config *Config.Config
}

// ConnectDB tries to connect DB and on succcesful it returns
// DB connection string and nil error, otherwise return empty DB and the corresponding error.
func (db *Db) ConnectDB() (*Db, error) {

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable",
		db.Config.DB.Host,
		db.Config.DB.Port,
		db.Config.DB.User,
		db.Config.DB.Password,
		db.Config.DB.Dbname)

	var err error

	db.Sql, err = sql.Open("postgres", connString)

	if err != nil {
		logger.CreateLogError(fmt.Sprintf("failed to connect to database: %v", err))
		return &Db{}, err
	}

	rows, errorIs := db.Sql.Query("SELECT version();")

	if errorIs != nil {
		panic(errorIs.(interface{}))
	}
	//
	for rows.Next() {
		var ver string
		rows.Scan(&ver)
		fmt.Println(ver)
	}

	return db, nil
}

// CreateService creates an instance of membership interface with the necessary dependencies
func CreateDb(cfg *Config.Config) *Db {

	Db := &Db{
		Sql:    new(sql.DB),
		Config: cfg,
	}

	Db.ConnectDB()

	return Db
}

// ConnectDB tries to connect DB and on succcesful it returns
// DB connection string and nil error, otherwise return empty DB and the corresponding error.
func (db *Db) ConnectDBTest() (*Db, error) {

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable",
		db.Config.DB.Host,
		db.Config.DB.Port,
		db.Config.DB.User,
		db.Config.DB.Password,
		db.Config.DB.Dbname)

	var err error

	// Create a test database (pseudo-code)
	db.Sql, err = sql.Open("postgres", connString)
	if err != nil {
		logger.CreateLogError(err.Error())
		return nil, err
	}
	rows, errorIs := db.Sql.Query("SELECT version();")

	if errorIs != nil {
		panic(errorIs.(interface{}))
	}
	//
	for rows.Next() {
		var ver string
		rows.Scan(&ver)
	}

	migrateTest(db.Config.DB)

	return db, nil
}

func applyMigrations(dbURL string, migrationsPath string) error {
	m, err := migrate.New(
		migrationsPath,
		dbURL,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func migrateTest(db Config.DB) {

	connStringMigrate := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Dbname)

	// Apply migrations
	if err := downMigrations(connStringMigrate, "file://../../../../../../database/migration"); err != nil {
		panic(err.(interface{}))
	}
	// Apply migrations
	if err := applyMigrations(connStringMigrate, "file://../../../../../../database/migration"); err != nil {
		panic(err.(interface{}))
	}
}

func downMigrations(dbURL string, migrationsPath string) error {
	m, err := migrate.New(
		migrationsPath,
		dbURL,
	)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
