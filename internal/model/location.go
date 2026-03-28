package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateLocationRequest struct {
	WarehouseID  uuid.UUID `json:"warehouseId"`
	Zone         string    `json:"zone"`
	Aisle        string    `json:"aisle,omitempty"`
	Rack         string    `json:"rack,omitempty"`
	Bin          string    `json:"bin,omitempty"`
	LocationCode string    `json:"locationCode"`
	LocationType string    `json:"locationType"`
	IsPickFace   bool      `json:"isPickFace"`
	MaxWeight    *float64  `json:"maxWeight,omitempty"`
}

type UpdateLocationRequest struct {
	Zone         string   `json:"zone,omitempty"`
	Aisle        string   `json:"aisle,omitempty"`
	Rack         string   `json:"rack,omitempty"`
	Bin          string   `json:"bin,omitempty"`
	LocationCode string   `json:"locationCode,omitempty"`
	LocationType string   `json:"locationType,omitempty"`
	IsPickFace   *bool    `json:"isPickFace,omitempty"`
	MaxWeight    *float64 `json:"maxWeight,omitempty"`
	IsActive     *bool    `json:"isActive,omitempty"`
}

type LocationResponse struct {
	ID           uuid.UUID `json:"id"`
	WarehouseID  uuid.UUID `json:"warehouseId"`
	Zone         string    `json:"zone"`
	Aisle        string    `json:"aisle,omitempty"`
	Rack         string    `json:"rack,omitempty"`
	Bin          string    `json:"bin,omitempty"`
	LocationCode string    `json:"locationCode"`
	LocationType string    `json:"locationType"`
	IsPickFace   bool      `json:"isPickFace"`
	MaxWeight    *float64  `json:"maxWeight,omitempty"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
