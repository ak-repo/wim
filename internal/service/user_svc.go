package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/pkg/auth"
	apperrors "github.com/ak-repo/wim/pkg/errors"
)

type UserService interface {
	CreateUser(ctx context.Context, input *model.UserRequest) (int, error)
	GetUserByID(ctx context.Context, userID int) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, userID int, input *model.UserRequest) error
	DeleteUser(ctx context.Context, userID int) error
	ListUsers(ctx context.Context, params *model.UserParams) ([]*model.UserResponse, int, error)
}

type userService struct {
	repos     *repository.Repositories
	passwords auth.PasswordHasher
}

func NewUserService(repositories *repository.Repositories, passwords auth.PasswordHasher) UserService {
	return &userService{
		repos:     repositories,
		passwords: passwords,
	}
}

func (s *userService) CreateUser(ctx context.Context, input *model.UserRequest) (int, error) {
	if input == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// Required field validation
	if input.Username == nil || input.Email == nil || input.PasswordHash == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "username, email and password are required")
	}

	if strings.TrimSpace(*input.Username) == "" || strings.TrimSpace(*input.Email) == "" || strings.TrimSpace(*input.PasswordHash) == "" {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "username, email and password cannot be empty")
	}

	if len(strings.TrimSpace(*input.PasswordHash)) < 8 {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "password must be at least 8 characters")
	}

	if !strings.Contains(*input.Email, "@") {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid email format")
	}

	// Check email uniqueness
	exists, err := s.repos.User.ExistsByEmail(ctx, *input.Email)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check email availability")
	}
	if exists {
		return 0, apperrors.New(apperrors.CodeAlreadyExists, "user with this email already exists")
	}

	// Check username uniqueness
	exists, err = s.repos.User.ExistsByUsername(ctx, *input.Username)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check username availability")
	}
	if exists {
		return 0, apperrors.New(apperrors.CodeAlreadyExists, "user with this username already exists")
	}

	passwordHash, err := s.passwords.Hash(ctx, *input.PasswordHash)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeInternal, "failed to process password")
	}
	input.PasswordHash = &passwordHash

	// Refcode
	refCode, err := s.repos.RefCode.GenerateUserRefCode(ctx)
	if err != nil {
		return 0, err
	}
	input.RefCode = refCode

	id, err := s.repos.User.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "user with this email already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create user")
	}
	return id, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID int) (*model.UserResponse, error) {
	user, err := s.repos.User.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "user not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load user")
	}
	return user.ToAPIResponse(), nil
}

func (s *userService) UpdateUser(ctx context.Context, userID int, input *model.UserRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	// PATCH validation - only validate provided fields
	if input.Username != nil {
		if strings.TrimSpace(*input.Username) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "username cannot be empty")
		}
		// Check username uniqueness if being updated
		exists, err := s.repos.User.ExistsByUsername(ctx, *input.Username)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check username availability")
		}
		if exists {
			return apperrors.New(apperrors.CodeAlreadyExists, "user with this username already exists")
		}
	}

	if input.Email != nil {
		if strings.TrimSpace(*input.Email) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "email cannot be empty")
		}
		if !strings.Contains(*input.Email, "@") {
			return apperrors.New(apperrors.CodeInvalidInput, "invalid email format")
		}
		// Check email uniqueness if being updated
		exists, err := s.repos.User.ExistsByEmail(ctx, *input.Email)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check email availability")
		}
		if exists {
			return apperrors.New(apperrors.CodeAlreadyExists, "user with this email already exists")
		}
	}

	if input.PasswordHash != nil {
		if strings.TrimSpace(*input.PasswordHash) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "password cannot be empty")
		}
		if len(strings.TrimSpace(*input.PasswordHash)) < 8 {
			return apperrors.New(apperrors.CodeInvalidInput, "password must be at least 8 characters")
		}
		passwordHash, err := s.passwords.Hash(ctx, *input.PasswordHash)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeInternal, "failed to process password")
		}
		input.PasswordHash = &passwordHash
	}

	if input.Role != nil {
		if strings.TrimSpace(*input.Role) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "role cannot be empty")
		}
	}

	if err := s.repos.User.Update(ctx, userID, input); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "user not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update user")
	}
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	if err := s.repos.User.Delete(ctx, userID); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "user not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete user")
	}
	return nil
}

func (s *userService) ListUsers(ctx context.Context, params *model.UserParams) ([]*model.UserResponse, int, error) {
	// Validate params
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	users, err := s.repos.User.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list users")
	}

	count, err := s.repos.User.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count users")
	}
	return users.ToAPIResponse(), count, nil
}
