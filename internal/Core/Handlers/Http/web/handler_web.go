package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short-link/internal/Config"
	_ "short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Logic/Db/Serialization"
	service_interface "short-link/internal/Core/Ports"
	"short-link/pkg/errorMsg"
	"short-link/pkg/logger"
)

type HandlerWeb struct {
	loggerInstance *logger.StandardLogger
	LinkService    service_interface.LinkServiceInterface
}

func (h *HandlerWeb) HandleIndex(c *gin.Context) {

	errorMsg := errorMsg.GetErrorMsg(c, "error_msg")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
		"error": errorMsg,
		"url":   Config.GetBaseUrl(),
	})

}

func (h *HandlerWeb) HandleShorten(c *gin.Context) {

	link := c.PostForm("link")

	// Generate a unique shortened key for the original URL
	h.LinkService.SetUrl(link)

	c.Redirect(http.StatusMovedPermanently, Config.GetBaseUrl()+"/list/all")

}

func (h *HandlerWeb) HandleRedirect(c *gin.Context) {

	shortKey := c.Param("url")

	// Retrieve the original URL from the `urls` map using the shortened key
	link := h.LinkService.FindValidUrlByShortKey(shortKey)

	if link != nil {
		// Redirect the user to the original URL
		c.Redirect(http.StatusMovedPermanently, link.Link)
	}

	c.HTML(http.StatusNotFound, "404.html", nil)

}

func (h *HandlerWeb) HandleList(c *gin.Context) {

	linksDb, _ := h.LinkService.GetAllUrlV2()

	dataLinkSerialized := Serialization.DeserializeAllLink(linksDb)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"data": dataLinkSerialized,
		"url":  Config.GetBaseUrl(),
	})
}
