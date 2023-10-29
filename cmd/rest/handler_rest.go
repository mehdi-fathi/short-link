package rest

import (
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
	h.LinkService.SetUrl(link)

	c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/list/all")

}

func (h *Handler) HandleRedirect(c *gin.Context) {

	shortKey := c.Param("url")

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL := h.LinkService.GetUrl(shortKey)

	// Redirect the user to the original URL
	c.Redirect(http.StatusMovedPermanently, originalURL.Link)
}

func (h *Handler) HandleList(c *gin.Context) {

	//allUrl := h.LinkService.GetAllUrl()

	all, _ := h.LinkService.GetAllUrlV2()

	log.Println(all[0].ShortKey)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"data": all,
	})
}
