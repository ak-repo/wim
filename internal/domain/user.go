package domain

import (
	"time"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(username, email, passwordHash string, role string) *User {
	if role == "" {
		role = constants.RoleStaff
	}

	now := time.Now()

	return &User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
