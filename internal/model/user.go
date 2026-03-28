package model

import (
	"database/sql"
	"time"

	"github.com/ak-repo/wim/pkg/utils"
	"github.com/google/uuid"
)

type UserDTO struct {
	ID           uuid.UUID      `db:"id"`
	Username     sql.NullString `db:"username"`
	Email        sql.NullString `db:"email"`
	PasswordHash sql.NullString `db:"password_hash"`
	Role         sql.NullString `db:"role"`
	Contact      sql.NullString `db:"contact,omitempty"`
	IsActive     bool           `db:"isActive"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	DeletedAt    sql.NullTime   `db:"deleted_at"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Contact   *string   `json:"contact,omitempty"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRequest struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	Contact      *string   `json:"contact,omitempty"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// API conversion
type UserDTOs []*UserDTO

func (m *UserDTOs) ToAPIRequest() []*UserResponse {
	var responses []*UserResponse
	for _, dto := range *m {
		responses = append(responses, dto.ToAPI())
	}
	return responses
}

func (m *UserDTO) ToAPI() *UserResponse {
	return &UserResponse{
		ID:        m.ID,
		Username:  m.Username.String,
		Email:     m.Email.String,
		Role:      m.Role.String,
		Contact:   utils.StringNil(m.Contact),
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

}
