package service

import (
	"context"
	"fmt"

	"password_manager/internal/model"
	"password_manager/internal/repository"
	"password_manager/internal/security"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(
	userRepository *repository.UserRepository,
) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) RegisterUser(
	ctx context.Context,
	request model.RegisterRequest,
) (*model.User, error) {

	if request.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if request.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	passwordHash, err := security.HashPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.userRepository.CreateUser(
		ctx,
		request.Email,
		passwordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) LoginUser(
	ctx context.Context,
	request model.LoginRequest,
) (*model.User, error) {

	if request.Email == "" {
		return nil, fmt.Errorf("email is required")
	}

	if request.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	user, err := s.userRepository.GetUserByEmail(
		ctx,
		request.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	isValid, err := security.VerifyPassword(
		request.Password,
		user.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password: %w", err)
	}

	if !isValid {
		return nil, fmt.Errorf("invalid email or password")
	}

	return user, nil
}