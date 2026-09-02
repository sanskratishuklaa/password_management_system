package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"password_manager/internal/model"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
) (*model.User, error) {

	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, created_at, updated_at
	`

	user := &model.User{}

	err := r.db.QueryRow(
		ctx,
		query,
		email,
		passwordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}