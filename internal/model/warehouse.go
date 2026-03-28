package model

import (
	"time"

	"github.com/google/uuid"
)

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
