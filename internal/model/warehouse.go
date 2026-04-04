package model

import (
	"database/sql"
	"time"

	"github.com/ak-repo/wim/pkg/utils"
)

type WarehouseDTO struct {
	ID           int        `db:"id"`
	RefCode      string     `db:"ref_code"`
	Code         string     `db:"code"`
	Name         string     `db:"name"`
	AddressLine1 *string    `db:"address_line1"`
	AddressLine2 *string    `db:"address_line2"`
	City         *string    `db:"city"`
	State        *string    `db:"state"`
	PostalCode   *string    `db:"postal_code"`
	Country      string     `db:"country"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type WarehouseRequest struct {
	ID           int     `json:"id,omitempty"`
	RefCode      string  `json:"refCode,omitempty"`
	Code         *string `json:"code,omitempty"`
	Name         *string `json:"name,omitempty"`
	AddressLine1 *string `json:"addressLine1,omitempty"`
	AddressLine2 *string `json:"addressLine2,omitempty"`
	City         *string `json:"city,omitempty"`
	State        *string `json:"state,omitempty"`
	PostalCode   *string `json:"postalCode,omitempty"`
	Country      *string `json:"country,omitempty"`
	IsActive     *bool   `json:"isActive,omitempty"`
}

type WarehouseResponse struct {
	ID           int       `json:"id"`
	RefCode      string    `json:"refCode"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	AddressLine1 *string   `json:"addressLine1,omitempty"`
	AddressLine2 *string   `json:"addressLine2,omitempty"`
	City         *string   `json:"city,omitempty"`
	State        *string   `json:"state,omitempty"`
	PostalCode   *string   `json:"postalCode,omitempty"`
	Country      string    `json:"country"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type WarehouseParams struct {
	Active *bool `json:"active"`
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
}

type WarehouseDTOs []*WarehouseDTO

func (m *WarehouseDTOs) ToAPIResponse() []*WarehouseResponse {
	var responses []*WarehouseResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *WarehouseDTO) ToAPIResponse() *WarehouseResponse {
	return &WarehouseResponse{
		ID:           m.ID,
		RefCode:      m.RefCode,
		Code:         m.Code,
		Name:         m.Name,
		AddressLine1: m.AddressLine1,
		AddressLine2: m.AddressLine2,
		City:         m.City,
		State:        m.State,
		PostalCode:   m.PostalCode,
		Country:      m.Country,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (m *WarehouseDTO) ApplyNullScalars(isActive sql.NullBool, createdAt, updatedAt sql.NullTime, deletedAt sql.NullTime) {
	if isActive.Valid {
		m.IsActive = isActive.Bool
	} else {
		m.IsActive = true
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

// Get helpers for nullable fields
func (m *WarehouseDTO) GetAddressLine1() string {
	return utils.NilOrString(m.AddressLine1)
}

func (m *WarehouseDTO) GetAddressLine2() string {
	return utils.NilOrString(m.AddressLine2)
}

func (m *WarehouseDTO) GetCity() string {
	return utils.NilOrString(m.City)
}

func (m *WarehouseDTO) GetState() string {
	return utils.NilOrString(m.State)
}

func (m *WarehouseDTO) GetPostalCode() string {
	return utils.NilOrString(m.PostalCode)
}
