package handler

import (
	"net/http"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/service"
	"github.com/ak-repo/wim/pkg/response"
)

type DashboardHandler struct {
	services *service.Services
}

func NewDashboardHandler(services *service.Services) *DashboardHandler {
	return &DashboardHandler{services: services}
}

func (h *DashboardHandler) TotalCounts(w http.ResponseWriter, r *http.Request) {

	resp := model.TotalCount{}
	var err error
	
	// Get total products count
	resp.TotalProducts, err = h.services.Product.GetProductCount(r.Context(), &model.ProductParams{})
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	// Get total users count
	resp.TotalUsers, err = h.services.User.GetUserCount(r.Context(), &model.UserParams{})
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	// Get total warehouses count
	resp.TotalWarehouses, err = h.services.Warehouse.GetWarehouseCount(r.Context(), &model.WarehouseParams{})
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	// Get total locations count
	resp.TotalLocations, err = h.services.Location.GetLocationCount(r.Context(), &model.LocationParams{})
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, resp)
}
