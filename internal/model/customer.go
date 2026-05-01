package model

import (
	"database/sql"
	"time"
)

type CustomerDTO struct {
	ID        int        `db:"id"`
	RefCode   string     `db:"ref_code"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Contact   *string    `db:"contact"`
	Address   *string    `db:"address"`
	IsActive  bool       `db:"is_active"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type CustomerResponse struct {
	ID        int       `json:"id"`
	RefCode   string    `json:"refCode"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Contact   *string   `json:"contact,omitempty"`
	Address   *string   `json:"address,omitempty"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CustomerRequest struct {
	ID       int     `json:"id,omitempty"`
	RefCode  string  `json:"refCode,omitempty"`
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Contact  *string `json:"contact,omitempty"`
	Address  *string `json:"address,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}

type CustomerParams struct {
	Active *bool `json:"active"`
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
}

type CustomerDTOs []*CustomerDTO

func (m *CustomerDTOs) ToAPIResponse() []*CustomerResponse {
	var responses []*CustomerResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *CustomerDTO) ToAPIResponse() *CustomerResponse {
	return &CustomerResponse{
		ID:        m.ID,
		RefCode:   m.RefCode,
		Name:      m.Name,
		Email:     m.Email,
		Contact:   m.Contact,
		Address:   m.Address,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *CustomerDTO) ApplyNullScalars(isActive sql.NullBool, createdAt, updatedAt, deletedAt sql.NullTime) {
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
