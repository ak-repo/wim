package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	apperrors "github.com/ak-repo/wim/pkg/errors"
)

type UserRoleService interface {
	CreateUserRole(ctx context.Context, input *model.UserRoleRequest) (int, error)
	GetUserRoleByID(ctx context.Context, roleID int) (*model.UserRoleResponse, error)
	UpdateUserRole(ctx context.Context, roleID int, input *model.UserRoleRequest) error
	DeleteUserRole(ctx context.Context, roleID int) error
	ListUserRoles(ctx context.Context, params *model.UserRoleParams) ([]*model.UserRoleResponse, int, error)
}

type userRoleService struct {
	repos *repository.Repositories
}

func NewUserRoleService(repositories *repository.Repositories) UserRoleService {
	return &userRoleService{repos: repositories}
}

func (s *userRoleService) CreateUserRole(ctx context.Context, input *model.UserRoleRequest) (int, error) {
	if input == nil {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.Name == nil || strings.TrimSpace(*input.Name) == "" {
		return 0, apperrors.New(apperrors.CodeInvalidInput, "name is required")
	}

	exists, err := s.repos.UserRole.ExistsByName(ctx, *input.Name)
	if err != nil {
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check role name")
	}
	if exists {
		return 0, apperrors.New(apperrors.CodeAlreadyExists, "user role with this name already exists")
	}

	refCode, err := s.repos.RefCode.GenerateUserRoleRefCode(ctx)
	if err != nil {
		return 0, err
	}
	input.RefCode = refCode

	id, err := s.repos.UserRole.Create(ctx, input)
	if err != nil {
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return 0, apperrors.New(apperrors.CodeAlreadyExists, "user role with this name already exists")
		}
		return 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create user role")
	}

	return id, nil
}

func (s *userRoleService) GetUserRoleByID(ctx context.Context, roleID int) (*model.UserRoleResponse, error) {
	role, err := s.repos.UserRole.GetByID(ctx, roleID)
	if err != nil {
		if errors.Is(err, repository.ErrUserRoleNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "user role not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load user role")
	}
	return role.ToAPIResponse(), nil
}

func (s *userRoleService) UpdateUserRole(ctx context.Context, roleID int, input *model.UserRoleRequest) error {
	if input == nil {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	if input.Name != nil {
		if strings.TrimSpace(*input.Name) == "" {
			return apperrors.New(apperrors.CodeInvalidInput, "name cannot be empty")
		}
		exists, err := s.repos.UserRole.ExistsByName(ctx, *input.Name)
		if err != nil {
			return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to check role name")
		}
		if exists {
			return apperrors.New(apperrors.CodeAlreadyExists, "user role with this name already exists")
		}
	}

	if err := s.repos.UserRole.Update(ctx, roleID, input); err != nil {
		if errors.Is(err, repository.ErrUserRoleNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "user role not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to update user role")
	}

	return nil
}

func (s *userRoleService) DeleteUserRole(ctx context.Context, roleID int) error {
	if err := s.repos.UserRole.Delete(ctx, roleID); err != nil {
		if errors.Is(err, repository.ErrUserRoleNotFound) {
			return apperrors.New(apperrors.CodeNotFound, "user role not found")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to delete user role")
	}
	return nil
}

func (s *userRoleService) ListUserRoles(ctx context.Context, params *model.UserRoleParams) ([]*model.UserRoleResponse, int, error) {
	if params == nil {
		params = &model.UserRoleParams{}
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	roles, err := s.repos.UserRole.List(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to list user roles")
	}

	count, err := s.repos.UserRole.Count(ctx, params)
	if err != nil {
		return nil, 0, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to count user roles")
	}

	return roles.ToAPIResponse(), count, nil
}
