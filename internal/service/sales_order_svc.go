package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/event"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/internal/errs"
)

type SalesOrderService interface {
	CreateSalesOrder(ctx context.Context, input *model.CreateSalesOrderRequest, createdBy *int) (*model.SalesOrderResponse, error)
	GetSalesOrderByID(ctx context.Context, orderID int) (*model.SalesOrderResponse, error)
	GetSalesOrderByRefCode(ctx context.Context, refCode string) (*model.SalesOrderResponse, error)
	ListSalesOrders(ctx context.Context, params *model.SalesOrderParams) ([]*model.SalesOrderResponse, int, error)
	UpdateSalesOrder(ctx context.Context, orderID int, input *model.CreateSalesOrderRequest) (*model.SalesOrderResponse, error)
	CancelSalesOrder(ctx context.Context, orderID int) error

	// Allocation
	AllocateSalesOrder(ctx context.Context, orderID int, performedBy *int) error
	DeallocateSalesOrder(ctx context.Context, orderID int) error

	// Shipping
	ShipSalesOrder(ctx context.Context, orderID int, input *model.ShipSalesOrderRequest, performedBy *int) error
}

type salesOrderService struct {
	repos  *repository.Repositories
	events event.EventPublisher
}

func NewSalesOrderService(repositories *repository.Repositories, eventPublisher event.EventPublisher) SalesOrderService {
	return &salesOrderService{
		repos:  repositories,
		events: eventPublisher,
	}
}

func (s *salesOrderService) CreateSalesOrder(ctx context.Context, input *model.CreateSalesOrderRequest, createdBy *int) (*model.SalesOrderResponse, error) {
	if input == nil {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.CustomerID <= 0 {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "customerId is required")
	}
	if input.WarehouseID <= 0 {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "warehouseId is required")
	}
	if len(input.Items) == 0 {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "at least one item is required")
	}

	// Verify warehouse exists
	if _, err := s.repos.Warehouse.GetByID(ctx, input.WarehouseID); err != nil {
		if errors.Is(err, repository.ErrWarehouseNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "warehouse not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load warehouse")
	}

	// Generate ref code
	refCode, err := s.repos.RefCode.GenerateSalesOrderRefCode(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeRefCodeFailed, "failed to generate reference code")
	}

	// Prepare items
	var items []*model.SalesOrderItemDTO
	for _, itemReq := range input.Items {
		if itemReq.ProductID <= 0 {
			return nil, apperrors.New(apperrors.CodeInvalidInput, "productId is required for all items")
		}
		if itemReq.QuantityOrdered <= 0 {
			return nil, apperrors.New(apperrors.CodeInvalidInput, "quantityOrdered must be greater than 0")
		}

		// Verify product exists
		if _, err := s.repos.Product.GetByID(ctx, itemReq.ProductID); err != nil {
			if errors.Is(err, repository.ErrProductNotFound) {
				return nil, apperrors.New(apperrors.CodeNotFound, fmt.Sprintf("product not found: %d", itemReq.ProductID))
			}
			return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load product")
		}

		item := &model.SalesOrderItemDTO{
			ProductID:        itemReq.ProductID,
			QuantityOrdered:  itemReq.QuantityOrdered,
			AllocationStatus: constants.StatusUnallocated,
		}
		if itemReq.UnitPrice != nil {
			item.UnitPrice = itemReq.UnitPrice
		}
		items = append(items, item)
	}

	order := &model.SalesOrderDTO{
		RefCode:          refCode,
		CustomerID:       input.CustomerID,
		WarehouseID:      input.WarehouseID,
		Status:           constants.StatusPending,
		AllocationStatus: constants.StatusUnallocated,
		OrderDate:        time.Now(),
	}

	if input.RequiredDate != nil {
		order.RequiredDate = input.RequiredDate
	}
	if strings.TrimSpace(input.ShippingMethod) != "" {
		order.ShippingMethod = &input.ShippingMethod
	}
	if strings.TrimSpace(input.ShippingAddress) != "" {
		order.ShippingAddress = &input.ShippingAddress
	}
	if strings.TrimSpace(input.BillingAddress) != "" {
		order.BillingAddress = &input.BillingAddress
	}
	if strings.TrimSpace(input.Notes) != "" {
		order.Notes = &input.Notes
	}
	order.CreatedBy = createdBy

	createdOrder, err := s.repos.SalesOrder.Create(ctx, order, items)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create sales order")
	}

	// Load items for response
	orderItems, err := s.repos.SalesOrder.GetItemsByOrderID(ctx, createdOrder.ID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load order items")
	}

	response := createdOrder.ToAPIResponse()
	response.Items = orderItems.ToAPIResponse()

	// Publish event
	if s.events != nil {
		eventData, _ := json.Marshal(map[string]interface{}{
			"orderId":     createdOrder.ID,
			"refCode":     createdOrder.RefCode,
			"customerId":  createdOrder.CustomerID,
			"warehouseId": createdOrder.WarehouseID,
			"status":      createdOrder.Status,
			"itemCount":   len(items),
		})
		s.events.Publish(ctx, event.Event{
			Type:      event.EventOrderCreated,
			Payload:   eventData,
			Timestamp: time.Now(),
		})
	}

	return response, nil
}

func (s *salesOrderService) GetSalesOrderByID(ctx context.Context, orderID int) (*model.SalesOrderResponse, error) {
	order, err := s.repos.SalesOrder.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrSalesOrderNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "sales order not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	items, err := s.repos.SalesOrder.GetItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load order items")
	}

	response := order.ToAPIResponse()
	response.Items = items.ToAPIResponse()
	return response, nil
}

func (s *salesOrderService) GetSalesOrderByRefCode(ctx context.Context, refCode string) (*model.SalesOrderResponse, error) {
	order, err := s.repos.SalesOrder.GetByRefCode(ctx, refCode)
	if err != nil {
		if errors.Is(err, repository.ErrSalesOrderNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "sales order not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	items, err := s.repos.SalesOrder.GetItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load order items")
	}

	response := order.ToAPIResponse()
	response.Items = items.ToAPIResponse()
	return response, nil
}

func (s *salesOrderService) ListSalesOrders(ctx context.Context, params *model.SalesOrderParams) ([]*model.SalesOrderResponse, int, error) {
	if params == nil {
		params = &model.SalesOrderParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	orders, err := s.repos.SalesOrder.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list sales orders")
	}

	count, err := s.repos.SalesOrder.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count sales orders")
	}

	return orders.ToAPIResponse(), count, nil
}

func (s *salesOrderService) UpdateSalesOrder(ctx context.Context, orderID int, input *model.CreateSalesOrderRequest) (*model.SalesOrderResponse, error) {
	if input == nil {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	order, err := s.repos.SalesOrder.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrSalesOrderNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "sales order not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	if order.Status == constants.StatusShipped {
		return nil, apperrors.New(apperrors.CodeInvalidOperation, "cannot update shipped order")
	}
	if order.Status == constants.StatusCancelled {
		return nil, apperrors.New(apperrors.CodeInvalidOperation, "cannot update cancelled order")
	}

	// Update basic fields
	if input.CustomerID > 0 {
		order.CustomerID = input.CustomerID
	}
	if input.WarehouseID > 0 {
		order.WarehouseID = input.WarehouseID
	}
	if input.RequiredDate != nil {
		order.RequiredDate = input.RequiredDate
	}
	if strings.TrimSpace(input.ShippingMethod) != "" {
		order.ShippingMethod = &input.ShippingMethod
	}
	if strings.TrimSpace(input.ShippingAddress) != "" {
		order.ShippingAddress = &input.ShippingAddress
	}
	if strings.TrimSpace(input.BillingAddress) != "" {
		order.BillingAddress = &input.BillingAddress
	}
	if strings.TrimSpace(input.Notes) != "" {
		order.Notes = &input.Notes
	}

	if err := s.repos.SalesOrder.Update(ctx, order); err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update sales order")
	}

	return s.GetSalesOrderByID(ctx, orderID)
}

func (s *salesOrderService) CancelSalesOrder(ctx context.Context, orderID int) error {
	order, err := s.repos.SalesOrder.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrSalesOrderNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "sales order not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	if order.Status == constants.StatusShipped {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot cancel shipped order")
	}
	if order.Status == constants.StatusCancelled {
		return apperrors.New(apperrors.CodeInvalidOperation, "order already cancelled")
	}

	// Deallocate any reserved stock first
	if order.AllocationStatus == constants.StatusPartiallyAllocated || order.AllocationStatus == constants.StatusFullyAllocated {
		if err := s.repos.SalesOrder.DeallocateStock(ctx, orderID); err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to deallocate stock")
		}
	}

	if err := s.repos.SalesOrder.UpdateStatus(ctx, orderID, constants.StatusCancelled); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to cancel sales order")
	}

	return nil
}

func (s *salesOrderService) AllocateSalesOrder(ctx context.Context, orderID int, performedBy *int) error {
	order, err := s.repos.SalesOrder.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrSalesOrderNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "sales order not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	if order.Status == constants.StatusShipped {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot allocate shipped order")
	}
	if order.Status == constants.StatusCancelled {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot allocate cancelled order")
	}
	if order.AllocationStatus == constants.StatusFullyAllocated {
		return apperrors.New(apperrors.CodeInvalidOperation, "order already fully allocated")
	}

	if err := s.repos.SalesOrder.AllocateStock(ctx, orderID, performedBy); err != nil {
		if errors.Is(err, apperrors.ErrInsufficientStock) {
			return apperrors.New(apperrors.CodeInsufficientStock, "insufficient stock for allocation")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to allocate stock")
	}

	// Publish event
	if s.events != nil {
		eventData, _ := json.Marshal(map[string]interface{}{
			"orderId": orderID,
			"refCode": order.RefCode,
		})
		s.events.Publish(ctx, event.Event{
			Type:      event.EventOrderAllocated,
			Payload:   eventData,
			Timestamp: time.Now(),
		})
	}

	return nil
}

func (s *salesOrderService) DeallocateSalesOrder(ctx context.Context, orderID int) error {
	order, err := s.repos.SalesOrder.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrSalesOrderNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "sales order not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	if order.Status == constants.StatusShipped {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot deallocate shipped order")
	}
	if order.AllocationStatus == constants.StatusUnallocated {
		return apperrors.New(apperrors.CodeInvalidOperation, "order has no allocated stock")
	}

	if err := s.repos.SalesOrder.DeallocateStock(ctx, orderID); err != nil {
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to deallocate stock")
	}

	return nil
}

func (s *salesOrderService) ShipSalesOrder(ctx context.Context, orderID int, input *model.ShipSalesOrderRequest, performedBy *int) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	order, err := s.repos.SalesOrder.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, repository.ErrSalesOrderNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "sales order not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	if order.Status == constants.StatusShipped {
		return apperrors.New(apperrors.CodeInvalidOperation, "order already shipped")
	}
	if order.Status == constants.StatusCancelled {
		return apperrors.New(apperrors.CodeInvalidOperation, "cannot ship cancelled order")
	}
	if order.AllocationStatus != constants.StatusFullyAllocated {
		return apperrors.New(apperrors.CodeInvalidOperation, "order must be fully allocated before shipping")
	}

	var shipDate *time.Time
	if input.ShippedDate != nil {
		shipDate = input.ShippedDate
	}

	if err := s.repos.SalesOrder.ShipOrder(ctx, orderID, shipDate, performedBy); err != nil {
		if errors.Is(err, apperrors.ErrInsufficientStock) {
			return apperrors.New(apperrors.CodeInsufficientStock, "insufficient stock for shipping")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to ship order")
	}

	// Publish event
	if s.events != nil {
		eventData, _ := json.Marshal(map[string]interface{}{
			"orderId":     orderID,
			"refCode":     order.RefCode,
			"shippedDate": shipDate,
		})
		s.events.Publish(ctx, event.Event{
			Type:      event.EventOrderShipped,
			Payload:   eventData,
			Timestamp: time.Now(),
		})
	}

	return nil
}
