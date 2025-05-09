package user

import (
	"errors"
	"fmt"
	"sync"

	"github.com/eirka/eirka-libs/config"
	e "github.com/eirka/eirka-libs/errors"
)

const (
	// MinSecretLength defines the minimum allowed length for a secret in production
	MinSecretLength = 16
	// MinSecretLengthTest is a shorter length used for tests
	MinSecretLengthTest = 6
)

var (
	// errEmptySecret is returned when an empty secret is provided
	errEmptySecret = errors.New("secret cannot be empty")
	// isTestMode indicates whether we're in test mode with relaxed validation
	isTestMode = false
)

// SetTestMode enables or disables test mode for secret validation
func SetTestMode(enabled bool) {
	isTestMode = enabled
}

// SecretManager handles JWT signing secrets with support for rotation
// This is now a simple wrapper around the config.Settings.Session
// All secrets are directly read from config rather than stored internally
type SecretManager struct {
	mu sync.RWMutex
}

// secretManager is the singleton instance of SecretManager
var secretManager = &SecretManager{}

// GetPrimarySecret returns the primary (new) secret for signing new tokens
func GetPrimarySecret() (string, error) {
	secretManager.mu.RLock()
	defer secretManager.mu.RUnlock()

	// Ensure Settings is initialized
	if config.Settings == nil {
		return "", e.ErrNoSecret
	}

	newSecret := config.Settings.Session.NewSecret
	if newSecret == "" {
		return "", e.ErrNoSecret
	}

	// Validate the secret for security
	if err := validateSecret(newSecret); err != nil {
		return "", err
	}

	return newSecret, nil
}

// GetSecrets returns all active secrets for token validation
// The first return value is always the new secret (primary)
// The second return value is the old secret (if exists)
func GetSecrets() ([]string, error) {
	secretManager.mu.RLock()
	defer secretManager.mu.RUnlock()

	// Ensure Settings is initialized
	if config.Settings == nil {
		return nil, e.ErrNoSecret
	}

	newSecret := config.Settings.Session.NewSecret
	oldSecret := config.Settings.Session.OldSecret

	// Require at least a new secret
	if newSecret == "" {
		return nil, e.ErrNoSecret
	}

	// Validate new secret
	if err := validateSecret(newSecret); err != nil {
		return nil, err
	}

	// Create the slice of secrets
	secrets := []string{newSecret}

	// Only add old secret if it exists and passes validation
	if oldSecret != "" {
		if err := validateSecret(oldSecret); err == nil {
			secrets = append(secrets, oldSecret)
		}
	}

	return secrets, nil
}

// validateSecret checks if a secret meets the required criteria
func validateSecret(secret string) error {
	if secret == "" {
		return errEmptySecret
	}

	// Use relaxed validation in test mode
	minLength := MinSecretLength
	if isTestMode {
		minLength = MinSecretLengthTest
	}

	if len(secret) < minLength {
		// Create a fresh error each time with the appropriate length
		if isTestMode {
			return fmt.Errorf("secret must be at least %d characters", MinSecretLengthTest)
		}
		return fmt.Errorf("secret must be at least %d characters", MinSecretLength)
	}

	return nil
}

// IsInitialized returns true if the new secret is properly configured
func IsInitialized() bool {
	secretManager.mu.RLock()
	defer secretManager.mu.RUnlock()

	// Ensure Settings is initialized
	if config.Settings == nil {
		return false
	}

	newSecret := config.Settings.Session.NewSecret
	if newSecret == "" {
		return false
	}

	// Also check if it passes validation
	return validateSecret(newSecret) == nil
}

// IsRotationActive returns true if both old and new secrets are set
func IsRotationActive() bool {
	secretManager.mu.RLock()
	defer secretManager.mu.RUnlock()

	// Ensure Settings is initialized
	if config.Settings == nil {
		return false
	}

	return config.Settings.Session.OldSecret != "" && config.Settings.Session.NewSecret != ""
}
