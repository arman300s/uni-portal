package services

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/arman300s/uni-portal/internal/core/contracts"
	"github.com/arman300s/uni-portal/internal/core/repositories"
	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/auth"
	"github.com/arman300s/uni-portal/pkg/queue"
	"github.com/arman300s/uni-portal/pkg/tasks"
)

// AuthService coordinates signup/login flows.
type AuthService struct {
	users repositories.UserRepository
	roles repositories.RoleRepository
}

func NewAuthService(users repositories.UserRepository, roles repositories.RoleRepository) *AuthService {
	return &AuthService{users: users, roles: roles}
}

func (s *AuthService) Signup(ctx context.Context, input contracts.SignupInput) (*contracts.AuthResponse, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Name = strings.TrimSpace(input.Name)

	if errs := validateSignupInput(input); len(errs) > 0 {
		return nil, errs
	}

	if _, err := s.users.FindByEmail(ctx, input.Email); err == nil {
		return nil, contracts.ErrEmailInUse
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashed,
	}
	if err := s.users.Create(ctx, &user); err != nil {
		return nil, err
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	payload := tasks.SendWelcomeEmailPayload{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	}
	_ = queue.Enqueue(tasks.TypeSendWelcomeEmail, payload, 0)

	return &contracts.AuthResponse{Token: token, UserID: user.ID}, nil
}

func (s *AuthService) Login(ctx context.Context, input contracts.LoginInput) (*contracts.AuthResponse, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))

	if errs := validateLoginInput(input); len(errs) > 0 {
		return nil, errs
	}

	user, err := s.users.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contracts.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := auth.CheckPassword(user.Password, input.Password); err != nil {
		return nil, contracts.ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	payload := tasks.SendWelcomeEmailPayload{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	}
	_ = queue.Enqueue(tasks.TypeSendWelcomeEmail, payload, 0)

	return &contracts.AuthResponse{Token: token, UserID: user.ID}, nil
}
