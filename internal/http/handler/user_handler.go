package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/response"
	"github.com/ak-repo/wim/pkg/utils"
)

type UserHandler struct {
	services *service.Services
}

func NewUserHandler(services *service.Services) *UserHandler {
	return &UserHandler{services: services}
}

// CREATE USER
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.UserRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	id, err := h.services.User.CreateUser(r.Context(), &req)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]int{
		"id": id,
	})
}

// GET USER BY ID
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	user, err := h.services.User.GetUserByID(r.Context(), id)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// UPDATE USER (PATCH)
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.UserRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.User.UpdateUser(r.Context(), id, &req); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user updated",
	})
}

// DELETE USER
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.User.DeleteUser(r.Context(), id); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user deleted",
	})
}

// LIST USERS (WITH FILTER + PAGINATION)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := &model.UserParams{
		Page:   utils.GetInt(query, "page", 1),
		Limit:  utils.GetInt(query, "limit", 10),
		Active: utils.GetBoolPtr(query, "active"),
	}

	data, count, err := h.services.User.ListUsers(r.Context(), params)
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
