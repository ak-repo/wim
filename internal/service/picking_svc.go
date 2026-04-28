package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/constants"
	apperrors "github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
)

type PickingService interface {
	CreatePickingTask(ctx context.Context, input *model.CreatePickingTaskRequest, createdBy *int) (*model.PickingTaskResponse, error)
	GetPickingTaskByID(ctx context.Context, taskID int) (*model.PickingTaskResponse, error)
	GetPickingTaskByRefCode(ctx context.Context, refCode string) (*model.PickingTaskResponse, error)
	ListPickingTasks(ctx context.Context, params *model.PickingTaskParams) ([]*model.PickingTaskResponse, int, error)
	AssignPickingTask(ctx context.Context, taskID int, input *model.AssignPickingTaskRequest) error
	StartPickingTask(ctx context.Context, taskID int) error
	PickItem(ctx context.Context, taskID int, input *model.PickItemRequest, performedBy *int) error
	CompletePickingTask(ctx context.Context, taskID int, notes string) error
	CancelPickingTask(ctx context.Context, taskID int, notes string) error
}

type pickingService struct {
	repos *repository.Repositories
}

func NewPickingService(repositories *repository.Repositories) PickingService {
	return &pickingService{repos: repositories}
}

func (s *pickingService) CreatePickingTask(ctx context.Context, input *model.CreatePickingTaskRequest, createdBy *int) (*model.PickingTaskResponse, error) {
	if input == nil {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if input.SalesOrderID <= 0 {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "salesOrderId is required")
	}

	salesOrder, err := s.repos.SalesOrder.GetByID(ctx, input.SalesOrderID)
	if err != nil {
		if errors.Is(err, repository.ErrSalesOrderNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "sales order not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order")
	}

	if salesOrder.Status != constants.StatusProcessing {
		return nil, apperrors.New(apperrors.CodeInvalidOperation, "sales order must be in processing status")
	}

	if salesOrder.AllocationStatus != constants.StatusFullyAllocated {
		return nil, apperrors.New(apperrors.CodeInvalidOperation, "sales order must be fully allocated")
	}

	salesOrderItems, err := s.repos.SalesOrder.GetItemsByOrderID(ctx, salesOrder.ID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load sales order items")
	}

	task := &model.PickingTaskDTO{
		SalesOrderID: salesOrder.ID,
		WarehouseID:  salesOrder.WarehouseID,
		Status:       string(model.PickingStatusPending),
		Priority:     input.Priority,
	}
	if task.Priority == "" {
		task.Priority = string(model.PickingPriorityMedium)
	}
	if createdBy != nil {
		task.CreatedBy = sql.NullInt64{Int64: int64(*createdBy), Valid: true}
	}

	var items []model.PickingTaskItemDTO
	for _, soItem := range salesOrderItems {
		if soItem.QuantityReserved > 0 {
			item := model.PickingTaskItemDTO{
				SalesOrderItemID: soItem.ID,
				ProductID:        soItem.ProductID,
				QuantityRequired: soItem.QuantityReserved,
				Status:           string(model.PickingStatusPending),
			}
			items = append(items, item)
		}
	}

	if len(items) == 0 {
		return nil, apperrors.New(apperrors.CodeInvalidOperation, "no reserved items to pick")
	}

	createdTask, err := s.repos.Picking.CreateTask(ctx, task, items)
	if err != nil {
		return nil, err
	}

	return s.GetPickingTaskByID(ctx, createdTask.ID)
}

func (s *pickingService) GetPickingTaskByID(ctx context.Context, taskID int) (*model.PickingTaskResponse, error) {
	task, err := s.repos.Picking.GetTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrPickingTaskNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "picking task not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load picking task")
	}

	items, err := s.repos.Picking.GetTaskItems(ctx, taskID)
	if err != nil {
		return nil, err
	}

	response := task.ToAPIResponse()
	var itemResponses []*model.PickingTaskItemResponse
	for _, item := range items {
		itemResp := item.ToAPIResponse()

		product, err := s.repos.Product.GetByID(ctx, item.ProductID)
		if err == nil && product.Name != "" {
			itemResp.ProductName = product.Name
		}

		if item.LocationID.Valid {
			location, err := s.repos.Location.GetByID(ctx, int(item.LocationID.Int64))
			if err == nil && location.LocationCode != "" {
				itemResp.LocationCode = &location.LocationCode
			}
		}

		itemResponses = append(itemResponses, itemResp)
	}
	response.Items = itemResponses

	if task.AssignedTo.Valid {
		user, err := s.repos.User.GetByID(ctx, int(task.AssignedTo.Int64))
		if err == nil && user.Username != "" {
			response.AssignedUser = &user.Username
		}
	}

	return response, nil
}

func (s *pickingService) GetPickingTaskByRefCode(ctx context.Context, refCode string) (*model.PickingTaskResponse, error) {
	task, err := s.repos.Picking.GetTaskByRefCode(ctx, refCode)
	if err != nil {
		if errors.Is(err, repository.ErrPickingTaskNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "picking task not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load picking task")
	}

	return s.GetPickingTaskByID(ctx, task.ID)
}

func (s *pickingService) ListPickingTasks(ctx context.Context, params *model.PickingTaskParams) ([]*model.PickingTaskResponse, int, error) {
	if params == nil {
		params = &model.PickingTaskParams{}
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	tasks, err := s.repos.Picking.ListTasks(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.repos.Picking.CountTasks(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	var responses []*model.PickingTaskResponse
	for _, task := range tasks {
		response := task.ToAPIResponse()
		responses = append(responses, response)
	}

	return responses, count, nil
}

func (s *pickingService) AssignPickingTask(ctx context.Context, taskID int, input *model.AssignPickingTaskRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if input.AssignedTo <= 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "assignedTo is required")
	}

	_, err := s.repos.User.GetByID(ctx, input.AssignedTo)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "user not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load user")
	}

	return s.repos.Picking.AssignTask(ctx, taskID, input.AssignedTo, input.Notes)
}

func (s *pickingService) StartPickingTask(ctx context.Context, taskID int) error {
	return s.repos.Picking.StartTask(ctx, taskID)
}

func (s *pickingService) PickItem(ctx context.Context, taskID int, input *model.PickItemRequest, performedBy *int) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if input.Quantity <= 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "quantity must be greater than 0")
	}
	if input.LocationID <= 0 {
		return apperrors.New(apperrors.CodeInvalidInput, "locationId is required")
	}

	var performedByVal int
	if performedBy != nil {
		performedByVal = *performedBy
	}

	return s.repos.Picking.PickItem(ctx, taskID, input.PickingTaskItemID, input.Quantity, input.LocationID, input.BatchID, performedByVal)
}

func (s *pickingService) CompletePickingTask(ctx context.Context, taskID int, notes string) error {
	if strings.TrimSpace(notes) != "" {
		notes = strings.TrimSpace(notes)
	}
	return s.repos.Picking.CompleteTask(ctx, taskID, notes)
}

func (s *pickingService) CancelPickingTask(ctx context.Context, taskID int, notes string) error {
	if strings.TrimSpace(notes) != "" {
		notes = strings.TrimSpace(notes)
	}
	return s.repos.Picking.CancelTask(ctx, taskID, notes)
}