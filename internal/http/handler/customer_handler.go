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

type CustomerHandler struct {
	services *service.Services
}

func NewCustomerHandler(services *service.Services) *CustomerHandler {
	return &CustomerHandler{services: services}
}

func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var req model.CustomerRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	id, err := h.services.Customer.CreateCustomer(r.Context(), &req)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	customer, err := h.services.Customer.GetCustomerByID(r.Context(), id)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) GetCustomerByEmail(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		response.WriteError(w, http.StatusBadRequest, apperrors.CodeInvalidInput, "email is required")
		return
	}

	customer, err := h.services.Customer.GetCustomerByEmail(r.Context(), email)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.CustomerRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.Customer.UpdateCustomer(r.Context(), id, &req); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "customer updated"})
}

func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.Customer.DeleteCustomer(r.Context(), id); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "customer deleted"})
}

func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := &model.CustomerParams{
		Page:   utils.GetInt(query, "page", 1),
		Limit:  utils.GetInt(query, "limit", 10),
		Active: utils.GetBoolPtr(query, "active"),
	}

	data, count, err := h.services.Customer.ListCustomers(r.Context(), params)
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
