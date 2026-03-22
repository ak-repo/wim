package domain

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID           uuid.UUID
	Code         string
	Name         string
	AddressLine1 string
	AddressLine2 string
	City         string
	State        string
	PostalCode   string
	Country      string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewWarehouse(code, name, country string) *Warehouse {
	return &Warehouse{
		ID:        uuid.New(),
		Code:      code,
		Name:      name,
		Country:   country,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
