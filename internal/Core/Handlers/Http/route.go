package Http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short-link/internal/Config"
	"short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Handlers/Http/web"
)

func SetupRouter(handler *rest.HandlerRest, handlerWeb *web.HandlerWeb) *gin.Engine {

	router := gin.Default()

	router.LoadHTMLGlob("tmp/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
			"url":   Config.GetBaseUrl(),
		})
	})

	router.POST("/make", handlerWeb.HandleShorten)
	router.GET("/short/:url", handlerWeb.HandleRedirect)
	router.GET("/list/all", handlerWeb.HandleList)

	v1 := router.Group("/v1")
	{
		v1.GET("/list/all", handler.HandleListJson)

	}

	return router
}
