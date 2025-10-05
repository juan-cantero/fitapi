package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juan-cantero/fitapi/internal/models"
	"github.com/juan-cantero/fitapi/internal/services"
)

// EquipmentHandler handles HTTP requests for equipment endpoints
type EquipmentHandler struct {
	service *services.EquipmentService
}

// NewEquipmentHandler creates a new equipment handler
func NewEquipmentHandler(service *services.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{service: service}
}

// Create handles POST /api/equipment
func (h *EquipmentHandler) Create(c *gin.Context) {
	var req models.CreateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	equipment, err := h.service.CreateEquipment(c.Request.Context(), userID, &req)
	if err != nil {
		// Log the actual error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create equipment",
			"detail": err.Error(), // Add this temporarily for debugging
		})
		return
	}

	c.JSON(http.StatusCreated, equipment)
}

// GetByID handles GET /api/equipment/:id
func (h *EquipmentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	equipment, err := h.service.GetEquipment(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, services.ErrEquipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "equipment not found"})
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to access this equipment"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get equipment"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// List handles GET /api/equipment
func (h *EquipmentHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	equipment, err := h.service.ListEquipment(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list equipment"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// Update handles PUT /api/equipment/:id
func (h *EquipmentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.UpdateEquipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	equipment, err := h.service.UpdateEquipment(c.Request.Context(), id, userID, &req)
	if err != nil {
		if errors.Is(err, services.ErrEquipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "equipment not found"})
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to update this equipment"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update equipment"})
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// Delete handles DELETE /api/equipment/:id
func (h *EquipmentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	err := h.service.DeleteEquipment(c.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, services.ErrEquipmentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "equipment not found"})
			return
		}
		if errors.Is(err, services.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to delete this equipment"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete equipment"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
