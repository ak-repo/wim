package model

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseDTO struct {
	ID           uuid.UUID `db:"id"`
	Code         string    `db:"code"`
	Name         string    `db:"name"`
	AddressLine1 *string   `db:"address_line1"`
	AddressLine2 *string   `db:"address_line2"`
	City         *string   `db:"city"`
	State        *string   `db:"state"`
	PostalCode   *string   `db:"postal_code"`
	Country      string    `db:"country"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type CreateWarehouseRequest struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country"`
}

type UpdateWarehouseRequest struct {
	Name         string `json:"name,omitempty"`
	AddressLine1 string `json:"addressLine1,omitempty"`
	AddressLine2 string `json:"addressLine2,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PostalCode   string `json:"postalCode,omitempty"`
	Country      string `json:"country,omitempty"`
	IsActive     *bool  `json:"isActive,omitempty"`
}

type WarehouseResponse struct {
	ID           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	AddressLine1 string    `json:"addressLine1,omitempty"`
	AddressLine2 string    `json:"addressLine2,omitempty"`
	City         string    `json:"city,omitempty"`
	State        string    `json:"state,omitempty"`
	PostalCode   string    `json:"postalCode,omitempty"`
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
	var addr1, addr2, city, state, postal string
	if m.AddressLine1 != nil {
		addr1 = *m.AddressLine1
	}
	if m.AddressLine2 != nil {
		addr2 = *m.AddressLine2
	}
	if m.City != nil {
		city = *m.City
	}
	if m.State != nil {
		state = *m.State
	}
	if m.PostalCode != nil {
		postal = *m.PostalCode
	}
	return &WarehouseResponse{
		ID:           m.ID,
		Code:         m.Code,
		Name:         m.Name,
		AddressLine1: addr1,
		AddressLine2: addr2,
		City:         city,
		State:        state,
		PostalCode:   postal,
		Country:      m.Country,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
