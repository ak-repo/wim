package model

import (
	"database/sql"
	"time"

	"github.com/ak-repo/wim/pkg/utils"
)

type UserDTO struct {
	ID           int        `db:"id"`
	RefCode      string     `db:"ref_code"`
	Username     string     `db:"username"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Role         string     `db:"role"`
	Contact      *string    `db:"contact"`
	IsActive     bool       `db:"is_active"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	RefCode   string    `json:"refCode"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Contact   *string   `json:"contact,omitempty"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserRequest struct {
	ID           int     `json:"id,omitempty"`
	RefCode      string  `json:"refCode,omitempty"`
	Username     *string `json:"username,omitempty"`
	Email        *string `json:"email,omitempty"`
	PasswordHash *string `json:"passwordHash,omitempty"`
	Role         *string `json:"role,omitempty"`
	Contact      *string `json:"contact,omitempty"`
	IsActive     *bool   `json:"isActive,omitempty"`
}

// Parameter struct for List API
type UserParams struct {
	Active *bool `json:"active"`
	Page   int   `json:"page"`
	Limit  int   `json:"limit"`
}

// API conversion
type UserDTOs []*UserDTO

func (m *UserDTOs) ToAPIResponse() []*UserResponse {
	var responses []*UserResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPIResponse())
	}
	return responses
}

func (m *UserDTO) ToAPIResponse() *UserResponse {
	return &UserResponse{
		ID:        m.ID,
		RefCode:   m.RefCode,
		Username:  m.Username,
		Email:     m.Email,
		Role:      m.Role,
		Contact:   m.Contact,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (m *UserDTO) ApplyNullScalars(isActive sql.NullBool, createdAt, updatedAt sql.NullTime, deletedAt sql.NullTime) {
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

// GetContact returns string value or empty string if nil
func (m *UserDTO) GetContact() string {
	return utils.NilOrString(m.Contact)
}
