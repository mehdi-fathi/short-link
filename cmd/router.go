package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"short-link/cmd/rest"
)

func SetupRouter(handler *rest.Handler) *gin.Engine {

	router := gin.Default()

	router.LoadHTMLGlob("tmp/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.POST("/make", handler.HandleShorten)
	router.GET("/short/:url", handler.HandleRedirect)
	router.GET("/list/all", handler.HandleList)

	v1 := router.Group("/v1")
	{
		v1.GET("/list/all", handler.HandleListJson)

	}

	return router
}
