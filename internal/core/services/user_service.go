package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/arman300s/uni-portal/internal/core/contracts"
	"github.com/arman300s/uni-portal/internal/core/repositories"
	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/auth"
	"github.com/arman300s/uni-portal/pkg/cache"
)

// UserService encapsulates admin/user flows.
type UserService struct {
	users repositories.UserRepository
	roles repositories.RoleRepository
}

func NewUserService(users repositories.UserRepository, roles repositories.RoleRepository) *UserService {
	return &UserService{users: users, roles: roles}
}

func (s *UserService) GetCurrentUser(ctx context.Context, id uint) (*contracts.UserDTO, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contracts.ErrUserNotFound
		}
		return nil, err
	}
	return mapToUserDTO(user), nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]contracts.UserDTO, error) {
	cacheKey := "users:all"
	cached, err := cache.RDB.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var dtos []contracts.UserDTO
		if json.Unmarshal(cached, &dtos) == nil {
			return dtos, nil
		}
	}

	users, err := s.users.List(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]contracts.UserDTO, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, *mapToUserDTO(&u))
	}

	data, _ := json.Marshal(dtos)
	cache.RDB.Set(ctx, cacheKey, data, 5*time.Minute)

	return dtos, nil
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*contracts.UserDTO, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contracts.ErrUserNotFound
		}
		return nil, err
	}
	return mapToUserDTO(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, input contracts.CreateUserInput) (*contracts.UserDTO, error) {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	input.RoleName = strings.TrimSpace(strings.ToLower(input.RoleName))

	if errs := validateCreateUserInput(input); len(errs) > 0 {
		return nil, errs
	}

	if _, err := s.users.FindByEmail(ctx, input.Email); err == nil {
		return nil, contracts.ErrEmailInUse
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role, err := s.roles.FindByName(ctx, input.RoleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, contracts.ErrRoleNotFound
		}
		return nil, err
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		RoleID:   &role.ID,
	}

	if err := s.users.Create(ctx, &user); err != nil {
		return nil, err
	}

	return mapToUserDTO(&user), nil
}

func (s *UserService) UpdateUserRole(ctx context.Context, id uint, input contracts.UpdateUserInput) error {
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.RoleName = strings.TrimSpace(strings.ToLower(input.RoleName))

	if errs := validateUpdateUserInput(input); len(errs) > 0 {
		return errs
	}

	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.ErrUserNotFound
		}
		return err
	}

	role, err := s.roles.FindByName(ctx, input.RoleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.ErrRoleNotFound
		}
		return err
	}

	user.RoleID = &role.ID
	return s.users.Save(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	if err := s.users.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contracts.ErrUserNotFound
		}
		return err
	}
	return nil
}

func mapToUserDTO(user *models.User) *contracts.UserDTO {
	dto := &contracts.UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	if user.Role != nil {
		dto.Role = user.Role.Name
	}
	return dto
}
