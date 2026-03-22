package domain

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	ID           uuid.UUID
	WarehouseID  uuid.UUID
	Zone         string
	Aisle        string
	Rack         string
	Bin          string
	LocationCode string
	LocationType string
	IsPickFace   bool
	MaxWeight    *float64
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewLocation(warehouseID uuid.UUID, zone, aisle, rack, bin, locationCode, locationType string) *Location {
	return &Location{
		ID:           uuid.New(),
		WarehouseID:  warehouseID,
		Zone:         zone,
		Aisle:        aisle,
		Rack:         rack,
		Bin:          bin,
		LocationCode: locationCode,
		LocationType: locationType,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
