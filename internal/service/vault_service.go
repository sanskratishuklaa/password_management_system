package service

import (
	"context"
	"fmt"
	"strings"

	"password_manager/internal/model"
	"password_manager/internal/repository"
	"password_manager/internal/security"
)

var (
	ErrVaultItemNotFound = fmt.Errorf("vault item not found")
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
) (*model.VaultItemResponse, error) {

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

	return model.ToVaultItemResponse(
		createdItem,
		request.Password,
	), nil
}

func (s *VaultService) GetVaultItems(
	ctx context.Context,
	userID string,
) ([]*model.VaultItemResponse, error) {

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

	responses := make(
		[]*model.VaultItemResponse,
		0,
		len(vaultItems),
	)

	for _, vaultItem := range vaultItems {

		password, err := security.Decrypt(
			vaultItem.EncryptedPassword,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to decrypt vault item password: %w",
				err,
			)
		}

		response := model.ToVaultItemResponse(
			vaultItem,
			password,
		)

		responses = append(responses, response)
	}

	return responses, nil
}

func (s *VaultService) GetVaultItemByID(
	ctx context.Context,
	itemID string,
	userID string,
) (*model.VaultItemResponse, error) {

	if strings.TrimSpace(itemID) == "" {
		return nil, fmt.Errorf("vault item ID is required")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	vaultItem, err := s.vaultRepository.GetVaultItemByID(
		ctx,
		itemID,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get vault item: %w", err)
	}

	password, err := security.Decrypt(
		vaultItem.EncryptedPassword,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decrypt password: %w",
			err,
		)
	}

	return model.ToVaultItemResponse(
		vaultItem,
		password,
	), nil
}

func (s *VaultService) UpdateVaultItem(
	ctx context.Context,
	itemID string,
	userID string,
	title string,
	username string,
	password string,
	website *string,
	notes *string,
) (*model.VaultItemResponse, error) {

	if strings.TrimSpace(itemID) == "" {
		return nil, fmt.Errorf("vault item ID is required")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title is required")
	}

	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("username is required")
	}

	if password == "" {
		return nil, fmt.Errorf("password is required")
	}

	encryptedPassword, err := security.Encrypt(password)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to encrypt password: %w",
			err,
		)
	}

	vaultItem, err := s.vaultRepository.UpdateVaultItem(
		ctx,
		itemID,
		userID,
		title,
		username,
		encryptedPassword,
		website,
		notes,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to update vault item: %w",
			err,
		)
	}

	return model.ToVaultItemResponse(
		vaultItem,
		password,
	), nil
}

func (s *VaultService) DeleteVaultItem(
	ctx context.Context,
	itemID string,
	userID string,
) error {

	if strings.TrimSpace(itemID) == "" {
		return fmt.Errorf("vault item ID is required")
	}

	if strings.TrimSpace(userID) == "" {
		return fmt.Errorf("user ID is required")
	}

	err := s.vaultRepository.DeleteVaultItem(
		ctx,
		itemID,
		userID,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to delete vault item: %w",
			err,
		)
	}

	return nil
}
