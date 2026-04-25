package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/pkg/auth"
)

type UserService interface {
	CreateUser(ctx context.Context, input *model.UserRequest) (int, error)
	GetUserByID(ctx context.Context, userID int) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, userID int, input *model.UserRequest) error
	DeleteUser(ctx context.Context, userID int) error
	ListUsers(ctx context.Context, params *model.UserParams) ([]*model.UserResponse, int, error)
	GetUserCount(ctx context.Context, params *model.UserParams) (int, error)
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

const op = "service/UserService"

func (s *userService) CreateUser(ctx context.Context, input *model.UserRequest) (int, error) {
	const opCreate = op + ".Create"
	if input == nil {
		return 0, errs.E(opCreate, errs.InvalidRequest, errors.New("invalid input"))
	}

	if input.Username == nil || input.Email == nil || input.PasswordHash == nil {
		return 0, errs.E(opCreate, errs.InvalidRequest, errors.New("username, email and password are required"))
	}

	if strings.TrimSpace(*input.Username) == "" || strings.TrimSpace(*input.Email) == "" || strings.TrimSpace(*input.PasswordHash) == "" {
		return 0, errs.E(opCreate, errs.InvalidRequest, errors.New("username, email and password cannot be empty"))
	}

	if len(strings.TrimSpace(*input.PasswordHash)) < 8 {
		return 0, errs.E(opCreate, errs.InvalidRequest, errors.New("password must be at least 8 characters"))
	}

	if !strings.Contains(*input.Email, "@") {
		return 0, errs.E(opCreate, errs.InvalidRequest, errors.New("invalid email format"))
	}

	exists, err := s.repos.User.ExistsByEmail(ctx, *input.Email)
	if err != nil {
		return 0, errs.E(opCreate, errs.Database, err)
	}
	if exists {
		return 0, errs.E(opCreate, errs.Conflict, errors.New("user with this email already exists"))
	}

	exists, err = s.repos.User.ExistsByUsername(ctx, *input.Username)
	if err != nil {
		return 0, errs.E(opCreate, errs.Database, err)
	}
	if exists {
		return 0, errs.E(opCreate, errs.Conflict, errors.New("user with this username already exists"))
	}

	passwordHash, err := s.passwords.Hash(ctx, *input.PasswordHash)
	if err != nil {
		return 0, errs.E(opCreate, errs.Internal, err)
	}
	input.PasswordHash = &passwordHash

	refCode, err := s.repos.RefCode.GenerateUserRefCode(ctx)
	if err != nil {
		return 0, errs.E(opCreate, errs.Internal, err)
	}
	input.RefCode = refCode

	id, err := s.repos.User.Create(ctx, input)
	if err != nil {
		return 0, errs.E(opCreate, errs.Database, err)
	}
	return id, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID int) (*model.UserResponse, error) {
	const opGet = op + ".GetByID"
	user, err := s.repos.User.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(errs.TopError(err), repository.ErrUserNotFound) {
			return nil, errs.E(opGet, errs.NotFound, err)
		}
		return nil, errs.E(opGet, errs.Database, err)
	}
	return user.ToAPIResponse(), nil
}

func (s *userService) UpdateUser(ctx context.Context, userID int, input *model.UserRequest) error {
	const opUpdate = op + ".Update"
	if input == nil {
		return errs.E(opUpdate, errs.InvalidRequest, errors.New("invalid input"))
	}

	if input.Username != nil {
		if strings.TrimSpace(*input.Username) == "" {
			return errs.E(opUpdate, errs.InvalidRequest, errors.New("username cannot be empty"))
		}
		exists, err := s.repos.User.ExistsByUsername(ctx, *input.Username)
		if err != nil {
			return errs.E(opUpdate, errs.Database, err)
		}
		if exists {
			return errs.E(opUpdate, errs.Conflict, errors.New("user with this username already exists"))
		}
	}

	if input.Email != nil {
		if strings.TrimSpace(*input.Email) == "" {
			return errs.E(opUpdate, errs.InvalidRequest, errors.New("email cannot be empty"))
		}
		if !strings.Contains(*input.Email, "@") {
			return errs.E(opUpdate, errs.InvalidRequest, errors.New("invalid email format"))
		}
		exists, err := s.repos.User.ExistsByEmail(ctx, *input.Email)
		if err != nil {
			return errs.E(opUpdate, errs.Database, err)
		}
		if exists {
			return errs.E(opUpdate, errs.Conflict, errors.New("user with this email already exists"))
		}
	}

	if input.PasswordHash != nil {
		if strings.TrimSpace(*input.PasswordHash) == "" {
			return errs.E(opUpdate, errs.InvalidRequest, errors.New("password cannot be empty"))
		}
		if len(strings.TrimSpace(*input.PasswordHash)) < 8 {
			return errs.E(opUpdate, errs.InvalidRequest, errors.New("password must be at least 8 characters"))
		}
		passwordHash, err := s.passwords.Hash(ctx, *input.PasswordHash)
		if err != nil {
			return errs.E(opUpdate, errs.Internal, err)
		}
		input.PasswordHash = &passwordHash
	}

	if input.Role != nil {
		if strings.TrimSpace(*input.Role) == "" {
			return errs.E(opUpdate, errs.InvalidRequest, errors.New("role cannot be empty"))
		}
	}

	if err := s.repos.User.Update(ctx, userID, input); err != nil {
		if errors.Is(errs.TopError(err), repository.ErrUserNotFound) {
			return errs.E(opUpdate, errs.NotFound, err)
		}
		return errs.E(opUpdate, errs.Database, err)
	}
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, userID int) error {
	const opDelete = op + ".Delete"
	if err := s.repos.User.Delete(ctx, userID); err != nil {
		if errors.Is(errs.TopError(err), repository.ErrUserNotFound) {
			return errs.E(opDelete, errs.NotFound, err)
		}
		return errs.E(opDelete, errs.Database, err)
	}
	return nil
}

func (s *userService) ListUsers(ctx context.Context, params *model.UserParams) ([]*model.UserResponse, int, error) {
	const opList = op + ".List"
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	users, err := s.repos.User.List(ctx, params)
	if err != nil {
		return nil, 0, errs.E(opList, errs.Database, err)
	}

	count, err := s.repos.User.Count(ctx, params)
	if err != nil {
		return nil, 0, errs.E(opList, errs.Database, err)
	}
	return users.ToAPIResponse(), count, nil
}

func (s *userService) GetUserCount(ctx context.Context, params *model.UserParams) (int, error) {
	const opCount = op + ".GetUserCount"
	count, err := s.repos.User.Count(ctx, params)
	if err != nil {
		return 0, errs.E(opCount, errs.Database, err)
	}
	return count, nil
}
