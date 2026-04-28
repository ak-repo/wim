package handler

import (
	"errors"
	"net/http"

	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/httpx"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/auth"
	"github.com/ak-repo/wim/pkg/utils"
)

const opPurchaseOrder errs.Op = "handler/PurchaseOrderHandler"

type PurchaseOrderHandler struct {
	services *service.Services
}

func NewPurchaseOrderHandler(services *service.Services) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{services: services}
}

func (h *PurchaseOrderHandler) CreatePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePurchaseOrderRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	createdBy, _ := auth.UserIDFromContext(r.Context())
	var createdByPtr *int
	if createdBy > 0 {
		createdByPtr = &createdBy
	}

	order, err := h.services.PurchaseOrder.CreatePurchaseOrder(r.Context(), &req, createdByPtr)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, order)
}

func (h *PurchaseOrderHandler) GetPurchaseOrderByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	order, err := h.services.PurchaseOrder.GetPurchaseOrderByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, order)
}

func (h *PurchaseOrderHandler) GetPurchaseOrderByRefCode(w http.ResponseWriter, r *http.Request) {
	refCode := r.URL.Query().Get("refCode")
	if refCode == "" {
		httpx.WriteError(w, r, errs.E(opPurchaseOrder+".GetPurchaseOrderByRefCode", errs.InvalidRequest, errors.New("refCode is required")))
		return
	}

	order, err := h.services.PurchaseOrder.GetPurchaseOrderByRefCode(r.Context(), refCode)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, order)
}

func (h *PurchaseOrderHandler) ListPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := &model.PurchaseOrderParams{
		Page:        utils.GetInt(query, "page", 1),
		Limit:       utils.GetInt(query, "limit", 10),
		SupplierID:  utils.GetIntPtr(query, "supplierId"),
		WarehouseID: utils.GetIntPtr(query, "warehouseId"),
		Status:      utils.GetStringPtr(query, "status"),
	}

	data, count, err := h.services.PurchaseOrder.ListPurchaseOrders(r.Context(), params)
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

func (h *PurchaseOrderHandler) ReceivePurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.ReceivePurchaseOrderRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	performedBy, _ := auth.UserIDFromContext(r.Context())
	var performedByPtr *int
	if performedBy > 0 {
		performedByPtr = &performedBy
	}

	if err := h.services.PurchaseOrder.ReceivePurchaseOrder(r.Context(), id, &req, performedByPtr); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "purchase order received"})
}

func (h *PurchaseOrderHandler) PutAwayPurchaseOrder(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.PutAwayPurchaseOrderRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	performedBy, _ := auth.UserIDFromContext(r.Context())
	var performedByPtr *int
	if performedBy > 0 {
		performedByPtr = &performedBy
	}

	if err := h.services.PurchaseOrder.PutAwayPurchaseOrder(r.Context(), id, &req, performedByPtr); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "purchase order put away"})
}
