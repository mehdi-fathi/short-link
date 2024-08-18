package errorMsg

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "log"
)

func GetErrorMsg(c *gin.Context, key string) interface{} {

	session := sessions.Default(c)
	errorMsg := session.Get(key)
	session.Delete(key)
	session.Save()

	return errorMsg
}
