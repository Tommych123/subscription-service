package api

import (
	_ "github.com/Tommych123/subscription-service/internal/docs"
	"github.com/Tommych123/subscription-service/models"
	"github.com/Tommych123/subscription-service/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type SubscriptionHandler struct {
	svc    *service.SubscriptionService
	logger *zap.Logger
}

func NewSubscriptionHandler(svc *service.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, logger: logger}
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

// Create подписку
// @Summary Создать подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Подписка"
// @Success 201 {object} map[string]string "id новой подписки"
// @Failure 400 {object} map[string]string "Ошибка валидации входных данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/ [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		h.logger.Warn("Invalid input for Create", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}
	id, err := h.svc.Create(&sub)
	if err != nil {
		h.logger.Error("Failed to create subscription", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create: " + err.Error()})
		return
	}
	h.logger.Info("Subscription created", zap.String("id", id))
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetByID получить подписку по ID
// @Summary Получить подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	sub, err := h.svc.GetByID(id)
	if err != nil {
		h.logger.Error("Failed to get subscription by ID", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub == nil {
		h.logger.Info("Subscription not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Update обновить подписку
// @Summary Обновить подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param subscription body model.Subscription true "Подписка"
// @Success 200 "Обновление успешно"
// @Failure 400 {object} map[string]string "Ошибка валидации входных данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		h.logger.Warn("Invalid input for Update", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub.ID = id
	if err := h.svc.Update(&sub); err != nil {
		h.logger.Error("Failed to update subscription", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Subscription updated", zap.String("id", id))
	c.Status(http.StatusOK)
}

// Delete удалить подписку
// @Summary Удалить подписку
// @Tags subscriptions
// @Param id path string true "ID подписки"
// @Success 204 "Удаление успешно"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(id); err != nil {
		h.logger.Error("Failed to delete subscription", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Subscription deleted", zap.String("id", id))
	c.Status(http.StatusNoContent)
}

// List получить список подписок
// @Summary Получить список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/ [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	subs, err := h.svc.List()
	if err != nil {
		h.logger.Error("Failed to list subscriptions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// GetTotalCost вычислить суммарную стоимость подписок
// @Summary Получить суммарную стоимость подписок за период
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Param from query string true "Дата начала периода (MM-YYYY)" example(01-2023)
// @Param to query string true "Дата окончания периода (MM-YYYY)" example(12-2023)
// @Success 200 {object} map[string]int "Суммарная стоимость"
// @Failure 400 {object} map[string]string "Ошибка валидации входных данных"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /total [get]
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")
	fromStr := c.Query("from")
	toStr := c.Query("to")

	from, err := parseMonthYear(fromStr)
	if err != nil {
		h.logger.Warn("Invalid from date format", zap.String("from", fromStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format"})
		return
	}
	to, err := parseMonthYear(toStr)
	if err != nil {
		h.logger.Warn("Invalid to date format", zap.String("to", toStr), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format"})
		return
	}

	total, err := h.svc.GetTotalCost(userID, serviceName, from, to)
	if err != nil {
		h.logger.Error("Failed to get total cost", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Total cost calculated", zap.String("user_id", userID), zap.String("service_name", serviceName), zap.Int("total_cost", total))
	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}
