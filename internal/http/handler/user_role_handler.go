package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/response"
	"github.com/ak-repo/wim/pkg/utils"
)

type UserRoleHandler struct {
	services *service.Services
}

func NewUserRoleHandler(services *service.Services) *UserRoleHandler {
	return &UserRoleHandler{services: services}
}

func (h *UserRoleHandler) CreateUserRole(w http.ResponseWriter, r *http.Request) {
	var req model.UserRoleRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	id, err := h.services.UserRole.CreateUserRole(r.Context(), &req)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *UserRoleHandler) GetUserRoleByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	userRole, err := h.services.UserRole.GetUserRoleByID(r.Context(), id)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, userRole)
}

func (h *UserRoleHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.UserRoleRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.UserRole.UpdateUserRole(r.Context(), id, &req); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "user role updated"})
}

func (h *UserRoleHandler) DeleteUserRole(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.UserRole.DeleteUserRole(r.Context(), id); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "user role deleted"})
}

func (h *UserRoleHandler) ListUserRoles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := &model.UserRoleParams{
		Page:   utils.GetInt(query, "page", 1),
		Limit:  utils.GetInt(query, "limit", 10),
		Active: utils.GetBoolPtr(query, "active"),
	}

	data, count, err := h.services.UserRole.ListUserRoles(r.Context(), params)
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
