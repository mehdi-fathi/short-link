package rest

import (
	"github.com/gin-gonic/gin"
	service_interface "short-link/internal/Core/Ports"
	"short-link/pkg/logger"
)

type HandlerRest struct {
	loggerInstance *logger.StandardLogger
	LinkService    service_interface.LinkServiceInterface
}

func (h *HandlerRest) HandleListJson(c *gin.Context) {

	all, _ := h.LinkService.GetAllLinkApi()

	c.JSON(200, all)
}
