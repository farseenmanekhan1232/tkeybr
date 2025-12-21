package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	deviceAuthURL = "https://oauth2.googleapis.com/device/code"
	tokenURL      = "https://oauth2.googleapis.com/token"
	userInfoURL   = "https://www.googleapis.com/oauth2/v2/userinfo"
)

// DeviceCodeResponse represents the response from the device code request
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// TokenResponse represents the OAuth token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

// UserInfo represents Google user info
type UserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// StoredAuth represents stored authentication data
type StoredAuth struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

// GoogleAuth handles Google OAuth device flow authentication
type GoogleAuth struct {
	clientID     string
	storagePath  string
	currentAuth  *StoredAuth
	httpClient   *http.Client
}

// NewGoogleAuth creates a new Google auth handler
func NewGoogleAuth(clientID string) *GoogleAuth {
	homeDir, _ := os.UserHomeDir()
	storagePath := filepath.Join(homeDir, ".tkeybr", "auth.json")

	return &GoogleAuth{
		clientID:    clientID,
		storagePath: storagePath,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
	}
}

// RequestDeviceCode initiates the device code flow
func (g *GoogleAuth) RequestDeviceCode() (*DeviceCodeResponse, error) {
	data := url.Values{}
	data.Set("client_id", g.clientID)
	data.Set("scope", "openid email profile")

	resp, err := g.httpClient.PostForm(deviceAuthURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to request device code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var dcResp DeviceCodeResponse
	if err := json.Unmarshal(body, &dcResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &dcResp, nil
}

// PollForToken polls for the access token after user authorization
func (g *GoogleAuth) PollForToken(deviceCode string, interval int) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", g.clientID)
	data.Set("device_code", deviceCode)
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	resp, err := g.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to poll token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &tokenResp, nil
}

// GetUserInfo fetches user info using the access token
func (g *GoogleAuth) GetUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo UserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}

// SaveAuth saves authentication data to disk
func (g *GoogleAuth) SaveAuth(auth *StoredAuth) error {
	// Ensure directory exists
	dir := filepath.Dir(g.storagePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create auth directory: %w", err)
	}

	data, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal auth data: %w", err)
	}

	if err := os.WriteFile(g.storagePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write auth file: %w", err)
	}

	g.currentAuth = auth
	return nil
}

// LoadAuth loads authentication data from disk
func (g *GoogleAuth) LoadAuth() (*StoredAuth, error) {
	data, err := os.ReadFile(g.storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read auth file: %w", err)
	}

	var auth StoredAuth
	if err := json.Unmarshal(data, &auth); err != nil {
		return nil, fmt.Errorf("failed to parse auth file: %w", err)
	}

	g.currentAuth = &auth
	return &auth, nil
}

// IsTokenValid checks if the current token is still valid
func (g *GoogleAuth) IsTokenValid() bool {
	if g.currentAuth == nil {
		return false
	}
	return time.Now().Before(g.currentAuth.ExpiresAt)
}

// RefreshAccessToken refreshes the access token using the refresh token
func (g *GoogleAuth) RefreshAccessToken() error {
	if g.currentAuth == nil || g.currentAuth.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	data := url.Values{}
	data.Set("client_id", g.clientID)
	data.Set("refresh_token", g.currentAuth.RefreshToken)
	data.Set("grant_type", "refresh_token")

	resp, err := g.httpClient.PostForm(tokenURL, data)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if tokenResp.Error != "" {
		return fmt.Errorf("token refresh failed: %s", tokenResp.ErrorDesc)
	}

	g.currentAuth.AccessToken = tokenResp.AccessToken
	g.currentAuth.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return g.SaveAuth(g.currentAuth)
}

// Logout clears stored authentication
func (g *GoogleAuth) Logout() error {
	g.currentAuth = nil
	if err := os.Remove(g.storagePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove auth file: %w", err)
	}
	return nil
}

// GetCurrentUser returns the current authenticated user, or nil if not logged in
func (g *GoogleAuth) GetCurrentUser() *UserInfo {
	if g.currentAuth == nil {
		return nil
	}
	return &g.currentAuth.User
}

// completeAuth completes the authentication flow with the token response
func (g *GoogleAuth) CompleteAuth(tokenResp *TokenResponse) error {
	userInfo, err := g.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		return err
	}

	auth := &StoredAuth{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		User:         *userInfo,
	}

	return g.SaveAuth(auth)
}

// Placeholder for when no client ID is configured
func (g *GoogleAuth) IsConfigured() bool {
	return g.clientID != "" && g.clientID != "YOUR_GOOGLE_CLIENT_ID"
}

// Helper to read response body
func readBody(resp *http.Response) ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(resp.Body)
	return buf.Bytes(), err
}
