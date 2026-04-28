package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/httpx"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/auth"
	"github.com/ak-repo/wim/pkg/utils"
)

const opPicking errs.Op = "handler/PickingHandler"

type PickingHandler struct {
	services *service.Services
}

func NewPickingHandler(services *service.Services) *PickingHandler {
	return &PickingHandler{services: services}
}

func (h *PickingHandler) CreatePickingTask(w http.ResponseWriter, r *http.Request) {
	var req model.CreatePickingTaskRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	createdBy, _ := auth.UserIDFromContext(r.Context())
	var createdByPtr *int
	if createdBy > 0 {
		createdByPtr = &createdBy
	}

	task, err := h.services.Picking.CreatePickingTask(r.Context(), &req, createdByPtr)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, task)
}

func (h *PickingHandler) GetPickingTaskByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	task, err := h.services.Picking.GetPickingTaskByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, task)
}

func (h *PickingHandler) GetPickingTaskByRefCode(w http.ResponseWriter, r *http.Request) {
	refCode := r.URL.Query().Get("refCode")
	if refCode == "" {
		httpx.WriteError(w, r, errs.E(opPicking+".GetPickingTaskByRefCode", errs.InvalidRequest, errors.New("refCode is required")))
		return
	}

	task, err := h.services.Picking.GetPickingTaskByRefCode(r.Context(), refCode)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, task)
}

func (h *PickingHandler) ListPickingTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	params := &model.PickingTaskParams{
		Page:  utils.GetInt(query, "page", 1),
		Limit: utils.GetInt(query, "limit", 10),
	}

	if warehouseID := query.Get("warehouseId"); warehouseID != "" {
		if id, err := strconv.Atoi(warehouseID); err == nil {
			params.WarehouseID = &id
		}
	}
	if status := query.Get("status"); status != "" {
		params.Status = &status
	}
	if priority := query.Get("priority"); priority != "" {
		params.Priority = &priority
	}
	if assignedTo := query.Get("assignedTo"); assignedTo != "" {
		if id, err := strconv.Atoi(assignedTo); err == nil {
			params.AssignedTo = &id
		}
	}

	tasks, count, err := h.services.Picking.ListPickingTasks(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	totalPage := (count + params.Limit - 1) / params.Limit
	httpx.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":         tasks,
		"total_count":  count,
		"total_page":   totalPage,
		"current_page": params.Page,
		"limit":        params.Limit,
	})
}

func (h *PickingHandler) AssignPickingTask(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.AssignPickingTaskRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	err := h.services.Picking.AssignPickingTask(r.Context(), id, &req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "picking task assigned"})
}

func (h *PickingHandler) StartPickingTask(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	err := h.services.Picking.StartPickingTask(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "picking task started"})
}

func (h *PickingHandler) PickItem(w http.ResponseWriter, r *http.Request) {
	taskID, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.PickItemRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	performedBy, _ := auth.UserIDFromContext(r.Context())
	var performedByPtr *int
	if performedBy > 0 {
		performedByPtr = &performedBy
	}

	err := h.services.Picking.PickItem(r.Context(), taskID, &req, performedByPtr)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "item picked"})
}

func (h *PickingHandler) CompletePickingTask(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.CompletePickingRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	err := h.services.Picking.CompletePickingTask(r.Context(), id, req.Notes)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "picking task completed"})
}

func (h *PickingHandler) CancelPickingTask(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	err := h.services.Picking.CancelPickingTask(r.Context(), id, req.Notes)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{"message": "picking task cancelled"})
}