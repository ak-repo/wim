package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/httpx"
	"github.com/go-chi/chi"
)

// JSON decoder
// decodeJSON safely decodes request body
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	defer r.Body.Close()

	// limit request size (1MB)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		origErr := err.Error()
		httpx.WriteError(w, r, errs.E("utils/DecodeJSON", errs.InvalidRequest, errors.New("invalid request body: "+origErr), errs.WithCode(errs.CodeInvalidRequest)))
		return false
	}

	return true
}

// parseID safely parses path ID
func ParseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		httpx.WriteError(w, r, errs.E("utils/ParseID", errs.InvalidRequest, errors.New("invalid id"), errs.WithCode(errs.CodeInvalidRequest)))
		return 0, false
	}
	return id, true
}

// Int
func GetInt(q url.Values, key string, def int) int {
	val := q.Get(key)
	if v, err := strconv.Atoi(val); err == nil && v > 0 {
		return v
	}
	return def
}

func GetString(q url.Values, key, def string) string {
	val := strings.TrimSpace(q.Get(key))
	if val == "" {
		return def
	}
	return val
}

func GetBool(q url.Values, key string, def bool) bool {
	val := q.Get(key)
	if val == "" {
		return def
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return b
}

// pointer values
func GetBoolPtr(q url.Values, key string) *bool {
	val := q.Get(key)
	if val == "" {
		return nil
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return nil
	}
	return &b
}

func GetIntPtr(q url.Values, key string) *int {
	val := q.Get(key)
	if val == "" {
		return nil
	}
	v, err := strconv.Atoi(val)
	if err != nil || v <= 0 {
		return nil
	}
	return &v
}

func GetStringPtr(q url.Values, key string) *string {
	val := strings.TrimSpace(q.Get(key))
	if val == "" {
		return nil
	}
	return &val
}
