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
	resp.TotalProducts, err = h.services.Product.GetProductCount(r.Context(), &model.ProductParams{})
	if err != nil {
		response.WriteServiceError(w, err)
		return
	}

	//TODO: like this implement other counts need adding the proper service functions .

	response.WriteJSON(w, http.StatusOK, resp)
}
