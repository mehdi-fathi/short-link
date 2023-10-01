package Db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Db holds database connection to Postgres
type Db struct {
	Sql *sql.DB
}

// database variables
// usually we should get them from env like os.Getenv("variableName")
const (
	host     = "localhost"
	port     = 5432
	user     = "default"
	password = "secret"
	dbname   = "slink"
)

// ConnectDB tries to connect DB and on succcesful it returns
// DB connection string and nil error, otherwise return empty DB and the corresponding error.
func (db *Db) ConnectDB() (*Db, error) {

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable", host, port, user, password, dbname)

	var err error

	db.Sql, err = sql.Open("postgres", connString)
	log.Println("err", err)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return &Db{}, err
	}

	rows, errorIs := db.Sql.Query("SELECT version();")

	log.Println("errorIs", errorIs)
	if errorIs != nil {
		panic(errorIs)
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
func CreateDb() *Db {

	Db := &Db{
		Sql: new(sql.DB),
	}

	Db.ConnectDB()

	return Db
}
