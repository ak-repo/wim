package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/ak-repo/wim/pkg/response"
	"github.com/ak-repo/wim/pkg/utils"
	"github.com/go-chi/chi"
)

type WarehouseHandler struct {
	services *service.Services
}

func NewWarehouseHandler(services *service.Services) *WarehouseHandler {
	return &WarehouseHandler{services: services}
}

func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var req model.WarehouseRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	id, err := h.services.Warehouse.CreateWarehouse(r.Context(), &req)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *WarehouseHandler) GetWarehouseByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	warehouse, err := h.services.Warehouse.GetWarehouseByID(r.Context(), id)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, warehouse)
}

func (h *WarehouseHandler) GetWarehouseByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		response.WriteError(w, http.StatusBadRequest, apperrors.CodeInvalidInput, "code is required")
		return
	}

	warehouse, err := h.services.Warehouse.GetWarehouseByCode(r.Context(), code)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, warehouse)
}

func (h *WarehouseHandler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.WarehouseRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.Warehouse.UpdateWarehouse(r.Context(), id, &req); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "warehouse updated",
	})
}

func (h *WarehouseHandler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.Warehouse.DeleteWarehouse(r.Context(), id); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "warehouse deleted",
	})
}

func (h *WarehouseHandler) ListWarehouses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := &model.WarehouseParams{
		Page:   utils.GetInt(query, "page", 1),
		Limit:  utils.GetInt(query, "limit", 10),
		Active: utils.GetBoolPtr(query, "active"),
	}

	data, count, err := h.services.Warehouse.ListWarehouses(r.Context(), params)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	totalPage := (count + params.Limit - 1) / params.Limit
	responseDTO := map[string]any{
		"data":         data,
		"total_count":  count,
		"total_page":   totalPage,
		"current_page": params.Page,
		"limit":        params.Limit,
	}
	response.WriteJSON(w, http.StatusOK, responseDTO)
}
