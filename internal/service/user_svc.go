package service

import (
	"context"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/pkg/auth"
	"github.com/ak-repo/wim/pkg/errors"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, input *model.UserRequest) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	CreateUser(ctx context.Context, input *model.UserRequest) error

	ListUsers(ctx context.Context, params *model.UserParams) ([]*model.UserResponse, int, error)
}

type userService struct {
	repos     *repository.Repositories
	passwords auth.PasswordHasher
}

func NewUserService(repositories *repository.Repositories) UserService {
	return &userService{
		repos: repositories,
	}
}

func (s *userService) CreateUser(ctx context.Context, input *model.UserRequest) error {

	if strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.PasswordHash) == "" {
		return apperrors.ErrInvalidInput
	}
	if len(strings.TrimSpace(input.PasswordHash)) < 8 || !strings.Contains(input.Email, "@") {
		return apperrors.ErrInvalidInput
	}

	exists, err := s.repos.User.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return apperrors.ErrCheckingFaild
	}
	if exists {
		return apperrors.ErrAlreadyExists
	}

	passwordHash, err := s.passwords.Hash(ctx, input.PasswordHash)
	if err != nil {
		return err
	}
	input.PasswordHash = passwordHash

	return s.repos.User.Create(ctx, input)
}

func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error) {
	user, err := s.repos.User.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrDataBase
	}
	return user.ToAPI(), nil
}

func (s *userService) UpdateUser(ctx context.Context, userID uuid.UUID, input *model.UserRequest) error {
	input.ID = userID
	return s.repos.User.Update(ctx, input)
}

func (s *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.repos.User.Delete(ctx, userID)
}

func (s *userService) ListUsers(ctx context.Context, params *model.UserParams) ([]*model.UserResponse, int, error) {

	//Validate params
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	users, err := s.repos.User.List(ctx, params)
	if err != nil {
		return nil, 0, errors.ErrDataBase
	}

	count, err := s.repos.User.Count(ctx, params)
	if err != nil {
		return nil, 0, errors.ErrDataBase
	}
	return users.ToAPIResponse(), count, nil
}
