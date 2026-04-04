package response

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/ak-repo/wim/pkg/errors"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// --------------------
// SERVICE ERROR HANDLER
// --------------------
func WriteServiceError(w http.ResponseWriter, err error) {
	// Check if structured error
	if e := apperrors.GetError(err); e != nil {
		WriteJSON(w, statusFromCode(e.Code), ErrorResponse{
			Error: ErrorBody{
				Code:    e.Code,
				Message: e.Message,
			},
		})
		return
	}

	// fallback for sentinel errors (fixed client messages; never echo err.Error())
	switch {
	case err == apperrors.ErrInvalidInput:
		WriteError(w, http.StatusBadRequest, apperrors.CodeInvalidInput, "invalid input")

	case err == apperrors.ErrAlreadyExists:
		WriteError(w, http.StatusConflict, apperrors.CodeAlreadyExists, "resource already exists")

	case err == apperrors.ErrUnauthorized:
		WriteError(w, http.StatusUnauthorized, apperrors.CodeUnauthorized, "unauthorized")

	case err == apperrors.ErrForbidden:
		WriteError(w, http.StatusForbidden, apperrors.CodeForbidden, "forbidden")

	case err == apperrors.ErrNotFound:
		WriteError(w, http.StatusNotFound, apperrors.CodeNotFound, "resource not found")

	case err == apperrors.ErrDatabase:
		WriteError(w, http.StatusInternalServerError, apperrors.CodeDatabase, "database error")

	case err == apperrors.ErrCheckingFailed:
		WriteError(w, http.StatusInternalServerError, apperrors.CodeCheckFailed, "failed to verify existing record")

	case err == apperrors.ErrRefCodeGeneration:
		WriteError(w, http.StatusInternalServerError, apperrors.CodeRefCodeFailed, "failed to generate reference code")

	case err == apperrors.ErrInternal:
		WriteError(w, http.StatusInternalServerError, apperrors.CodeInternal, "internal server error")

	default:
		// 🚨 never expose raw error
		WriteError(w, http.StatusInternalServerError, apperrors.CodeInternal, "internal server error")
	}
}

// --------------------
// GENERIC ERROR
// --------------------
func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

// --------------------
// JSON RESPONSE
// --------------------
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// --------------------
// STATUS MAPPING
// --------------------
func statusFromCode(code string) int {
	switch code {
	case apperrors.CodeInvalidInput:
		return http.StatusBadRequest
	case apperrors.CodeAlreadyExists:
		return http.StatusConflict
	case apperrors.CodeUnauthorized:
		return http.StatusUnauthorized
	case apperrors.CodeForbidden:
		return http.StatusForbidden
	case apperrors.CodeNotFound:
		return http.StatusNotFound
	case apperrors.CodeConcurrentUpdate:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
