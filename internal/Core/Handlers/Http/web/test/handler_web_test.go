package test

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"short-link/internal/Cache"
	"short-link/internal/Cache/MemCache"
	"short-link/internal/Config"
	"short-link/internal/Core/Domin"
	"short-link/internal/Core/Handlers/Http/web"
	"short-link/internal/Core/Logic/Db"
	"short-link/internal/Core/Logic/Db/Repository"
	"short-link/internal/Core/Logic/Service"
	"short-link/internal/Infrastructure"
	"short-link/internal/Queue"
	"short-link/pkg/logger"
	"strings"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// test_utils.go
func initTest() (*web.HandlerWeb, *Infrastructure.Handler, *gin.Engine, *Db.Db, error) {
	os.Setenv("GO_ENV", "test")
	cfg := Config.LoadConfigEnvApp()

	db, err := initTestDB(cfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	handlerWeb, handlerMain, router := setupRouterAndHandler(cfg, db)
	return handlerWeb, handlerMain, router, db, nil
}

func initTestDB(cfg *Config.Config) (*Db.Db, error) {

	// connect to DB first
	var errDb error

	Db := Db.CreateDb(cfg)

	Db, err := Db.ConnectDBTest()
	if err != nil {
		log.Fatalf("setupTestDB failed: %v", err)
	}
	if errDb != nil {
		log.Fatalf("failed to start the server: %v", errDb)
	}

	logger.CreateLogger(cfg.Graylog)

	return Db, nil
}

func setupRouterAndHandler(cfg *Config.Config, db *Db.Db) (*web.HandlerWeb, *Infrastructure.Handler, *gin.Engine) {

	repoLink := Repository.CreateLinkRepository(cfg, db)
	repoShortKey := Repository.CreateShortKeyRepository(cfg, db)

	//httpHandler := &Handler{
	//	LinkService: internal.CreateLinkService(cfg),
	//}

	cache := Cache.CreateCache(cfg)

	queue := Queue.CreateQueue(cfg)

	memCache := MemCache.CreateMemCache(cfg)

	//LinkService := tt.initService()
	//var configServer ConfigModel

	var ser = Service.CreateLinkService(cfg, repoLink, repoShortKey, cache, memCache, queue)
	handlerWeb := Infrastructure.CreateHandlerWeb(ser)
	handler := Infrastructure.CreateHandlerMain()
	//handler := CreateHandler(LinkService,bookstore.CreateLinkService(nil))
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = ioutil.Discard
	router := gin.Default()
	log.Println(cfg.AppPath)
	router.LoadHTMLGlob(cfg.AppPath + "tmp/*")

	return handlerWeb, handler, router
}

// TestHandleShorten tests the POST /shorten endpoint
func TestHandleShorten(t *testing.T) {

	runTest(t, func(t *testing.T, handler *web.HandlerWeb, handlerMain *Infrastructure.Handler, router *gin.Engine, dbLayer *Db.Db) {

		log.Println("router", router)
		router.POST("/make", handler.HandleShorten)

		// Create the form data for the POST body
		formData := url.Values{
			"link": {"https://www.google.com"},
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
	})

}

func runTest(t *testing.T, testFunc func(t *testing.T, handler *web.HandlerWeb, handlerMain *Infrastructure.Handler, router *gin.Engine, dbLayer *Db.Db)) {
	handler, handlerMain, router, dbLayer, err := initTest()

	if err != nil {
		t.Fatalf("initTest failed: %v", err)
	}
	//defer teardownTestDB(dbLayer.Sql)

	testFunc(t, handler, handlerMain, router, dbLayer)
}

func TestHandleRedirectNotFound(t *testing.T) {

	runTest(t, func(t *testing.T, handler *web.HandlerWeb, handlerMain *Infrastructure.Handler, router *gin.Engine, dbLayer *Db.Db) {

		router.GET("/short/:url", handler.HandleRedirect)

		// Create a new HTTP request with the form data
		req, err := http.NewRequest(http.MethodGet, "/short/xhdts", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Record the response
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestHandleRedirectSuccess(t *testing.T) {

	runTest(t, func(t *testing.T, handler *web.HandlerWeb, handlerMain *Infrastructure.Handler, router *gin.Engine, dbLayer *Db.Db) {

		handler.LinkService.SetUrl("https://www.google.com")

		router.GET("/short/:url", handler.HandleRedirect)

		shortLink := handler.LinkService.GenerateShortLink(1, true)
		status := Domin.LINK_STATUS_APPROVE

		handler.LinkService.UpdateStatusShortKey(status, shortLink, "https://www.google.com")

		// Create a new HTTP request with the form data
		req, err := http.NewRequest(http.MethodGet, "/short/"+shortLink, nil)
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
	})
}

func TestHandleListAll(t *testing.T) {

	runTest(t, func(t *testing.T, handler *web.HandlerWeb, handlerMain *Infrastructure.Handler, router *gin.Engine, dbLayer *Db.Db) {

		handler.LinkService.SetUrl("https://www.google.com")

		router.GET("/list/all", handler.HandleList)

		// Create a new HTTP request with the form data
		req, err := http.NewRequest(http.MethodGet, "/list/all", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// Record the response
		w := httptest.NewRecorder()

		// Perform the request
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
	})

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
