package router

import (
	"mini/config"
	"mini/database/ent"
	"net/http"

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
