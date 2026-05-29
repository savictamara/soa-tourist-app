package handlers

import (
	"errors"
	"net/http"

	"purchase-service/models"
	"purchase-service/service"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	service *service.PurchaseService
}

func NewPurchaseHandler(service *service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

func (h *PurchaseHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *PurchaseHandler) GetCart(c *gin.Context) {
	cart, err := h.service.GetCart(c.Request.Context(), c.Param("touristId"))
	if err != nil {
		writeError(c, err, "failed to fetch cart")
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *PurchaseHandler) AddItem(c *gin.Context) {
	var item models.OrderItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	cart, err := h.service.AddItem(c.Request.Context(), c.Param("touristId"), item)
	if err != nil {
		writeError(c, err, "failed to add item")
		return
	}
	c.JSON(http.StatusCreated, cart)
}

func (h *PurchaseHandler) RemoveItem(c *gin.Context) {
	cart, err := h.service.RemoveItem(c.Request.Context(), c.Param("touristId"), c.Param("tourId"))
	if err != nil {
		writeError(c, err, "failed to remove item")
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (h *PurchaseHandler) Checkout(c *gin.Context) {
	tokens, err := h.service.Checkout(c.Request.Context(), c.Param("touristId"))
	if err != nil {
		writeError(c, err, "failed to checkout")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"tokens": tokens})
}

func (h *PurchaseHandler) GetTokens(c *gin.Context) {
	tokens, err := h.service.GetTokens(c.Request.Context(), c.Param("touristId"))
	if err != nil {
		writeError(c, err, "failed to fetch tokens")
		return
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *PurchaseHandler) IsPurchased(c *gin.Context) {
	purchased, err := h.service.IsPurchased(c.Request.Context(), c.Param("touristId"), c.Param("tourId"))
	if err != nil {
		writeError(c, err, "failed to check purchase")
		return
	}
	c.JSON(http.StatusOK, gin.H{"purchased": purchased})
}

func writeError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, service.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": fallback})
}
