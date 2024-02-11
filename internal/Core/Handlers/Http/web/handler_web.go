package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short-link/internal/Config"
	"short-link/internal/Core/Domin"
	_ "short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Logic/Db/serialization"
	service_interface "short-link/internal/Core/Ports"
	"short-link/pkg/logger"
)

type HandlerWeb struct {
	loggerInstance *logger.StandardLogger
	LinkService    service_interface.ServiceInterface
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
	originalURL := h.LinkService.GetUrl(shortKey)

	if originalURL != nil && originalURL.Status == Domin.LINK_STATUS_APPROVE {
		// Redirect the user to the original URL
		c.Redirect(http.StatusMovedPermanently, originalURL.Link)
	}

	c.HTML(http.StatusNotFound, "404.html", nil)

}

func (h *HandlerWeb) HandleList(c *gin.Context) {

	linksDb, _ := h.LinkService.GetAllUrlV2()

	dataLinkSerialized := serialization.DeserializeAllLink(linksDb)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"data": dataLinkSerialized,
		"url":  Config.GetBaseUrl(),
	})
}
