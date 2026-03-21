package audit

import (
	"net/http"
	"strconv"
	"strings"

	auditSvc "github.com/ak-repo/wim/internal/service/audit"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *auditSvc.Service
}

func NewHandler(s *auditSvc.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) List(c *gin.Context) {
	if userID := c.Query("user_id"); userID != "" {
		uuidVal, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}

		limit := 50
		if v := c.Query("limit"); v != "" {
			parsed, err := strconv.Atoi(v)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
				return
			}
			limit = parsed
		}

		logs, err := h.service.GetByUser(c.Request.Context(), uuidVal, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, logs)
		return
	}

	entityType := strings.TrimSpace(c.Query("entity_type"))
	entityID := c.Query("entity_id")
	if entityType == "" || entityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_type and entity_id are required when user_id is not provided"})
		return
	}

	uuidVal, err := uuid.Parse(entityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entity_id"})
		return
	}

	logs, err := h.service.GetByEntity(c.Request.Context(), entityType, uuidVal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
