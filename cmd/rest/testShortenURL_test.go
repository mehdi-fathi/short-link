package rest

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"short-link/internal"
	"short-link/internal/Cache"
	"short-link/internal/Config"
	"short-link/internal/Db"
	"short-link/internal/Db/Repository"
	"short-link/internal/Queue"
	"short-link/pkg/logger"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

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

type DbNew struct {
	*Db.Db
}

func usage() {
	usageStr := `
Usage: server [options]
Options:
	-c,  --config   <config file name>   Path of yaml configuration file
`
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

// TestShortenURL tests the POST /shorten endpoint
func TestShortenURL(t *testing.T) {

	// Default Config file based on the environment variable
	defaultConfigFile := "/Users/mehdi/Sites/short-link/config/config-test.yaml"
	if env := os.Getenv("APP_MODE"); env != "" {
		defaultConfigFile = fmt.Sprintf("config/config-%s.yaml", env)
	}

	// Load Master Config File
	var configFile string
	flag.StringVar(&configFile, "config", defaultConfigFile, "The environment configuration file of application")
	flag.Usage = usage
	flag.Parse()

	// Loading the config file
	cfg, err := Config.LoadConfig(configFile)
	if err != nil {
		log.Println(errors.Wrapf(err, "failed to load config: %s", "CreateService"))
	}

	log.Println(cfg.Logger)

	loggerInstance := logger.CreateLogger(cfg.Logger)
	loggerInstance.Info("[OK] Logger Configured")

	//loggerInstance := logrus.Logger{}
	//loggerInstance.Info("[OK] Logger Configured")

	// Setting up the main context

	// connect to DB first
	var errDb error

	Db := &Db.Db{
		Sql:    new(sql.DB),
		Config: cfg,
	}

	dbLayer, err := Db.ConnectDBTest()
	if err != nil {
		t.Fatalf("setupTestDB failed: %v", err)
	}
	if errDb != nil {
		log.Fatalf("failed to start the server: %v", errDb)
	}

	repo := Repository.CreateRepository(cfg, dbLayer)

	defer dbLayer.Sql.Close()

	defer func() {

		if err := teardownTestDB(dbLayer.Sql); err != nil {
			t.Fatalf("teardownTestDB failed: %v", err)
		}
	}()

	//httpHandler := &Handler{
	//	Service: internal.CreateService(cfg),
	//}

	client := Cache.CreateCache()

	queue := Queue.CreateQueue(cfg)

	//service := tt.initService()
	//var configServer ConfigModel

	var ser = internal.CreateService(cfg, repo, client, queue)
	handler := CreateHandler(ser)
	//handler := CreateHandler(service,bookstore.CreateService(nil))
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = ioutil.Discard
	router := gin.Default()

	router.POST("/make", handler.HandleShorten)

	// Create the form data for the POST body
	formData := url.Values{
		"link": {"https://www.example.com"},
	}

	// Create a new HTTP request with the form data
	req, err := http.NewRequest(http.MethodPost, "/make", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Record the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status code %d, got %d", http.StatusMovedPermanently, w.Code)
	}
}

func teardownTestDB(db *sql.DB) error {

	dropAllTables(db)
	return nil
}

func dropAllTables(db *sql.DB) error {
	// Query to get the list of all tables in the current database
	query := `SELECT table_name FROM information_schema.tables WHERE table_schema='public'`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Slice to hold all table names
	var tables []string

	// Iterate over the rows and append table names to the slice
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		tables = append(tables, tableName)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return err
	}

	// Disable foreign key checks to avoid issues with dependent tables
	_, err = db.Exec("SET CONSTRAINTS ALL DEFERRED")
	if err != nil {
		return err
	}

	// Drop each table
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			return err
		}
	}

	// Re-enable foreign key checks
	_, err = db.Exec("SET CONSTRAINTS ALL IMMEDIATE")
	if err != nil {
		return err
	}

	return nil
}
