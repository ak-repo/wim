package service

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/ak-repo/wim/internal/constants"
	apperrors "github.com/ak-repo/wim/internal/errs"
	"github.com/ak-repo/wim/internal/model"
	"github.com/ak-repo/wim/internal/repository"
	"github.com/ak-repo/wim/pkg/auth"
)

type AuthService interface {
	Register(ctx context.Context, input *model.RegisterRequest) error
	Login(ctx context.Context, input *model.LoginRequest) (*model.AuthResponse, error)
	Me(ctx context.Context) (*model.UserResponse, error)
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
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if len(strings.TrimSpace(input.Password)) < 8 || !strings.Contains(input.Email, "@") {
		return apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	exists, err := s.repos.User.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeCheckFailed, "failed to check email availability")
	}
	if exists {
		return apperrors.New(apperrors.CodeAlreadyExists, "user with this email already exists")
	}

	passwordHash, err := s.passwords.Hash(ctx, input.Password)
	if err != nil {
		return apperrors.Wrap(err, apperrors.CodeInternal, "failed to process password")
	}

	username := strings.TrimSpace(input.Username)
	email := input.Email
	role := constants.RoleAdmin
	isActive := true
	user := &model.UserRequest{
		Username:     &username,
		Email:        &email,
		PasswordHash: &passwordHash,
		Role:         &role,
		IsActive:     &isActive,
	}

	// Refcode
	refCode, err := s.repos.RefCode.GenerateUserRefCode(ctx)
	if err != nil {
		return err
	}

	user.RefCode = refCode
	_, err = s.repos.User.Create(ctx, user)
	if err != nil {
		log.Println("error creating user", err)
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return apperrors.New(apperrors.CodeAlreadyExists, "user with this email already exists")
		}
		return apperrors.Wrap(err, apperrors.CodeDatabase, "failed to create user")
	}

	return nil
}

func (s *authService) Login(ctx context.Context, input *model.LoginRequest) (*model.AuthResponse, error) {
	if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}
	if !strings.Contains(input.Email, "@") {
		return nil, apperrors.New(apperrors.CodeInvalidInput, "invalid input")
	}

	user, err := s.repos.User.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.New(apperrors.CodeUnauthorized, "invalid email or password")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load user")
	}

	if !user.IsActive {
		return nil, apperrors.New(apperrors.CodeForbidden, "account is disabled")
	}

	if err := s.passwords.Compare(ctx, user.PasswordHash, input.Password); err != nil {
		return nil, apperrors.New(apperrors.CodeUnauthorized, "invalid email or password")
	}

	return s.issueTokens(ctx, user)
}

func (s *authService) issueTokens(ctx context.Context, user *model.UserDTO) (*model.AuthResponse, error) {
	accessToken, err := s.tokenManager.IssueAccessToken(ctx, auth.Claims{
		Subject: user.ID,
		Role:    user.Role,
	})
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeInternal, "failed to issue access token")
	}

	refreshToken, refreshTokenTTL, err := s.tokenManager.IssueRefreshToken(ctx, auth.Claims{Subject: user.ID, Role: user.Role})
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeInternal, "failed to issue refresh token")
	}

	refCode, err := s.repos.RefCode.GenerateRefreshTokenRefCode(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.CodeRefCodeFailed, "failed to generate refresh token reference")
	}

	now := time.Now().UTC()
	storedToken := &model.RefreshTokenDTO{
		RefCode:   refCode,
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

func (s *authService) Me(ctx context.Context) (*model.UserResponse, error) {
	userID, ok := auth.UserIDFromContext(ctx)
	if !ok || userID <= 0 {
		return nil, apperrors.New(apperrors.CodeUnauthorized, "unauthorized")
	}

	user, err := s.repos.User.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, apperrors.New(apperrors.CodeNotFound, "user not found")
		}
		return nil, apperrors.Wrap(err, apperrors.CodeDatabase, "failed to load user")
	}

	return user.ToAPIResponse(), nil
}
