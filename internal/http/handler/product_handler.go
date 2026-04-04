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

type ProductHandler struct {
	services *service.Services
}

func NewProductHandler(services *service.Services) *ProductHandler {
	return &ProductHandler{services: services}
}

// CREATE PRODUCT
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req model.ProductRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	id, err := h.services.Product.CreateProduct(r.Context(), &req)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]int{
		"id": id,
	})
}

// GET PRODUCT BY ID
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	product, err := h.services.Product.GetProductByID(r.Context(), id)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, product)
}

// GET PRODUCT BY SKU
func (h *ProductHandler) GetProductBySKU(w http.ResponseWriter, r *http.Request) {
	sku := chi.URLParam(r, "sku")
	if sku == "" {
		response.WriteError(w, http.StatusBadRequest, apperrors.CodeInvalidInput, "sku is required")
		return
	}

	product, err := h.services.Product.GetProductBySKU(r.Context(), sku)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, product)
}

// UPDATE PRODUCT (PATCH)
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.ProductRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.Product.UpdateProduct(r.Context(), id, &req); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "product updated",
	})
}

// DELETE PRODUCT
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.Product.DeleteProduct(r.Context(), id); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "product deleted",
	})
}

// LIST PRODUCTS (WITH FILTER + PAGINATION)
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := &model.ProductParams{
		Page:     utils.GetInt(query, "page", 1),
		Limit:    utils.GetInt(query, "limit", 10),
		Active:   utils.GetBoolPtr(query, "active"),
		Category: utils.GetString(query, "category", ""),
	}

	data, count, err := h.services.Product.ListProducts(r.Context(), params)
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
