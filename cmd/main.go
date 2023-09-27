package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"short-link/internal"
	service_interface "short-link/internal/interface"
)

type Handler struct {
	Service service_interface.Service
}

func main() {

	// Default Config file based on the environment variable
	defaultConfigFile := "config/config-local.yaml"
	if env := os.Getenv("APP_MODE"); env != "" {
		defaultConfigFile = fmt.Sprintf("config/config-%s.yaml", env)
	}

	// Load Master Config File
	var configFile string
	flag.StringVar(&configFile, "config", defaultConfigFile, "The environment configuration file of application")
	flag.Usage = usage
	flag.Parse()

	// Loading the config file
	cfg, err := internal.LoadConfig(configFile)
	if err != nil {
		log.Println(errors.Wrapf(err, "failed to load config: %s", "CreateService"))
	}

	httpHandler := &Handler{
		Service: internal.CreateService(cfg),
	}

	router := gin.Default()

	router.LoadHTMLGlob("tmp/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.POST("/make", httpHandler.HandleShorten)
	router.GET("/short/:url", httpHandler.HandleRedirect)
	router.GET("/list/all", httpHandler.HandleList)

	router.Run() // listen and serve on 0.0.0.0:8080
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

func (us *Handler) HandleRedirect(c *gin.Context) {

	shortKey := c.Param("url")

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL := us.Service.GetUrl(shortKey)

	log.Println(originalURL, shortKey)

	// Redirect the user to the original URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (us *Handler) HandleList(c *gin.Context) {

	allUrl := us.Service.GetAllUrl()

	log.Println(allUrl,"hi")

	c.HTML(http.StatusOK, "list.html", gin.H{
		"data": allUrl,
	})
}

func (us *Handler) HandleShorten(c *gin.Context) {
	//if r.Method != http.MethodPost {
	//	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	//	return
	//}
	//
	//originalURL := r.FormValue("url")
	//if originalURL == "" {
	//	http.Error(w, "URL parameter is missing", http.StatusBadRequest)
	//	return
	//}

	link := c.PostForm("link")

	// Generate a unique shortened key for the original URL
	shortKey := us.Service.SetUrl(link)

	log.Println(link)

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	log.Println(shortenedURL)

}
