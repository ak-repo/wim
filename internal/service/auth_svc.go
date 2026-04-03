package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/ak-repo/wim/internal/constants"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/pkg/auth"
	apperrors "github.com/ak-repo/wim/pkg/errors"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, input *model.RegisterRequest) error
	Login(ctx context.Context, input *model.LoginRequest) (*model.AuthResponse, error)
}
type authService struct {
	repos        *repository.Repositories
	tokenManager auth.TokenManager
	passwords    auth.PasswordHasher
}

func NewAuthService(repositories *repository.Repositories, tokenManager auth.TokenManager, passwords auth.PasswordHasher,
) AuthService {
	return &authService{
		repos:        repositories,
		tokenManager: tokenManager,
		passwords:    passwords,
	}
}

func (s *authService) Register(ctx context.Context, input *model.RegisterRequest) error {
	if strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		return apperrors.ErrInvalidInput
	}
	if len(strings.TrimSpace(input.Password)) < 8 || !strings.Contains(input.Email, "@") {
		return apperrors.ErrInvalidInput
	}

	exists, err := s.repos.User.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return apperrors.ErrCheckingFaild
	}
	if exists {
		return apperrors.ErrAlreadyExists
	}

	passwordHash, err := s.passwords.Hash(ctx, input.Password)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	user := &model.UserRequest{
		ID:           uuid.New(),
		Username:     strings.TrimSpace(input.Username),
		Email:        input.Email,
		PasswordHash: passwordHash,
		Role:         constants.RoleWorker,
		IsActive:     true,
		UpdatedAt:    now,
		CreatedAt:    now,
	}

	if err := s.repos.User.Create(ctx, user); err != nil {
		log.Println("error creating user", err)
		return err
	}

	return nil
}

func (s *authService) Login(ctx context.Context, input *model.LoginRequest) (*model.AuthResponse, error) {
	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		return nil, apperrors.ErrInvalidInput
	}
	if !strings.Contains(input.Email, "@") {
		return nil, apperrors.ErrInvalidInput
	}

	user, err := s.repos.User.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.ErrUnauthorized
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, apperrors.ErrForbidden
	}

	if err := s.passwords.Compare(ctx, user.PasswordHash.String, input.Password); err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) issueTokens(ctx context.Context, user *model.UserDTO) (*model.AuthResponse, error) {
	accessToken, err := s.tokenManager.IssueAccessToken(ctx, auth.Claims{
		Subject: user.ID,
		Role:    user.Role.String,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenTTL, err := s.tokenManager.IssueRefreshToken(ctx, auth.Claims{Subject: user.ID, Role: user.Role.String})
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	storedToken := &model.RefreshTokenDTO{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshToken,
		ExpiresAt: now.Add(refreshTokenTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repos.Auth.StoreRefreshToken(ctx, storedToken); err != nil {
		log.Println("error storing refresh token", err)
		// return nil, err
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
