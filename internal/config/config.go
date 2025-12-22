package config

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
)

// Build-time variables (set via ldflags)
// These are embedded during CI/CD builds from GitHub Secrets
var (
	// EmbeddedGoogleClientID is the Google OAuth Client ID (set at build time)
	EmbeddedGoogleClientID string

	// EmbeddedFirebaseCredentials is base64-encoded Firebase service account JSON (set at build time)
	EmbeddedFirebaseCredentials string

	// EmbeddedFirebaseProjectID is the Firebase project ID (set at build time)
	EmbeddedFirebaseProjectID string
)

// Config holds application configuration
type Config struct {
	// Firebase configuration
	FirebaseCredentialsPath string `json:"firebase_credentials_path"`
	FirebaseProjectID       string `json:"firebase_project_id"`

	// Google OAuth (for device code flow)
	GoogleClientID string `json:"google_client_id"`

	// Decoded Firebase credentials (from embedded or file)
	firebaseCredentialsJSON []byte
}

// DefaultConfigPath returns the default config file path
func DefaultConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".tkeybr", "config.json")
}

// DefaultFirebaseCredentialsPath returns the default Firebase credentials path
func DefaultFirebaseCredentialsPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".tkeybr", "firebase-credentials.json")
}

// Load loads configuration from embedded values, file, and environment
func Load() (*Config, error) {
	cfg := &Config{}

	// Priority 1: Embedded credentials (from build)
	if EmbeddedGoogleClientID != "" {
		cfg.GoogleClientID = EmbeddedGoogleClientID
	}
	if EmbeddedFirebaseProjectID != "" {
		cfg.FirebaseProjectID = EmbeddedFirebaseProjectID
	}
	if EmbeddedFirebaseCredentials != "" {
		decoded, err := base64.StdEncoding.DecodeString(EmbeddedFirebaseCredentials)
		if err == nil {
			cfg.firebaseCredentialsJSON = decoded
		}
	}

	// Priority 2: Config file (can override embedded)
	configPath := DefaultConfigPath()
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Priority 3: Environment variables (highest priority, can override all)
	if envCreds := os.Getenv("TKEYBR_FIREBASE_CREDENTIALS"); envCreds != "" {
		cfg.FirebaseCredentialsPath = envCreds
	}
	if envProject := os.Getenv("TKEYBR_FIREBASE_PROJECT"); envProject != "" {
		cfg.FirebaseProjectID = envProject
	}
	if envClientID := os.Getenv("TKEYBR_GOOGLE_CLIENT_ID"); envClientID != "" {
		cfg.GoogleClientID = envClientID
	}

	// Default Firebase credentials path if not set and no embedded credentials
	if cfg.FirebaseCredentialsPath == "" && len(cfg.firebaseCredentialsJSON) == 0 {
		cfg.FirebaseCredentialsPath = DefaultFirebaseCredentialsPath()
	}

	return cfg, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath := DefaultConfigPath()

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// HasFirebaseCredentials checks if Firebase credentials are available
func (c *Config) HasFirebaseCredentials() bool {
	// Check embedded credentials first
	if len(c.firebaseCredentialsJSON) > 0 {
		return true
	}
	// Check file path
	if c.FirebaseCredentialsPath == "" {
		return false
	}
	_, err := os.Stat(c.FirebaseCredentialsPath)
	return err == nil
}

// GetFirebaseCredentialsJSON returns the Firebase credentials as JSON bytes
// Returns embedded credentials if available, otherwise reads from file
func (c *Config) GetFirebaseCredentialsJSON() ([]byte, error) {
	// Return embedded credentials if available
	if len(c.firebaseCredentialsJSON) > 0 {
		return c.firebaseCredentialsJSON, nil
	}
	// Read from file
	return os.ReadFile(c.FirebaseCredentialsPath)
}

// HasGoogleClientID checks if Google Client ID is configured
func (c *Config) HasGoogleClientID() bool {
	return c.GoogleClientID != "" && c.GoogleClientID != "YOUR_GOOGLE_CLIENT_ID"
}

// HasEmbeddedCredentials returns true if build-time credentials are embedded
func HasEmbeddedCredentials() bool {
	return EmbeddedGoogleClientID != "" || EmbeddedFirebaseCredentials != ""
}
