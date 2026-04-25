package httpx

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/observability"
)

var ExposeStack bool

type ErrorResponse struct {
	Error errs.ApiError `json:"error"`
}

func Wrap(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		WriteError(w, r, err)
	}
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	status, apiErr, _ := errs.HTTPErrorResponse(err, ExposeStack)

	stack := errs.OpStack(err)

	observability.Report(r.Context(), err, stack)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: apiErr})
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}
