package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/response"
)

type UserHandler struct {
	services *service.Services
}

func NewUserHandler(services *service.Services) *UserHandler {
	return &UserHandler{services: services}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	data, err := h.services.User.ListUsers(r.Context())
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}
	response.WriteJSON(w, http.StatusOK, data)

}
