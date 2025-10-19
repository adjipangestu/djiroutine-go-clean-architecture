package auth

import (
	"context"
)

// User adalah tipe data yang mewakili informasi pengguna yang diautentikasi
type User struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

// UseCase adalah interface untuk use case autentikasi
type UseCase interface {
	// ValidateToken memvalidasi token dan mengembalikan informasi pengguna
	ValidateToken(ctx context.Context, token string) (*User, error)

	// GetAuthorizationURL menghasilkan URL otorisasi untuk login
	GetAuthorizationURL(ctx context.Context) (string, string, error)

	// ProcessCallback memproses callback dari OAuth provider
	ProcessCallback(ctx context.Context, code, state string) (*User, string, error)

	// Logout mengeluarkan pengguna dari sistem
	Logout(ctx context.Context, token string) error
}
