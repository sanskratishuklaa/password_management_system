package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"password_manager/internal/model"
)

type VaultRepository struct {
	db *pgx.Conn
}

func NewVaultRepository(db *pgx.Conn) *VaultRepository {
	return &VaultRepository{
		db: db,
	}
}

func (r *VaultRepository) CreateVaultItem(
	ctx context.Context,
	item *model.VaultItem,
) (*model.VaultItem, error) {

	query := `
		INSERT INTO vault_items (
			user_id,
			title,
			username,
			encrypted_password,
			website,
			notes
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			user_id,
			title,
			username,
			encrypted_password,
			website,
			notes,
			created_at,
			updated_at
	`

	createdItem := &model.VaultItem{}

	err := r.db.QueryRow(
		ctx,
		query,
		item.UserID,
		item.Title,
		item.Username,
		item.EncryptedPassword,
		item.Website,
		item.Notes,
	).Scan(
		&createdItem.ID,
		&createdItem.UserID,
		&createdItem.Title,
		&createdItem.Username,
		&createdItem.EncryptedPassword,
		&createdItem.Website,
		&createdItem.Notes,
		&createdItem.CreatedAt,
		&createdItem.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return createdItem, nil
}

func (r *VaultRepository) GetVaultItemsByUserID(
	ctx context.Context,
	userID string,
) ([]*model.VaultItem, error) {

	query := `
		SELECT
			id,
			user_id,
			title,
			username,
			encrypted_password,
			website,
			notes,
			created_at,
			updated_at
		FROM vault_items
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var vaultItems []*model.VaultItem

	for rows.Next() {
		item := &model.VaultItem{}

		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Title,
			&item.Username,
			&item.EncryptedPassword,
			&item.Website,
			&item.Notes,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		vaultItems = append(vaultItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vaultItems, nil
}

func (r *VaultRepository) GetVaultItemByID(
	ctx context.Context,
	itemID string,
	userID string,
) (*model.VaultItem, error) {

	query := `
		SELECT
			id,
			user_id,
			title,
			username,
			encrypted_password,
			website,
			notes,
			created_at,
			updated_at
		FROM vault_items
		WHERE id = $1
		AND user_id = $2
	`

	item := &model.VaultItem{}

	err := r.db.QueryRow(
		ctx,
		query,
		itemID,
		userID,
	).Scan(
		&item.ID,
		&item.UserID,
		&item.Title,
		&item.Username,
		&item.EncryptedPassword,
		&item.Website,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}
func (r *VaultRepository) UpdateVaultItem(
	ctx context.Context,
	itemID string,
	userID string,
	title string,
	username string,
	encryptedPassword string,
	website *string,
	notes *string,
) (*model.VaultItem, error) {

	query := `
		UPDATE vault_items
		SET
			title = $1,
			username = $2,
			encrypted_password = $3,
			website = $4,
			notes = $5,
			updated_at = now()
		WHERE id = $6
		AND user_id = $7
		RETURNING
			id,
			user_id,
			title,
			username,
			encrypted_password,
			website,
			notes,
			created_at,
			updated_at
	`

	item := &model.VaultItem{}

	err := r.db.QueryRow(
		ctx,
		query,
		title,
		username,
		encryptedPassword,
		website,
		notes,
		itemID,
		userID,
	).Scan(
		&item.ID,
		&item.UserID,
		&item.Title,
		&item.Username,
		&item.EncryptedPassword,
		&item.Website,
		&item.Notes,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return item, nil
}
func (r *VaultRepository) DeleteVaultItem(
	ctx context.Context,
	itemID string,
	userID string,
) error {

	query := `
		DELETE FROM vault_items
		WHERE id = $1
		AND user_id = $2
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		itemID,
		userID,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("vault item not found")
	}

	return nil
}
