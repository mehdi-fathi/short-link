package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) HandleShorten(c *gin.Context) {

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
	shortKey := h.linkService.SetUrl(link)

	log.Println(link)

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	log.Println(shortenedURL)
}

func (h *Handler) HandleRedirect(c *gin.Context) {

	shortKey := c.Param("url")

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL := h.linkService.GetUrl(shortKey)

	log.Println(originalURL, shortKey)

	// Redirect the user to the original URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

func (h *Handler) HandleList(c *gin.Context) {

	allUrl := h.linkService.GetAllUrl()

	log.Println(allUrl, "hi")

	c.HTML(http.StatusOK, "list.html", gin.H{
		"data": allUrl,
	})
}
