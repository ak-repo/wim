package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/domain"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/google/uuid"
)

var ErrAuthWorkflowNotImplemented = errors.New("auth workflow not implemented")

type RegisterInput struct {
	Username string
	Email    string
	Password string
	Role     string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	AccessToken string       `json:"accessToken"`
	User        *domain.User `json:"user"`
}

type UserService interface {
	Register(ctx context.Context, input RegisterInput) (*domain.User, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type userService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) UserService {
	return &userService{
		users: users,
	}
}

func (s *userService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {

	return nil, ErrAuthWorkflowNotImplemented
}

func (s *userService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		return nil, apperrors.ErrInvalidInput
	}

	return nil, ErrAuthWorkflowNotImplemented
}

func (s *userService) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	if userID == uuid.Nil {
		return nil, apperrors.ErrInvalidInput
	}

	if s.users == nil {
		return nil, ErrAuthWorkflowNotImplemented
	}

	return nil, ErrAuthWorkflowNotImplemented
}
