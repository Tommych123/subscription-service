package handler

import (
	"net/http"
	"time"

	"github.com/Tommych123/subscription-service/internal/model"
	"github.com/Tommych123/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
}

func NewSubscriptionHandler(svc *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
}

func (h *SubscriptionHandler) RegisterRoutes(r *gin.Engine) {
	sub := r.Group("/subscriptions")
	{
		sub.POST("/", h.Create)
		sub.GET("/", h.List)
		sub.GET("/:id", h.GetByID)
		sub.PUT("/:id", h.Update)
		sub.DELETE("/:id", h.Delete)
	}
	r.GET("/total", h.GetTotalCost)
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	id, err := h.svc.Create(&sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	sub, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub.ID = id
	if err := h.svc.Update(&sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	subs, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	fromStr := c.Query("from")
	toStr := c.Query("to")
	from, err := parseMonthYear(fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format"})
		return
	}
	to, err := parseMonthYear(toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format"})
		return
	}
	total, err := h.svc.GetTotalCost(userID, serviceName, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}
