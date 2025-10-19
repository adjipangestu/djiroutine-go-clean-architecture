package sso

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// Environment represents the OAuth2 environment
type Environment string

const (
	Sandbox    Environment = "sandbox"
	Production Environment = "production"
)

// Config holds OAuth2 configuration constants
type Config struct {
	SandboxBaseURL    string
	ProductionBaseURL string
	DefaultRedisURL   string
	TokenEndpoint     string
	AuthorizeEndpoint string
	UserInfoEndpoint  string
	RevokeEndpoint    string
	UnauthorizedURL   string
}

// DefaultConfig provides default configuration values
var DefaultConfig = Config{
	SandboxBaseURL:    "https://sandbox.sso.example.com",
	ProductionBaseURL: "https://sso.example.com",
	DefaultRedisURL:   "redis://:1n1p45w0rd@127.0.0.1:6379/1", // Default Redis URL; replace with environment variable if needed
	TokenEndpoint:     "/o/token/",
	AuthorizeEndpoint: "/o/authorize/",
	UserInfoEndpoint:  "/o/userinfo/",
	RevokeEndpoint:    "/o/revoke-token/",
	UnauthorizedURL:   "/o/unauthorized/",
}

// OAuth2Client represents the OAuth2 client
type OAuth2Client struct {
	clientID     string
	clientSecret string
	redirectURI  string
	environment  Environment
	baseURL      string
	redisClient  *redis.Client
}

// NewOAuth2Client creates a new OAuth2 client instance
func NewOAuth2Client(clientID, clientSecret, redirectURI string, environment Environment, redisURL string) (*OAuth2Client, error) {
	if redisURL == "" {
		redisURL = DefaultConfig.DefaultRedisURL
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %v", err)
	}

	redisClient := redis.NewClient(opt)
	ctx := context.Background()

	// Test Redis connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	baseURL := DefaultConfig.SandboxBaseURL
	if environment == Production {
		baseURL = DefaultConfig.ProductionBaseURL
	}

	return &OAuth2Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		environment:  environment,
		baseURL:      baseURL,
		redisClient:  redisClient,
	}, nil
}

// GenerateCodeVerifier generates a code verifier for PKCE
func (c *OAuth2Client) GenerateCodeVerifier() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// GenerateCodeChallenge generates a code challenge from the code verifier
func (c *OAuth2Client) GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

type AuthorizationURL struct {
	URL string `json:"url"`
}

// GetAuthorizationURL returns the authorization URL for initiating the OAuth2 flow
func (c *OAuth2Client) GetAuthorizationURL(state string) (*AuthorizationURL, error) {
	verifier, err := c.GenerateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %v", err)
	}

	challenge := c.GenerateCodeChallenge(verifier)

	ctx := context.Background()
	key := fmt.Sprintf("oauth2_verifier_%s", state)

	err = c.redisClient.Set(ctx, key, verifier, 10*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to store code verifier: %v", err)
	}

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", c.clientID)
	params.Set("redirect_uri", c.redirectURI)
	params.Set("state", state)
	params.Set("scope", "read write profile email")
	params.Set("code_challenge", challenge)
	params.Set("code_challenge_method", "S256")

	authURL := fmt.Sprintf("%s%s?%s", c.baseURL, DefaultConfig.AuthorizeEndpoint, params.Encode())

	return &AuthorizationURL{
		URL: authURL,
	}, nil
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// GetAccessToken exchanges authorization code for access token
func (c *OAuth2Client) GetAccessToken(code, state string) (*TokenResponse, error) {
	ctx := context.Background()
	key := fmt.Sprintf("oauth2_verifier_%s", state)

	verifier, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("code verifier not found or expired: %v", err)
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", c.redirectURI)
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("code_verifier", verifier)

	resp, err := http.PostForm(c.baseURL+DefaultConfig.TokenEndpoint, data)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %v", err)
	}

	// Clean up Redis
	c.redisClient.Del(ctx, key)

	return &tokenResp, nil
}

// UserInfo represents the user information response
type UserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

// GetUserInfo retrieves user information using the access token
func (c *OAuth2Client) GetUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", c.baseURL+DefaultConfig.UserInfoEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	fmt.Println(req)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return &userInfo, nil
}

// IsAuthenticated checks if the user is authenticated using the access token
func (c *OAuth2Client) IsAuthenticated(accessToken string) bool {
	userInfo, err := c.GetUserInfo(accessToken)
	if err != nil {
		return false
	}
	return userInfo.Sub != ""
}

// Logout revokes the access token and logs out the user
func (c *OAuth2Client) Logout(accessToken, sessionKey string) error {
	data := url.Values{}
	data.Set("pjnhk_id", sessionKey)

	req, err := http.NewRequest("POST", c.baseURL+DefaultConfig.RevokeEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create logout request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("logout request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed with status: %d", resp.StatusCode)
	}

	return nil
}

// GetUnauthorizedURL returns the unauthorized URL
func (c *OAuth2Client) GetUnauthorizedURL() string {
	params := url.Values{}
	params.Set("client_id", c.clientID)
	return fmt.Sprintf("%s%s?%s", c.baseURL, DefaultConfig.UnauthorizedURL, params.Encode())
}
