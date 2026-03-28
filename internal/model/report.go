package model

import "time"

type InventoryReportRequest struct {
	WarehouseID string     `json:"warehouseId,omitempty"`
	ProductID   string     `json:"productId,omitempty"`
	FromDate    *time.Time `json:"fromDate,omitempty"`
	ToDate      *time.Time `json:"toDate,omitempty"`
}

type MovementReportRequest struct {
	MovementType string     `json:"movementType,omitempty"`
	WarehouseID  string     `json:"warehouseId,omitempty"`
	FromDate     *time.Time `json:"fromDate,omitempty"`
	ToDate       *time.Time `json:"toDate,omitempty"`
}

type ExpiryReportRequest struct {
	WarehouseID string `json:"warehouseId,omitempty"`
	DaysAhead   int    `json:"daysAhead,omitempty"`
}
