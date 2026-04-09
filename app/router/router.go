package router

import (
	"mini/config"
	"mini/database/ent"
	"net/http"

	ticketHandler "mini/app/ticket/handler"
	ticketRepo "mini/app/ticket/repository"
	ticketUC "mini/app/ticket/usecase"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	DB *ent.Client
	R  *gin.Engine
}

func (h *Handlers) Routes() {
	h.R.GET("/health", h.CheckHealth)

	v1 := h.R.Group(config.App.PrefixApi)
	v1.GET("/check-connection", h.CheckConnection)
	v1.GET("/version", h.Version)

	// Ticket
	repo := ticketRepo.NewTicketRepository()
	uc := ticketUC.NewTicketUsecase(h.DB, repo)
	ticketHandler.TicketRoute(uc, v1)
}

func (h *Handlers) CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "API is running"})
}

func (h *Handlers) CheckConnection(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Connected"})
}

func (h *Handlers) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "0.0.1"})
}
