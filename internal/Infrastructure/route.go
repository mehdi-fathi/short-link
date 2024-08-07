package Infrastructure

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "short-link/internal/Config"
	"short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Handlers/Http/web"
)


// Example struct for form data
type RequestData struct {
	Link  string `form:"link" binding:"required"`
}

// Validation middleware
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData RequestData

		// Bind incoming form data to requestData struct and validate
		if err := c.ShouldBind(&requestData); err != nil {
			// Set flash message
			session := sessions.Default(c)
			session.Set("error", err.Error())
			session.Save()
			c.Redirect(http.StatusFound, "/index")
			c.Abort()
			return
		}

		// Proceed if validation passes
		c.Next()
	}
}

func SetupRouter(handler *rest.HandlerRest, handlerWeb *web.HandlerWeb) *gin.Engine {

	router := gin.Default()

	// Setup the session store
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("tmp/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", handlerWeb.HandleIndex)

	router.POST("/make", ValidationMiddleware(), handlerWeb.HandleShorten)
	router.GET("/short/:url", handlerWeb.HandleRedirect)
	router.GET("/list/all", handlerWeb.HandleList)

	v1 := router.Group("/v1")
	{
		v1.GET("/list/all", handler.HandleListJson)

	}

	return router
}
