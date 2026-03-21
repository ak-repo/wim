package inventory

import (
	"net/http"
	"strconv"

	invt "github.com/ak-repo/wim/internal/service/inventory"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *invt.Service
}

func NewHandler(s *invt.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetByWarehouse(c *gin.Context) {
	warehouseID := c.Param("warehouse_id")
	limit := 50
	offset := 0

	inventory, err := h.service.GetByWarehouse(c.Request.Context(), warehouseID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *Handler) GetByProduct(c *gin.Context) {
	productID := c.Param("product_id")

	inventory, err := h.service.GetByProduct(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *Handler) List(c *gin.Context) {
	limit := 50
	offset := 0

	if v := c.Query("limit"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		limit = parsed
	}

	if v := c.Query("offset"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
		offset = parsed
	}

	inventory, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *Handler) Adjust(c *gin.Context) {
	var input invt.AdjustInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.AdjustInventory(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}
