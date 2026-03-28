package service

import (
	"context"

	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	ListUsers(ctx context.Context) ([]*model.UserResponse, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, input *model.UserRequest) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
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
		return nil, err
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

func (s *userService) ListUsers(ctx context.Context) ([]*model.UserResponse, error) {
	users, err := s.repos.User.List(ctx, 100, 0)
	if err != nil {
		return nil, err
	}
	return users.ToAPIRequest(), nil
}
