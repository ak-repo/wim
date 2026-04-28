package model

import (
	"database/sql"
	"time"
)

type PickingTaskStatus string

const (
	PickingStatusPending    PickingTaskStatus = "PENDING"
	PickingStatusInProgress PickingTaskStatus = "IN_PROGRESS"
	PickingStatusCompleted  PickingTaskStatus = "COMPLETED"
	PickingStatusCancelled  PickingTaskStatus = "CANCELLED"
)

type PickingPriority string

const (
	PickingPriorityLow    PickingPriority = "LOW"
	PickingPriorityMedium PickingPriority = "MEDIUM"
	PickingPriorityHigh   PickingPriority = "HIGH"
	PickingPriorityUrgent PickingPriority = "URGENT"
)

type PickingTaskDTO struct {
	ID           int
	RefCode      string        `db:"ref_code"`
	SalesOrderID int           `db:"sales_order_id"`
	WarehouseID  int           `db:"warehouse_id"`
	Status       string        `db:"status"`
	Priority     string        `db:"priority"`
	AssignedTo   sql.NullInt64 `db:"assigned_to"`
	StartedAt    sql.NullTime  `db:"started_at"`
	CompletedAt  sql.NullTime  `db:"completed_at"`
	Notes        sql.NullString `db:"notes"`
	CreatedBy    sql.NullInt64 `db:"created_by"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
}

type PickingTaskItemDTO struct {
	ID               int            `db:"id"`
	PickingTaskID    int            `db:"picking_task_id"`
	SalesOrderItemID int            `db:"sales_order_item_id"`
	ProductID        int            `db:"product_id"`
	LocationID       sql.NullInt64  `db:"location_id"`
	BatchID          sql.NullInt64   `db:"batch_id"`
	QuantityRequired int            `db:"quantity_required"`
	QuantityPicked   int            `db:"quantity_picked"`
	PickedAt         sql.NullTime   `db:"picked_at"`
	Status           string         `db:"status"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
}

type CreatePickingTaskRequest struct {
	SalesOrderID int    `json:"salesOrderId"`
	Priority     string `json:"priority,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type AssignPickingTaskRequest struct {
	AssignedTo int    `json:"assignedTo"`
	Notes      string `json:"notes,omitempty"`
}

type PickItemRequest struct {
	PickingTaskItemID int  `json:"pickingTaskItemId"`
	Quantity          int  `json:"quantity"`
	LocationID        int  `json:"locationId"`
	BatchID           *int `json:"batchId,omitempty"`
}

type CompletePickingRequest struct {
	Items []PickItemRequest `json:"items"`
	Notes string            `json:"notes,omitempty"`
}

type PickingTaskResponse struct {
	ID           int                       `json:"id"`
	RefCode      string                    `json:"refCode"`
	SalesOrderID int                       `json:"salesOrderId"`
	WarehouseID  int                       `json:"warehouseId"`
	Status       string                    `json:"status"`
	Priority     string                    `json:"priority"`
	AssignedTo   *int                      `json:"assignedTo,omitempty"`
	AssignedUser *string                   `json:"assignedUser,omitempty"`
	StartedAt    *string                   `json:"startedAt,omitempty"`
	CompletedAt  *string                   `json:"completedAt,omitempty"`
	Notes        *string                   `json:"notes,omitempty"`
	CreatedBy    *int                      `json:"createdBy,omitempty"`
	CreatedAt    string                    `json:"createdAt"`
	UpdatedAt    string                    `json:"updatedAt"`
	Items        []*PickingTaskItemResponse `json:"items,omitempty"`
}

type PickingTaskItemResponse struct {
	ID               int     `json:"id"`
	PickingTaskID    int     `json:"pickingTaskId"`
	SalesOrderItemID int     `json:"salesOrderItemId"`
	ProductID        int     `json:"productId"`
	ProductName      string  `json:"productName,omitempty"`
	LocationID       *int    `json:"locationId,omitempty"`
	LocationCode     *string `json:"locationCode,omitempty"`
	BatchID          *int    `json:"batchId,omitempty"`
	QuantityRequired int     `json:"quantityRequired"`
	QuantityPicked   int     `json:"quantityPicked"`
	PickedAt         *string `json:"pickedAt,omitempty"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type PickingTaskParams struct {
	WarehouseID *int `json:"warehouseId,omitempty"`
	Status      *string `json:"status,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	AssignedTo  *int    `json:"assignedTo,omitempty"`
	Page        int     `json:"page,omitempty"`
	Limit       int     `json:"limit,omitempty"`
}

type PickingTaskDTOs []*PickingTaskDTO

func (m *PickingTaskDTO) ToAPIResponse() *PickingTaskResponse {
	resp := &PickingTaskResponse{
		ID:           m.ID,
		RefCode:      m.RefCode,
		SalesOrderID: m.SalesOrderID,
		WarehouseID:  m.WarehouseID,
		Status:       m.Status,
		Priority:     m.Priority,
		CreatedAt:    m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    m.UpdatedAt.Format(time.RFC3339),
	}

	if m.AssignedTo.Valid {
		assignedTo := int(m.AssignedTo.Int64)
		resp.AssignedTo = &assignedTo
	}
	if m.StartedAt.Valid {
		startedAt := m.StartedAt.Time.Format(time.RFC3339)
		resp.StartedAt = &startedAt
	}
	if m.CompletedAt.Valid {
		completedAt := m.CompletedAt.Time.Format(time.RFC3339)
		resp.CompletedAt = &completedAt
	}
	if m.Notes.Valid {
		notes := m.Notes.String
		resp.Notes = &notes
	}
	if m.CreatedBy.Valid {
		createdBy := int(m.CreatedBy.Int64)
		resp.CreatedBy = &createdBy
	}

	return resp
}

func (m *PickingTaskItemDTO) ToAPIResponse() *PickingTaskItemResponse {
	resp := &PickingTaskItemResponse{
		ID:               m.ID,
		PickingTaskID:    m.PickingTaskID,
		SalesOrderItemID: m.SalesOrderItemID,
		ProductID:        m.ProductID,
		QuantityRequired: m.QuantityRequired,
		QuantityPicked:   m.QuantityPicked,
		Status:           m.Status,
		CreatedAt:        m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        m.UpdatedAt.Format(time.RFC3339),
	}

	if m.LocationID.Valid {
		locationID := int(m.LocationID.Int64)
		resp.LocationID = &locationID
	}
	if m.BatchID.Valid {
		batchID := int(m.BatchID.Int64)
		resp.BatchID = &batchID
	}
	if m.PickedAt.Valid {
		pickedAt := m.PickedAt.Time.Format(time.RFC3339)
		resp.PickedAt = &pickedAt
	}

	return resp
}

func (m *PickingTaskDTO) ApplyNullScalars(startedAt, completedAt, createdAt, updatedAt sql.NullTime) {
	if startedAt.Valid {
		m.StartedAt = startedAt
	}
	if completedAt.Valid {
		m.CompletedAt = completedAt
	}
	if createdAt.Valid {
		m.CreatedAt = createdAt.Time
	} else {
		m.CreatedAt = time.Time{}
	}
	if updatedAt.Valid {
		m.UpdatedAt = updatedAt.Time
	} else {
		m.UpdatedAt = time.Time{}
	}
}

func (m *PickingTaskItemDTO) ApplyNullScalars(createdAt, updatedAt sql.NullTime) {
	if createdAt.Valid {
		m.CreatedAt = createdAt.Time
	} else {
		m.CreatedAt = time.Time{}
	}
	if updatedAt.Valid {
		m.UpdatedAt = updatedAt.Time
	} else {
		m.UpdatedAt = time.Time{}
	}
}