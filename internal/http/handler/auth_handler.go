package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/httpx"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/utils"
)

type AuthHandler struct {
	services *service.Services
}

func NewAuthHandler(services *service.Services) *AuthHandler {
	return &AuthHandler{services: services}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	err := h.services.Auth.Register(r.Context(), &req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "user registered",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	data, err := h.services.Auth.Login(r.Context(), &req)
	if err != nil {
		httpx.WriteError(w, r, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, data)
}
