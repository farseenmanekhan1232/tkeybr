package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// App wraps Firebase services
type App struct {
	app       *firebase.App
	auth      *auth.Client
	firestore *firestore.Client
}

// NewApp creates a new Firebase app with credentials from a file path
func NewApp(ctx context.Context, credentialsPath string) (*App, error) {
	opt := option.WithCredentialsFile(credentialsPath)
	return newAppWithOption(ctx, opt)
}

// NewAppWithJSON creates a new Firebase app with credentials from JSON bytes
// This is used for embedded credentials from build-time injection
func NewAppWithJSON(ctx context.Context, credentialsJSON []byte) (*App, error) {
	opt := option.WithCredentialsJSON(credentialsJSON)
	return newAppWithOption(ctx, opt)
}

// newAppWithOption creates a Firebase app with the given option
func newAppWithOption(ctx context.Context, opt option.ClientOption) (*App, error) {
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase Auth: %w", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firestore: %w", err)
	}

	return &App{
		app:       app,
		auth:      authClient,
		firestore: firestoreClient,
	}, nil
}

// Auth returns the Firebase Auth client
func (a *App) Auth() *auth.Client {
	return a.auth
}

// Firestore returns the Firestore client
func (a *App) Firestore() *firestore.Client {
	return a.firestore
}

// Close closes the Firebase app and releases resources
func (a *App) Close() error {
	if a.firestore != nil {
		return a.firestore.Close()
	}
	return nil
}

// VerifyIDToken verifies a Firebase ID token and returns the user's UID
func (a *App) VerifyIDToken(ctx context.Context, idToken string) (string, error) {
	token, err := a.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return "", fmt.Errorf("failed to verify ID token: %w", err)
	}
	return token.UID, nil
}

// GetOrCreateUser gets or creates a Firebase user from Google OAuth info
func (a *App) GetOrCreateUser(ctx context.Context, email, displayName, photoURL string) (*auth.UserRecord, error) {
	// Try to get existing user by email
	user, err := a.auth.GetUserByEmail(ctx, email)
	if err == nil {
		return user, nil
	}

	// Create new user if not found
	params := (&auth.UserToCreate{}).
		Email(email).
		DisplayName(displayName).
		PhotoURL(photoURL).
		EmailVerified(true)

	user, err = a.auth.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// CreateCustomToken creates a custom token for a user (for client-side auth)
func (a *App) CreateCustomToken(ctx context.Context, uid string) (string, error) {
	token, err := a.auth.CustomToken(ctx, uid)
	if err != nil {
		return "", fmt.Errorf("failed to create custom token: %w", err)
	}
	return token, nil
}
