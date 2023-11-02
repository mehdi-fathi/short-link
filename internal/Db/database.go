package Db

import (
	"database/sql"
	"fmt"
	"log"
	"short-link/internal/Config"

	_ "github.com/lib/pq"
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
	log.Println(db.Config.DB.User)

	var err error

	db.Sql, err = sql.Open("postgres", connString)

	if err != nil {
		log.Printf("failed to connect to database: %v", err)
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
