package handler

import (
	"errors"
	"net/http"

	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/httpx"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/utils"
)

const opSalesOrder errs.Op = "handler/SalesOrderHandler"

type SalesOrderHandler struct {
	services *service.Services
}

func NewSalesOrderHandler(services *service.Services) *SalesOrderHandler {
	return &SalesOrderHandler{services: services}
}

// CreateSalesOrder handles the creation of a new sales order
func (h *SalesOrderHandler) CreateSalesOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreateSalesOrderRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	// Get user from context if available
	var createdBy *int
	// TODO: Extract user ID from auth context

	order, err := h.services.SalesOrder.CreateSalesOrder(r.Context(), &req, createdBy)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, order)
}

// GetSalesOrderByID retrieves a sales order by ID
func (h *SalesOrderHandler) GetSalesOrderByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	order, err := h.services.SalesOrder.GetSalesOrderByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, order)
}

// GetSalesOrderByRefCode retrieves a sales order by reference code
func (h *SalesOrderHandler) GetSalesOrderByRefCode(w http.ResponseWriter, r *http.Request) {
	refCode := r.URL.Query().Get("refCode")
	if refCode == "" {
		httpx.WriteError(w, r, errs.E(opSalesOrder+".GetSalesOrderByRefCode", errs.InvalidRequest, errors.New("refCode is required")))
		return
	}

	order, err := h.services.SalesOrder.GetSalesOrderByRefCode(r.Context(), refCode)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, order)
}

// ListSalesOrders lists sales orders with optional filtering
func (h *SalesOrderHandler) ListSalesOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := &model.SalesOrderParams{
		Page:             utils.GetInt(query, "page", 1),
		Limit:            utils.GetInt(query, "limit", 10),
		CustomerID:       utils.GetIntPtr(query, "customerId"),
		WarehouseID:      utils.GetIntPtr(query, "warehouseId"),
		Status:           utils.GetStringPtr(query, "status"),
		AllocationStatus: utils.GetStringPtr(query, "allocationStatus"),
	}

	data, count, err := h.services.SalesOrder.ListSalesOrders(r.Context(), params)
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

// UpdateSalesOrder updates an existing sales order
func (h *SalesOrderHandler) UpdateSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.CreateSalesOrderRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	order, err := h.services.SalesOrder.UpdateSalesOrder(r.Context(), id, &req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, order)
}

// CancelSalesOrder cancels a sales order
func (h *SalesOrderHandler) CancelSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.SalesOrder.CancelSalesOrder(r.Context(), id); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "sales order cancelled",
	})
}

// AllocateSalesOrder allocates stock for a sales order
func (h *SalesOrderHandler) AllocateSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.AllocateSalesOrderRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	// Get user from context if available
	var performedBy *int
	// TODO: Extract user ID from auth context

	if err := h.services.SalesOrder.AllocateSalesOrder(r.Context(), id, performedBy); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "sales order allocated",
	})
}

// DeallocateSalesOrder deallocates stock from a sales order
func (h *SalesOrderHandler) DeallocateSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.SalesOrder.DeallocateSalesOrder(r.Context(), id); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "sales order deallocated",
	})
}

// ShipSalesOrder ships a sales order
func (h *SalesOrderHandler) ShipSalesOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.ShipSalesOrderRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	// Get user from context if available
	var performedBy *int
	// TODO: Extract user ID from auth context

	if err := h.services.SalesOrder.ShipSalesOrder(r.Context(), id, &req, performedBy); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "sales order shipped",
	})
}
