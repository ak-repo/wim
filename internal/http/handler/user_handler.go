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

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	params := model.UserParams{
		Page:   utils.GetInt(query, "page", 1),
		Limit:  utils.GetInt(query, "limit", 10),
		Active: utils.GetBoolPtr(query, "active"),
	}

	data, count, err := h.services.User.ListUsers(r.Context(), &params)
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
