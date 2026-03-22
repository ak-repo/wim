package repository

import (
	"context"
	"errors"

	"github.com/ak-repo/wim/internal/db"
	"github.com/ak-repo/wim/internal/domain"
	"github.com/google/uuid"
)

var ErrUserRepositoryNotImplemented = errors.New("user repository not implemented")

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *db.DB
}

func NewUserRepository(database *db.DB) UserRepository {
	return &userRepository{db: database}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return ErrUserRepositoryNotImplemented
}

func (r *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return nil, ErrUserRepositoryNotImplemented
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, ErrUserRepositoryNotImplemented
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, ErrUserRepositoryNotImplemented
}
