package model

import (
	"database/sql"
	"time"
)

type ProductDTO struct {
	ID            int        `db:"id"`
	RefCode       string     `db:"ref_code"`
	SKU           string     `db:"sku"`
	Name          string     `db:"name"`
	Description   *string    `db:"description"`
	Category      *string    `db:"category"`
	UnitOfMeasure string     `db:"unit_of_measure"` // REQUIRED → value
	Weight        *float64   `db:"weight"`
	Length        *float64   `db:"length"`
	Width         *float64   `db:"width"`
	Height        *float64   `db:"height"`
	Barcode       *string    `db:"barcode"`
	IsActive      bool       `db:"is_active"`  // REQUIRED → value
	CreatedAt     time.Time  `db:"created_at"` // REQUIRED → value
	UpdatedAt     time.Time  `db:"updated_at"` // REQUIRED → value
	DeletedAt     *time.Time `db:"deleted_at"` // nullable
}

type ProductResponse struct {
	ID            int       `json:"id"`
	RefCode       string    `json:"refCode"`
	SKU           string    `json:"sku"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	Category      *string   `json:"category,omitempty"`
	UnitOfMeasure string    `json:"unitOfMeasure"`
	Weight        *float64  `json:"weight,omitempty"`
	Length        *float64  `json:"length,omitempty"`
	Width         *float64  `json:"width,omitempty"`
	Height        *float64  `json:"height,omitempty"`
	Barcode       *string   `json:"barcode,omitempty"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type ProductRequest struct {
	ID            int      `json:"id,omitempty"`
	RefCode       string   `json:"refCode,omitempty"`
	SKU           *string  `json:"sku,omitempty"`
	Name          *string  `json:"name,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Category      *string  `json:"category,omitempty"`
	UnitOfMeasure *string  `json:"unitOfMeasure,omitempty"`
	Weight        *float64 `json:"weight,omitempty"`
	Length        *float64 `json:"length,omitempty"`
	Width         *float64 `json:"width,omitempty"`
	Height        *float64 `json:"height,omitempty"`
	Barcode       *string  `json:"barcode,omitempty"`
	IsActive      *bool    `json:"isActive,omitempty"`
}

type ProductDTOs []*ProductDTO

// API conversion
func (m *ProductDTOs) ToAPIResponse() []*ProductResponse {
	var responses []*ProductResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}
func (m *ProductDTO) ToAPIResponse() *ProductResponse {
	return &ProductResponse{
		ID:            m.ID,
		RefCode:       m.RefCode,
		SKU:           m.SKU,
		Name:          m.Name,
		Description:   m.Description,
		Category:      m.Category,
		UnitOfMeasure: m.UnitOfMeasure,
		Weight:        m.Weight,
		Length:        m.Length,
		Width:         m.Width,
		Height:        m.Height,
		Barcode:       m.Barcode,
		IsActive:      m.IsActive,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

type ProductParams struct {
	Active   *bool  `json:"active"`
	Category string `json:"category,omitempty"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

func (m *ProductDTO) ApplyProductNullScalars(isActive sql.NullBool, createdAt, updatedAt sql.NullTime) {
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
}
