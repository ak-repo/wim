package response

import (
	"encoding/json"
	"errors"
	"net/http"

	apperrors "github.com/ak-repo/wim/pkg/errors"
)

func WriteServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrInvalidInput):
		WriteError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, apperrors.ErrAlreadyExists):
		WriteError(w, http.StatusConflict, err.Error())
	case errors.Is(err, apperrors.ErrUnauthorized):
		WriteError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, apperrors.ErrForbidden):
		WriteError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, apperrors.ErrNotFound):
		WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, apperrors.ErrCheckingFaild):
		WriteError(w, http.StatusNotFound, err.Error())

	default:
		WriteError(w, http.StatusInternalServerError, err.Error())
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
