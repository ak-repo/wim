package model

import (
	"database/sql"
	"time"
)

type LocationDTO struct {
	ID           int        `db:"id"`
	RefCode      string     `db:"ref_code"`
	WarehouseID  int        `db:"warehouse_id"`
	Zone         string     `db:"zone"`
	Aisle        *string    `db:"aisle"`
	Rack         *string    `db:"rack"`
	Bin          *string    `db:"bin"`
	LocationCode string     `db:"location_code"`
	LocationType string     `db:"location_type"`
	IsPickFace   bool       `db:"is_pick_face"`
	MaxWeight    *float64   `db:"max_weight"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type LocationResponse struct {
	ID           int       `json:"id"`
	RefCode      string    `json:"refCode"`
	WarehouseID  int       `json:"warehouseId"`
	Zone         string    `json:"zone"`
	Aisle        *string   `json:"aisle,omitempty"`
	Rack         *string   `json:"rack,omitempty"`
	Bin          *string   `json:"bin,omitempty"`
	LocationCode string    `json:"locationCode"`
	LocationType string    `json:"locationType"`
	IsPickFace   bool      `json:"isPickFace"`
	MaxWeight    *float64  `json:"maxWeight,omitempty"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LocationRequest struct {
	ID           int      `json:"id,omitempty"`
	RefCode      string   `json:"refCode,omitempty"`
	WarehouseID  int      `json:"warehouseId,omitempty"`
	Zone         *string  `json:"zone,omitempty"`
	Aisle        *string  `json:"aisle,omitempty"`
	Rack         *string  `json:"rack,omitempty"`
	Bin          *string  `json:"bin,omitempty"`
	LocationCode *string  `json:"locationCode,omitempty"`
	LocationType *string  `json:"locationType,omitempty"`
	IsPickFace   *bool    `json:"isPickFace,omitempty"`
	IsActive     *bool    `json:"isActive,omitempty"`
	MaxWeight    *float64 `json:"maxWeight,omitempty"`
}

type LocationParams struct {
	Active      *bool  `json:"active"`
	WarehouseID int    `json:"warehouseId,omitempty"`
	Zone        string `json:"zone,omitempty"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
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
		RefCode:      m.RefCode,
		WarehouseID:  m.WarehouseID,
		Zone:         m.Zone,
		Aisle:        m.Aisle,
		Rack:         m.Rack,
		Bin:          m.Bin,
		LocationCode: m.LocationCode,
		LocationType: m.LocationType,
		IsPickFace:   m.IsPickFace,
		MaxWeight:    m.MaxWeight,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (m *LocationDTO) ApplyNullScalars(isActive sql.NullBool, isPickFace sql.NullBool, createdAt, updatedAt sql.NullTime, deletedAt sql.NullTime) {
	if isActive.Valid {
		m.IsActive = isActive.Bool
	} else {
		m.IsActive = true
	}
	if isPickFace.Valid {
		m.IsPickFace = isPickFace.Bool
	} else {
		m.IsPickFace = false
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
	if deletedAt.Valid {
		m.DeletedAt = &deletedAt.Time
	} else {
		m.DeletedAt = nil
	}
}
