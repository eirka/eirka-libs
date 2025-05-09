package user

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eirka/eirka-libs/config"
	e "github.com/eirka/eirka-libs/errors"
)

// setupTestSecrets sets up the config secrets for testing
func setupTestSecrets() {
	// Make sure settings is initialized
	if config.Settings == nil {
		config.Settings = &config.Config{}
	}

	// Set test mode for secret validation
	SetTestMode(true)

	// Set secrets in config
	config.Settings.Session.NewSecret = "new-test-secret"
	config.Settings.Session.OldSecret = "old-test-secret"
}

// setupTestSecretsWithRotation sets up the config secrets for testing rotation
func setupTestSecretsWithRotation() {
	// Make sure settings is initialized
	if config.Settings == nil {
		config.Settings = &config.Config{}
	}

	// Set test mode for secret validation
	SetTestMode(true)

	// Set secrets in config for rotation
	config.Settings.Session.NewSecret = "new-test-secret"
	config.Settings.Session.OldSecret = "old-test-secret"
}

// setupTestSecretsNoRotation sets up the config with only a new secret
func setupTestSecretsNoRotation() {
	// Make sure settings is initialized
	if config.Settings == nil {
		config.Settings = &config.Config{}
	}

	// Set test mode for secret validation
	SetTestMode(true)

	// Set only new secret
	config.Settings.Session.NewSecret = "new-test-secret"
	config.Settings.Session.OldSecret = ""
}

// resetTestSecrets resets the config secrets for testing
func resetTestSecrets() {
	// Make sure settings is initialized
	if config.Settings == nil {
		config.Settings = &config.Config{}
	}

	// Reset secrets in config
	config.Settings.Session.NewSecret = ""
	config.Settings.Session.OldSecret = ""
}

func init() {
	// Enable test mode for secret validation in tests
	SetTestMode(true)

	// Create a Session struct in config.Settings if it doesn't exist
	if config.Settings == nil {
		config.Settings = &config.Config{}
	}

	// Initialize session if needed
	config.Settings.Session = config.Session{
		NewSecret: "",
		OldSecret: "",
	}
}

func TestSecretValidation(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid secret",
			secret:  "a-valid-secret-that-is-long-enough",
			wantErr: false,
		},
		{
			name:    "Empty secret",
			secret:  "",
			wantErr: true,
			errMsg:  "secret cannot be empty",
		},
		{
			name:    "Secret too short",
			secret:  "toosh",
			wantErr: true,
			errMsg:  fmt.Sprintf("secret must be at least %d characters", MinSecretLengthTest),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSecret(tt.secret)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetPrimarySecret(t *testing.T) {
	// Initialize with session struct
	config.Settings.Session = config.Session{
		NewSecret: "",
		OldSecret: "",
	}

	// With no secret set
	secret, err := GetPrimarySecret()
	assert.Error(t, err)
	assert.Equal(t, e.ErrNoSecret, err)
	assert.Empty(t, secret)

	// Set up valid secret
	config.Settings.Session.NewSecret = "a-valid-secret-that-is-long-enough"

	// Test getting the primary secret
	secret, err = GetPrimarySecret()
	assert.NoError(t, err)
	assert.Equal(t, "a-valid-secret-that-is-long-enough", secret)

	// Explicitly enable test mode
	SetTestMode(true)

	// Test with invalid secret (too short, but not empty)
	config.Settings.Session.NewSecret = "short"
	secret, err = GetPrimarySecret()
	assert.Error(t, err)
	// In test mode, minimum length should be 6
	assert.Contains(t, err.Error(), fmt.Sprintf("secret must be at least %d characters", MinSecretLengthTest))
	assert.Empty(t, secret)

	// Reset the config
	config.Settings.Session.NewSecret = ""
	config.Settings.Session.OldSecret = ""
}

func TestIsInitialized(t *testing.T) {
	// Initialize with session struct
	config.Settings.Session = config.Session{
		NewSecret: "",
		OldSecret: "",
	}

	// Not initialized
	assert.False(t, IsInitialized())

	// Set valid new secret
	config.Settings.Session.NewSecret = "a-valid-secret-that-is-long-enough"
	assert.True(t, IsInitialized())

	// Invalid secret
	config.Settings.Session.NewSecret = "short"
	assert.False(t, IsInitialized())

	// Empty secret
	config.Settings.Session.NewSecret = ""
	assert.False(t, IsInitialized())

	// Reset the config
	resetTestSecrets()
}

func TestIsRotationActive(t *testing.T) {
	// Initialize with session struct
	config.Settings.Session = config.Session{
		NewSecret: "",
		OldSecret: "",
	}

	// No secrets
	assert.False(t, IsRotationActive())

	// Only new secret
	config.Settings.Session.NewSecret = "a-valid-secret-that-is-long-enough"
	config.Settings.Session.OldSecret = ""
	assert.False(t, IsRotationActive())

	// Both old and new secrets
	config.Settings.Session.NewSecret = "a-valid-secret-that-is-long-enough"
	config.Settings.Session.OldSecret = "old-secret-that-is-long-enough"
	assert.True(t, IsRotationActive())

	// Only old secret
	config.Settings.Session.NewSecret = ""
	config.Settings.Session.OldSecret = "old-secret-that-is-long-enough"
	assert.False(t, IsRotationActive())

	// Reset the config
	resetTestSecrets()
}

func TestGetSecrets(t *testing.T) {
	// Initialize with session struct
	config.Settings.Session = config.Session{
		NewSecret: "",
		OldSecret: "",
	}

	// Before initialization
	secrets, err := GetSecrets()
	assert.Error(t, err)
	assert.Equal(t, e.ErrNoSecret, err)
	assert.Nil(t, secrets)

	// With only new secret
	config.Settings.Session.NewSecret = "new-secret-that-is-valid"
	config.Settings.Session.OldSecret = ""

	secrets, err = GetSecrets()
	assert.NoError(t, err)
	assert.Len(t, secrets, 1)
	assert.Equal(t, "new-secret-that-is-valid", secrets[0])

	// With both old and new secrets
	config.Settings.Session.NewSecret = "new-secret-that-is-valid"
	config.Settings.Session.OldSecret = "old-secret-that-is-valid"

	secrets, err = GetSecrets()
	assert.NoError(t, err)
	assert.Len(t, secrets, 2)
	assert.Equal(t, "new-secret-that-is-valid", secrets[0])
	assert.Equal(t, "old-secret-that-is-valid", secrets[1])

	// Invalid new secret
	config.Settings.Session.NewSecret = "short"
	config.Settings.Session.OldSecret = "old-secret-that-is-valid"

	secrets, err = GetSecrets()
	assert.Error(t, err)
	assert.Nil(t, secrets)

	// Reset the config
	resetTestSecrets()
}

func TestConcurrentSecretAccess(t *testing.T) {
	// Initialize with session struct
	config.Settings.Session = config.Session{
		NewSecret: "",
		OldSecret: "",
	}
	config.Settings.Session.NewSecret = "concurrent-test-secret"

	// Create several goroutines that read the secret concurrently
	const numGoroutines = 10
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				secret, err := GetPrimarySecret()
				assert.NoError(t, err)
				assert.Equal(t, "concurrent-test-secret", secret)
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// The secret should remain unchanged
	secret, err := GetPrimarySecret()
	assert.NoError(t, err)
	assert.Equal(t, "concurrent-test-secret", secret)

	// Reset the config
	resetTestSecrets()
}

func TestConfigBasedRotationScenario(t *testing.T) {
	// Initialize with session struct
	config.Settings.Session = config.Session{
		NewSecret: "",
		OldSecret: "",
	}

	// Step 1: Initialize with only new secret
	config.Settings.Session.NewSecret = "initial-secret-value"
	config.Settings.Session.OldSecret = ""

	// Verify single secret setup
	secrets, err := GetSecrets()
	assert.NoError(t, err)
	assert.Len(t, secrets, 1)
	assert.Equal(t, "initial-secret-value", secrets[0])
	assert.False(t, IsRotationActive())

	// Step 2: Perform rotation by updating config
	config.Settings.Session.OldSecret = "initial-secret-value"
	config.Settings.Session.NewSecret = "new-secret-value"

	// Verify rotation is active
	assert.True(t, IsRotationActive())

	// Verify both secrets are accessible
	secrets, err = GetSecrets()
	assert.NoError(t, err)
	assert.Len(t, secrets, 2)
	assert.Equal(t, "new-secret-value", secrets[0])
	assert.Equal(t, "initial-secret-value", secrets[1])

	// Step 3: Complete rotation by clearing old secret
	config.Settings.Session.OldSecret = ""

	// Verify only new secret remains
	secrets, err = GetSecrets()
	assert.NoError(t, err)
	assert.Len(t, secrets, 1)
	assert.Equal(t, "new-secret-value", secrets[0])
	assert.False(t, IsRotationActive())

	// Reset the config
	resetTestSecrets()
}
