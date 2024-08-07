package Infrastructure

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "short-link/internal/Config"
	"short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Handlers/Http/web"
	"short-link/internal/Core/Handlers/Validation/Link"
)


func SetupRouter(handler *rest.HandlerRest, handlerWeb *web.HandlerWeb) *gin.Engine {

	router := gin.Default()

	// Setup the session store
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("tmp/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", handlerWeb.HandleIndex)

	router.POST("/make", Link.ValidationMiddleware(), handlerWeb.HandleShorten)
	router.GET("/short/:url", handlerWeb.HandleRedirect)
	router.GET("/list/all", handlerWeb.HandleList)

	v1 := router.Group("/v1")
	{
		v1.GET("/list/all", handler.HandleListJson)

	}

	return router
}
