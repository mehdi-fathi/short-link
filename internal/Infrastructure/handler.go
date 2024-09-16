package Infrastructure

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"short-link/internal/Core/Handlers/Http/rest"
	"short-link/internal/Core/Handlers/Http/web"
	service_interface "short-link/internal/Core/Ports"
	"short-link/pkg/logger"
	"time"
)

type Handler struct {
	HTTPServer *http.Server
}

// CreateHandler Creates a new instance of REST handler
func CreateHandler(linkService service_interface.LinkServiceInterface) *rest.HandlerRest {
	return &rest.HandlerRest{
		LinkService: linkService,
	}
}

// CreateHandler Creates a new instance of REST handler
func CreateHandlerMain(r *gin.Engine, defaultPort int) *Handler {

	const op = "http.rest.start"

	addr := fmt.Sprintf(":%d", defaultPort)

	HTTPServer := &http.Server{
		Addr:    addr,
		Handler: r,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	return &Handler{
		HTTPServer: HTTPServer,
	}
}

// CreateHandler Creates a new instance of REST handler
func CreateHandlerWeb(linkService service_interface.LinkServiceInterface) *web.HandlerWeb {
	return &web.HandlerWeb{
		LinkService: linkService,
	}
}

// Start starts the http server
func (h *Handler) Start() {

	const op = "http.rest.start"

	//h.Graylog.Infof("[OK] Starting HTTP REST Server on %s ", addr)
	err := h.HTTPServer.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.CreateLogError(errors.WithMessage(err, op).Error())
	}
	//// Code Reach Here after HTTP Server Shutdown!
	//h.Graylog.Info("[OK] HTTP REST Server is shutting down!")
}

// Stop handles the http server in graceful shutdown
func (h *Handler) Stop() {
	const op = "http.rest.stop"

	// Create an 5s timeout context or waiting for app to shutdown after 5 seconds
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	h.HTTPServer.SetKeepAlivesEnabled(false)
	if err := h.HTTPServer.Shutdown(ctxTimeout); err != nil {
		logger.CreateLogError(errors.WithMessage(err, op).Error())
	}
	logger.CreateLogInfo("HTTP REST Server graceful shutdown completed")

}
