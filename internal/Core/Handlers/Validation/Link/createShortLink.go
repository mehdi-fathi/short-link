package Link

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Example struct for form data
type From struct {
	Link string `form:"link" binding:"required"`
}

// Validation middleware
func CreateShortLinkReqValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData From

		// Bind incoming form data to requestData struct and validate
		if err := c.ShouldBind(&formData); err != nil {
			// Set flash message
			session := sessions.Default(c)
			session.Set("error_msg", err.Error())
			session.Save()
			c.Redirect(http.StatusFound, "/index")
			c.Abort()
			return
		}

		// Proceed if validation passes
		c.Next()
	}
}
