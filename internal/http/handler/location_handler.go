package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/httpx"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/utils"
	"github.com/go-chi/chi"
)

const opLocation errs.Op = "handler/LocationHandler"

type LocationHandler struct {
	services *service.Services
}

func NewLocationHandler(services *service.Services) *LocationHandler {
	return &LocationHandler{services: services}
}

func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var req model.LocationRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	id, err := h.services.Location.CreateLocation(r.Context(), &req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *LocationHandler) GetLocationByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	location, err := h.services.Location.GetLocationByID(r.Context(), id)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, location)
}

func (h *LocationHandler) GetLocationByCode(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		httpx.WriteError(w, r, errs.E(opLocation+".GetLocationByCode", errs.InvalidRequest, errors.New("code is required")))
		return
	}

	location, err := h.services.Location.GetLocationByCode(r.Context(), code)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, location)
}

func (h *LocationHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.LocationRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.Location.UpdateLocation(r.Context(), id, &req); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "location updated",
	})
}

func (h *LocationHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.Location.DeleteLocation(r.Context(), id); err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "location deleted",
	})
}

func (h *LocationHandler) ListLocations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := &model.LocationParams{
		Page:        utils.GetInt(query, "page", 1),
		Limit:       utils.GetInt(query, "limit", 10),
		Active:      utils.GetBoolPtr(query, "active"),
		WarehouseID: utils.GetInt(query, "warehouseId", 0),
		Zone:        utils.GetString(query, "zone", ""),
	}

	data, count, err := h.services.Location.ListLocations(r.Context(), params)
	if err != nil {
		httpx.WriteError(w, r, err)
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
	httpx.WriteJSON(w, http.StatusOK, responseDTO)
}

func (h *LocationHandler) ListLocationsByWarehouse(w http.ResponseWriter, r *http.Request) {
	warehouseIDStr := chi.URLParam(r, "warehouseId")
	if warehouseIDStr == "" {
		httpx.WriteError(w, r, errs.E(opLocation+".ListLocationsByWarehouse", errs.InvalidRequest, errors.New("warehouseId is required")))
		return
	}

	warehouseID, err := strconv.Atoi(warehouseIDStr)
	if err != nil {
		httpx.WriteError(w, r, errs.E(opLocation+".ListLocationsByWarehouse", errs.InvalidRequest, errors.New("invalid warehouse id")))
		return
	}

	data, err := h.services.Location.ListLocationsByWarehouse(r.Context(), warehouseID)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, data)
}
