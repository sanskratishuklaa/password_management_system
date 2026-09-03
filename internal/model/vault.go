package model

import "time"

type VaultItem struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	Title             string    `json:"title"`
	Username          string    `json:"username"`
	EncryptedPassword string    `json:"-"`
	Website           string    `json:"website,omitempty"`
	Notes             string    `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CreateVaultRequest struct {
	Title    string `json:"title"`
	Username string `json:"username"`
	Password string `json:"password"`
	Website  string `json:"website"`
	Notes    string `json:"notes"`
}