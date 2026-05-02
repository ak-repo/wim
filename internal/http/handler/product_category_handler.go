package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/response"
	"github.com/ak-repo/wim/pkg/utils"
)

type ProductCategoryHandler struct {
	services *service.Services
}

func NewProductCategoryHandler(services *service.Services) *ProductCategoryHandler {
	return &ProductCategoryHandler{services: services}
}

func (h *ProductCategoryHandler) CreateProductCategory(w http.ResponseWriter, r *http.Request) {
	var req model.ProductCategoryRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	id, err := h.services.ProductCategory.CreateProductCategory(r.Context(), &req)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *ProductCategoryHandler) GetProductCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	productCategory, err := h.services.ProductCategory.GetProductCategoryByID(r.Context(), id)
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, productCategory)
}

func (h *ProductCategoryHandler) UpdateProductCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	var req model.ProductCategoryRequest
	if ok := utils.DecodeJSON(w, r, &req); !ok {
		return
	}

	if err := h.services.ProductCategory.UpdateProductCategory(r.Context(), id, &req); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "product category updated"})
}

func (h *ProductCategoryHandler) DeleteProductCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := utils.ParseID(w, r)
	if !ok {
		return
	}

	if err := h.services.ProductCategory.DeleteProductCategory(r.Context(), id); err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "product category deleted"})
}

func (h *ProductCategoryHandler) ListProductCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	params := &model.ProductCategoryParams{
		Page:   utils.GetInt(query, "page", 1),
		Limit:  utils.GetInt(query, "limit", 10),
		Active: utils.GetBoolPtr(query, "active"),
	}

	data, count, err := h.services.ProductCategory.ListProductCategories(r.Context(), params)
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
