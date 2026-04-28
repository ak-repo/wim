package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ak-repo/wim/internal/constants"
	apperrors "github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
)

type PurchaseOrderService interface {
	CreatePurchaseOrder(ctx context.Context, input *model.CreatePurchaseOrderRequest, createdBy *int) (*model.PurchaseOrderResponse, error)
	GetPurchaseOrderByID(ctx context.Context, orderID int) (*model.PurchaseOrderResponse, error)
	GetPurchaseOrderByRefCode(ctx context.Context, refCode string) (*model.PurchaseOrderResponse, error)
	ListPurchaseOrders(ctx context.Context, params *model.PurchaseOrderParams) ([]*model.PurchaseOrderResponse, int, error)
	ReceivePurchaseOrder(ctx context.Context, orderID int, input *model.ReceivePurchaseOrderRequest, performedBy *int) error
	PutAwayPurchaseOrder(ctx context.Context, orderID int, input *model.PutAwayPurchaseOrderRequest, performedBy *int) error
}

type purchaseOrderService struct {
	repos *repository.Repositories
}

func NewPurchaseOrderService(repositories *repository.Repositories) PurchaseOrderService {
	return &purchaseOrderService{repos: repositories}
}

func (s *purchaseOrderService) CreatePurchaseOrder(ctx context.Context, input *model.CreatePurchaseOrderRequest, createdBy *int) (*model.PurchaseOrderResponse, error) {
	if input == nil {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if input.SupplierID <= 0 {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "supplierId is required")
	}
	if input.WarehouseID <= 0 {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "warehouseId is required")
	}
	if len(input.Items) == 0 {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "at least one item is required")
	}

	if _, err := s.repos.Warehouse.GetByID(ctx, input.WarehouseID); err != nil {
		if errors.Is(err, repository.ErrWarehouseNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "warehouse not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load warehouse")
	}

	refCode, err := s.repos.RefCode.GeneratePurchaseOrderRefCode(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeRefCodeFailed, "failed to generate purchase order reference")
	}

	var items []*model.PurchaseOrderItemDTO
	for _, itemReq := range input.Items {
		if itemReq.ProductID <= 0 {
			return nil, apperrors.New(apperrors.CodeInvalidInput, "productId is required for all items")
		}
		if itemReq.QuantityOrdered <= 0 {
			return nil, apperrors.New(apperrors.CodeInvalidInput, "quantityOrdered must be greater than 0")
		}
		if _, err := s.repos.Product.GetByID(ctx, itemReq.ProductID); err != nil {
			if errors.Is(err, repository.ErrProductNotFound) {
				return nil, apperrors.New(apperrors.CodeNotFound, fmt.Sprintf("product not found: %d", itemReq.ProductID))
			}
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product")
		}

		item := &model.PurchaseOrderItemDTO{
			ProductID:       itemReq.ProductID,
			QuantityOrdered: itemReq.QuantityOrdered,
		}
		if strings.TrimSpace(itemReq.BatchNumber) != "" {
			batch := strings.TrimSpace(itemReq.BatchNumber)
			item.BatchNumber = &batch
		}
		if itemReq.UnitPrice != nil {
			item.UnitPrice = itemReq.UnitPrice
		}
		items = append(items, item)
	}

	order := &model.PurchaseOrderDTO{
		RefCode:      refCode,
		SupplierID:   input.SupplierID,
		WarehouseID:  input.WarehouseID,
		Status:       constants.StatusPending,
		ExpectedDate: input.ExpectedDate,
		CreatedBy:    createdBy,
	}
	if strings.TrimSpace(input.Notes) != "" {
		notes := strings.TrimSpace(input.Notes)
		order.Notes = &notes
	}

	createdOrder, err := s.repos.PurchaseOrder.Create(ctx, order, items)
	if err != nil {
		return nil, err
	}

	orderItems, err := s.repos.PurchaseOrder.GetItemsByOrderID(ctx, createdOrder.ID)
	if err != nil {
		return nil, err
	}

	response := createdOrder.ToAPIResponse()
	response.Items = orderItems.ToAPIResponse()
	return response, nil
}

func (s *purchaseOrderService) GetPurchaseOrderByID(ctx context.Context, orderID int) (*model.PurchaseOrderResponse, error) {
	order, err := s.repos.PurchaseOrder.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrPurchaseOrderNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "purchase order not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load purchase order")
	}
	items, err := s.repos.PurchaseOrder.GetItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	response := order.ToAPIResponse()
	response.Items = items.ToAPIResponse()
	return response, nil
}

func (s *purchaseOrderService) GetPurchaseOrderByRefCode(ctx context.Context, refCode string) (*model.PurchaseOrderResponse, error) {
	order, err := s.repos.PurchaseOrder.GetByRefCode(ctx, refCode)
	if err != nil {
		if errors.Is(err, repository.ErrPurchaseOrderNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "purchase order not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load purchase order")
	}
	items, err := s.repos.PurchaseOrder.GetItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	response := order.ToAPIResponse()
	response.Items = items.ToAPIResponse()
	return response, nil
}

func (s *purchaseOrderService) ListPurchaseOrders(ctx context.Context, params *model.PurchaseOrderParams) ([]*model.PurchaseOrderResponse, int, error) {
	if params == nil {
		params = &model.PurchaseOrderParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	orders, err := s.repos.PurchaseOrder.List(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repos.PurchaseOrder.Count(ctx, params)
	if err != nil {
		return nil, 0, err
	}
	return orders.ToAPIResponse(), count, nil
}

func (s *purchaseOrderService) ReceivePurchaseOrder(ctx context.Context, orderID int, input *model.ReceivePurchaseOrderRequest, performedBy *int) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if len(input.Items) == 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "at least one receive item is required")
	}
	receivedDate := time.Now().UTC()
	if input.ReceivedDate != nil {
		receivedDate = *input.ReceivedDate
	}
	if err := s.repos.PurchaseOrder.Receive(ctx, orderID, &receivedDate, input.Notes, input.Items, performedBy); err != nil {
		return err
	}
	return nil
}

func (s *purchaseOrderService) PutAwayPurchaseOrder(ctx context.Context, orderID int, input *model.PutAwayPurchaseOrderRequest, performedBy *int) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if len(input.Items) == 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "at least one put-away item is required")
	}
	return s.repos.PurchaseOrder.PutAway(ctx, orderID, input.Notes, input.Items, performedBy)
}
