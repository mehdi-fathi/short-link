package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"short-link/internal"
)

type UrlShortener struct {
	Urls map[string]string
}

type Handler struct {
	UrlShortener *internal.UrlShortener
}

func main() {

	shortenerUrl := &internal.UrlShortener{
		Urls: make(map[string]string),
	}

	shortener := &Handler{
		UrlShortener: shortenerUrl,
	}

	router := gin.Default()

	router.LoadHTMLGlob("tmp/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.POST("/make", shortener.HandleShorten)
	router.GET("/short/:url", shortener.HandleRedirect)

	router.Run() // listen and serve on 0.0.0.0:8080
}

func (us *Handler) HandleRedirect(c *gin.Context) {

	shortKey := c.Param("url")

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, _ := us.UrlShortener.Urls[shortKey]

	log.Println(originalURL, shortKey, us.UrlShortener.Urls)

	// Redirect the user to the original URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
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
	shortKey := internal.GenerateShortKey()
	us.UrlShortener.Urls[shortKey] = link

	log.Println(link)

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	log.Println(shortenedURL)

}
