package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type URLShortener struct {
	urls map[string]string
}

func main() {

	shortener := &URLShortener{
		urls: make(map[string]string),
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

func (us *URLShortener) HandleRedirect(c *gin.Context) {

	shortKey := c.Param("url")

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, _ := us.urls[shortKey]

	log.Println(originalURL,"sd",shortKey, us.urls)


	// Redirect the user to the original URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (us *URLShortener) HandleShorten(c *gin.Context) {
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
	shortKey := generateShortKey()
	us.urls[shortKey] = link

	log.Println(link)

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	log.Println(shortenedURL)

	//// Render the HTML response with the shortened URL
	//w.Header().Set("Content-Type", "text/html")
	//responseHTML := fmt.Sprintf(`
    //    <h2>URL Shortener</h2>
    //    <p>Original URL: %s</p>
    //    <p>Shortened URL: <a href="%s">%s</a></p>
    //    <form method="post" action="/shorten">
    //        <input type="text" name="url" placeholder="Enter a URL">
    //        <input type="submit" value="Shorten">
    //    </form>
    //`, originalURL, shortenedURL, shortenedURL)
	//fmt.Fprintf(w, responseHTML)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}
