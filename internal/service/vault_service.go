package service

import (
	"context"
	"fmt"
	"strings"

	"password_manager/internal/model"
	"password_manager/internal/repository"
	"password_manager/internal/security"
)

type VaultService struct {
	vaultRepository *repository.VaultRepository
}

func NewVaultService(
	vaultRepository *repository.VaultRepository,
) *VaultService {
	return &VaultService{
		vaultRepository: vaultRepository,
	}
}

func (s *VaultService) CreateVaultItem(
	ctx context.Context,
	userID string,
	request model.CreateVaultRequest,
) (*model.VaultItem, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	if strings.TrimSpace(request.Title) == "" {
		return nil, fmt.Errorf("title is required")
	}

	if strings.TrimSpace(request.Username) == "" {
		return nil, fmt.Errorf("username is required")
	}

	if request.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	encryptedPassword, err := security.Encrypt(request.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	vaultItem := &model.VaultItem{
		UserID:            userID,
		Title:             request.Title,
		Username:          request.Username,
		EncryptedPassword: encryptedPassword,
		Website:           request.Website,
		Notes:             request.Notes,
	}

	createdItem, err := s.vaultRepository.CreateVaultItem(
		ctx,
		vaultItem,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault item: %w", err)
	}

	return createdItem, nil
}
func (s *VaultService) GetVaultItems(
	ctx context.Context,
	userID string,
) ([]*model.VaultItem, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	vaultItems, err := s.vaultRepository.GetVaultItemsByUserID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get vault items: %w", err)
	}

	return vaultItems, nil
}
