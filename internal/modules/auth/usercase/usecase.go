package usecase

import (
	"context"
	"crypto/rand"
	"djiroutine-go-clean-architecture/internal/modules/auth"
	"djiroutine-go-clean-architecture/pkg/errors"
	"djiroutine-go-clean-architecture/pkg/sso"
	"encoding/base64"
)

type authUseCase struct {
	oauthClient *sso.OAuth2Client
}

// NewAuthUseCase membuat instance baru dari auth use case
func NewAuthUseCase(oauthClient *sso.OAuth2Client) auth.UseCase {
	return &authUseCase{
		oauthClient: oauthClient,
	}
}

// ValidateToken memvalidasi token dan mengembalikan informasi pengguna
func (uc *authUseCase) ValidateToken(ctx context.Context, token string) (*auth.User, error) {
	// Verifikasi token menggunakan SDK OAuth
	if !uc.oauthClient.IsAuthenticated(token) {
		return nil, errors.AuthError("Invalid or expired token", nil)
	}

	// Dapatkan informasi user dari token
	userInfo, err := uc.oauthClient.GetUserInfo(token)
	if err != nil {
		return nil, errors.InternalServerError("Failed to get user information", err)
	}

	// Konversi dari sso.UserInfo ke auth.User
	user := &auth.User{
		ID:      userInfo.Sub,
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Profile: userInfo.Profile,
	}

	return user, nil
}

// generateRandomState menghasilkan string acak untuk state
func generateRandomState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GetAuthorizationURL menghasilkan URL otorisasi untuk login
func (uc *authUseCase) GetAuthorizationURL(ctx context.Context) (string, string, error) {
	// Generate random state
	state, err := generateRandomState()
	if err != nil {
		return "", "", errors.InternalServerError("Failed to generate state", err)
	}

	// Get authorization URL from OAuth client
	authURL, err := uc.oauthClient.GetAuthorizationURL(state)
	if err != nil {
		return "", "", errors.InternalServerError("Failed to get authorization URL", err)
	}

	return authURL.URL, state, nil
}

// ProcessCallback memproses callback dari OAuth provider
func (uc *authUseCase) ProcessCallback(ctx context.Context, code, state string) (*auth.User, string, error) {
	// Exchange authorization code for access token
	tokenResp, err := uc.oauthClient.GetAccessToken(code, state)
	if err != nil {
		return nil, "", errors.AuthError("Failed to get access token", err)
	}

	// Get user info using the access token
	userInfo, err := uc.oauthClient.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, "", errors.InternalServerError("Failed to get user info", err)
	}

	// Convert to auth.User
	user := &auth.User{
		ID:      userInfo.Sub,
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Profile: userInfo.Profile,
	}

	return user, tokenResp.AccessToken, nil
}

// Logout mengeluarkan pengguna dari sistem
func (uc *authUseCase) Logout(ctx context.Context, token string) error {
	// Extract user ID from token to use as session key
	userInfo, err := uc.oauthClient.GetUserInfo(token)
	if err != nil {
		return errors.InternalServerError("Failed to get user info for logout", err)
	}

	// Use user ID as session key
	sessionKey := userInfo.Sub

	err = uc.oauthClient.Logout(token, sessionKey)
	if err != nil {
		return errors.InternalServerError("Failed to logout", err)
	}

	return nil
}
