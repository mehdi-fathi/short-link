package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	service_interface "short-link/internal/interface"
	"time"
)

type Handler struct {
	VersionInfo struct {
		GitCommit string
		BuildTime string
		StartTime time.Time
	}
	HTTPServer  *http.Server
	LinkService service_interface.ServiceInterface
}

// CreateHandler Creates a new instance of REST handler
func CreateHandler(linkService service_interface.ServiceInterface) *Handler {
	return &Handler{
		LinkService: linkService,
	}
}

// Start starts the http server
func (h *Handler) Start(ctx context.Context, r *gin.Engine, defaultPort int) {
	const op = "http.rest.start"

	addr := fmt.Sprintf(":%d", defaultPort)

	h.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: r,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	//h.Logger.Infof("[OK] Starting HTTP REST Server on %s ", addr)
	err := h.HTTPServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Println(errors.WithMessage(err, op))
	}
	//// Code Reach Here after HTTP Server Shutdown!
	//h.Logger.Info("[OK] HTTP REST Server is shutting down!")
}

// Stop handles the http server in graceful shutdown
func (h *Handler) Stop() {
	const op = "http.rest.stop"

	// Create an 5s timeout context or waiting for app to shutdown after 5 seconds
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	h.HTTPServer.SetKeepAlivesEnabled(false)
	if err := h.HTTPServer.Shutdown(ctxTimeout); err != nil {
		log.Println(errors.WithMessage(err, op))
	}
	log.Println("HTTP REST Server graceful shutdown completed")

}
