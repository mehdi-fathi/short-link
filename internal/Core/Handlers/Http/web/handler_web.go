package web

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"short-link/internal/Config"
	_ "short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Logic/Db/Serialization"
	service_interface "short-link/internal/Core/Ports"
	"short-link/pkg/errorMsg"
	"time"
)

var (
	// Histogram to track the duration of HTTP requests
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets, // Default bucket sizes
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestDuration)

}

type HandlerWeb struct {
	LinkService service_interface.LinkServiceInterface
}

func (h *HandlerWeb) HandleIndex(c *gin.Context) {
	start := time.Now()

	defer func() {
		duration := time.Since(start).Seconds()
		HttpRequestDuration.WithLabelValues("HandleIndex").Observe(duration)
	}()
	errorMsg := errorMsg.GetErrorMsg(c, "error_msg")

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main website",
		"error": errorMsg,
		"url":   Config.GetBaseUrl(),
	})

}

func (h *HandlerWeb) HandleShorten(c *gin.Context) {

	start := time.Now()

	defer func() {
		duration := time.Since(start).Seconds()
		HttpRequestDuration.WithLabelValues("HandleShorten").Observe(duration)
	}()

	link := c.PostForm("link")

	// Generate a unique shortened key for the original URL
	h.LinkService.SetUrl(link)

	c.Redirect(http.StatusMovedPermanently, Config.GetBaseUrl()+"/list/all")

}

func (h *HandlerWeb) HandleRedirect(c *gin.Context) {

	start := time.Now()

	defer func() {
		duration := time.Since(start).Seconds()
		HttpRequestDuration.WithLabelValues("HandleRedirect").Observe(duration)
	}()

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

	start := time.Now()

	defer func() {
		duration := time.Since(start).Seconds()
		HttpRequestDuration.WithLabelValues("HandleList").Observe(duration)
	}()

	linksDb, _ := h.LinkService.GetAllUrlV2()

	dataLinkSerialized := Serialization.DeserializeAllLink(linksDb)

	c.HTML(http.StatusOK, "list.html", gin.H{
		"data": dataLinkSerialized,
		"url":  Config.GetBaseUrl(),
	})
}
