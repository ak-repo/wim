package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/httpx"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/utils"
)

type InventoryHandler struct {
	services *service.Services
}

func NewInventoryHandler(services *service.Services) *InventoryHandler {
	return &InventoryHandler{services: services}
}

func (h *InventoryHandler) AdjustInventory(w http.ResponseWriter, r *http.Request) {
	var req model.AdjustInventoryRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.Inventory.AdjustInventory(r.Context(), &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "inventory adjusted",
	})
}

func (h *InventoryHandler) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	item, err := h.services.Inventory.GetInventoryByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, item)
}

func (h *InventoryHandler) ListInventory(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := &model.InventoryParams{
		Page:        utils.GetInt(query, "page", 1),
		Limit:       utils.GetInt(query, "limit", 10),
		ProductID:   utils.GetIntPtr(query, "productId"),
		WarehouseID: utils.GetIntPtr(query, "warehouseId"),
		LocationID:  utils.GetIntPtr(query, "locationId"),
		BatchID:     utils.GetIntPtr(query, "batchId"),
	}

	data, count, err := h.services.Inventory.ListInventory(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	totalPage := (count + params.Limit - 1) / params.Limit

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"data":         data,
		"total_count":  count,
		"total_page":   totalPage,
		"current_page": params.Page,
		"limit":        params.Limit,
	})
}

func (h *InventoryHandler) ListStockMovements(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := &model.StockMovementParams{
		Page:          utils.GetInt(query, "page", 1),
		Limit:         utils.GetInt(query, "limit", 10),
		MovementType:  utils.GetStringPtr(query, "movementType"),
		ProductID:     utils.GetIntPtr(query, "productId"),
		WarehouseID:   utils.GetIntPtr(query, "warehouseId"),
		LocationID:    utils.GetIntPtr(query, "locationId"),
		BatchID:       utils.GetIntPtr(query, "batchId"),
		ReferenceType: utils.GetStringPtr(query, "referenceType"),
		ReferenceID:   utils.GetIntPtr(query, "referenceId"),
	}

	data, count, err := h.services.Inventory.ListStockMovements(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	totalPage := (count + params.Limit - 1) / params.Limit

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"data":         data,
		"total_count":  count,
		"total_page":   totalPage,
		"current_page": params.Page,
		"limit":        params.Limit,
	})
}
