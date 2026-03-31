package service

import (
	"context"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/pkg/errors"
	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, input *model.UserRequest) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	ListUsers(ctx context.Context, params *model.UserParams) ([]*model.UserResponse, int, error)
}

type userService struct {
	repos *repository.Repositories
}

func NewUserService(repositories *repository.Repositories) UserService {
	return &userService{
		repos: repositories,
	}
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
	return users.ToAPIRequest(), count, nil
}
