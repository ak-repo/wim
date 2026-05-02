package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/response"
	"github.com/ak-repo/wim/pkg/utils"
)

type CustomerTypeHandler struct {
	services *service.Services
}

func NewCustomerTypeHandler(services *service.Services) *CustomerTypeHandler {
	return &CustomerTypeHandler{services: services}
}

func (h *CustomerTypeHandler) CreateCustomerType(w http.ResponseWriter, r *http.Request) {
	var req model.CustomerTypeRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	id, err := h.services.CustomerType.CreateCustomerType(r.Context(), &req)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *CustomerTypeHandler) GetCustomerTypeByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	customerType, err := h.services.CustomerType.GetCustomerTypeByID(r.Context(), id)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, customerType)
}

func (h *CustomerTypeHandler) UpdateCustomerType(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.CustomerTypeRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.CustomerType.UpdateCustomerType(r.Context(), id, &req); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "customer type updated"})
}

func (h *CustomerTypeHandler) DeleteCustomerType(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.CustomerType.DeleteCustomerType(r.Context(), id); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "customer type deleted"})
}

func (h *CustomerTypeHandler) ListCustomerTypes(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := &model.CustomerTypeParams{
		Page:   utils.GetInt(query, "page", 1),
		Limit:  utils.GetInt(query, "limit", 10),
		Active: utils.GetBoolPtr(query, "active"),
	}

	data, count, err := h.services.CustomerType.ListCustomerTypes(r.Context(), params)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	totalPage := (count + params.Limit - 1) / params.Limit
	response.WriteJSON(w, http.StatusOK, map[string]any{
		"data":         data,
		"total_count":  count,
		"total_page":   totalPage,
		"current_page": params.Page,
		"limit":        params.Limit,
	})
}
