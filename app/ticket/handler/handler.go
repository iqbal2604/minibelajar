package handler

import (
	"net/http"

	ticketDomain "mini/app/ticket"
	ticketUC "mini/app/ticket/usecase"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	uc ticketUC.TicketUsecase
}

func TicketRoute(uc ticketUC.TicketUsecase, r *gin.RouterGroup) {
	h := &TicketHandler{uc: uc}

	tickets := r.Group("/tickets")
	tickets.POST("", h.Create)
	tickets.GET("", h.GetAll)
	tickets.GET("/:id", h.GetByID)
	tickets.PUT("/:id", h.Update)
	tickets.DELETE("/:id", h.Delete)
}

func (h *TicketHandler) Create(c *gin.Context) {
	var f ticketDomain.Form
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if err := f.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	result, err := h.uc.Create(c.Request.Context(), &f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": result})
}

func (h *TicketHandler) GetAll(c *gin.Context) {
	q := ticketDomain.NewQueryParam(c)

	result, total, err := h.uc.GetAll(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
		"meta": gin.H{
			"total": total,
			"page":  q.Page,
			"limit": q.Limit,
		},
	})
}

func (h *TicketHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func (h *TicketHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var f ticketDomain.UpdateForm
	if err := c.ShouldBindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	if err := f.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	result, err := h.uc.Update(c.Request.Context(), id, &f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": result})
}

func (h *TicketHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "ticket deleted"})
}
