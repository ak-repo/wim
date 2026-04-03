package model

import (
	"time"

	"github.com/google/uuid"
)

type LocationDTO struct {
	ID           uuid.UUID `db:"id"`
	WarehouseID  uuid.UUID `db:"warehouse_id"`
	Zone         string    `db:"zone"`
	Aisle        *string   `db:"aisle"`
	Rack         *string   `db:"rack"`
	Bin          *string   `db:"bin"`
	LocationCode string    `db:"location_code"`
	LocationType string    `db:"location_type"`
	IsPickFace   bool      `db:"is_pick_face"`
	MaxWeight    *float64  `db:"max_weight"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

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

type LocationParams struct {
	Active      *bool     `json:"active"`
	WarehouseID uuid.UUID `json:"warehouseId,omitempty"`
	Zone        string    `json:"zone,omitempty"`
	Page        int       `json:"page"`
	Limit       int       `json:"limit"`
}

type LocationDTOs []*LocationDTO

func (m *LocationDTOs) ToAPIResponse() []*LocationResponse {
	var responses []*LocationResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *LocationDTO) ToAPIResponse() *LocationResponse {
	return &LocationResponse{
		ID:           m.ID,
		WarehouseID:  m.WarehouseID,
		Zone:         m.Zone,
		Aisle:        *m.Aisle,
		Rack:         *m.Rack,
		Bin:          *m.Bin,
		LocationCode: m.LocationCode,
		LocationType: m.LocationType,
		IsPickFace:   m.IsPickFace,
		MaxWeight:    m.MaxWeight,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
