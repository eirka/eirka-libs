package user

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eirka/eirka-libs/config"
	e "github.com/eirka/eirka-libs/errors"
	jwt "github.com/golang-jwt/jwt/v5"
)

func init() {
	// Enable test mode for secret validation
	SetTestMode(true)
}

// resetJwtTestConfig resets the config for JWT tests
func resetJwtTestConfig() {
	// Reset secrets
	config.Settings.Session.NewSecret = ""
	config.Settings.Session.OldSecret = ""
}

func TestMakeToken(t *testing.T) {
	// Reset config
	resetJwtTestConfig()

	// No secret set
	token, err := MakeToken(2)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrNoSecret, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	// Set a valid secret
	config.Settings.Session.NewSecret = "secret"

	// default user state should never get a token
	token, err = MakeToken(0)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	// a non authed user should never get a token
	token, err = MakeToken(1)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	token, err = MakeToken(2)
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "Token should not be empty")
	}
}

func TestMakeTokenValidateOutput(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	token, err := MakeToken(2)
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "Token should not be empty")
	}

	out, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, out, "Token should not be empty")
	}

	// get the claims from the token
	claims, ok := out.Claims.(*TokenClaims)
	assert.True(t, ok, "Should be true")

	assert.Equal(t, claims.User, uint(2), "Claim should match")
}

func TestCreateTokenAnonAuth(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	invalidUser := DefaultUser()
	invalidUser.SetID(1)
	invalidUser.SetAuthenticated()

	notoken, err := invalidUser.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, notoken, "token should not be returned")
	}
}

func TestCreateTokenZeroAuth(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	invalidUser := DefaultUser()
	invalidUser.SetID(0)
	invalidUser.SetAuthenticated()

	notoken, err := invalidUser.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, notoken, "token should not be returned")
	}
}

func TestCreateTokenZeroNoAuth(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	invalidUser := DefaultUser()
	invalidUser.SetID(0)

	notoken, err := invalidUser.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, notoken, "token should not be returned")
	}
}

func TestCreateTokenBadPassword(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	// a user that doesnt have a validated password should never get a token
	token, err := user.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrInvalidPassword, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	user.hash, err = HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, user.hash, "password should be returned")
	}

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

	token, err = user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "Token should not be empty")
	}
}

// TestAlgorithmConfusionAttack verifies protection against algorithm confusion attacks
func TestAlgorithmConfusionAttack(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()
	user.SetID(2)

	// Create token with none algorithm (which should be rejected)
	claims := TokenClaims{
		2,
		jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * jwtExpireDays)),
		},
	}

	// Create the token with "none" algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	token.Header[jwtHeaderKeyID] = 1

	// Sign with none method
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if assert.NoError(t, err, "An error was not expected during token creation") {
		// Try to parse this token
		_, err = jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return validateToken(token, &user)
		})
		// Should get an error about invalid signing method
		assert.Error(t, err, "Algorithm confusion attack should be detected")
		assert.Contains(t, err.Error(), "unexpected signing method", "Error should indicate wrong signing method")
	}
}

// TestTokenTampering tests that modified tokens are rejected
func TestTokenTampering(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	// Create a valid token
	validToken, err := MakeToken(2)
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, validToken, "Token should not be empty")
	}

	// Attempt to tamper with the token by changing a character in the signature
	parts := strings.Split(validToken, ".")
	if len(parts) == 3 {
		// Tamper with the signature part
		if len(parts[2]) > 0 {
			// Change the first character of the signature
			if parts[2][0] == 'a' {
				parts[2] = "b" + parts[2][1:]
			} else {
				parts[2] = "a" + parts[2][1:]
			}
		}

		tamperedToken := parts[0] + "." + parts[1] + "." + parts[2]

		// Try to validate the tampered token
		user := DefaultUser()
		_, err = jwt.ParseWithClaims(tamperedToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return validateToken(token, &user)
		})

		// Should get an error about invalid signature
		assert.Error(t, err, "Tampered token should be rejected")
		assert.Contains(t, err.Error(), "signature is invalid", "Error should indicate invalid signature")
	}
}

// TestSensitiveDataInClaims ensures no sensitive data is included in token
func TestSensitiveDataInClaims(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	// Create a valid token
	validToken, err := MakeToken(2)
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, validToken, "Token should not be empty")
	}

	// Get the primary secret for verification
	primarySecret, err := GetPrimarySecret()
	assert.NoError(t, err)

	// Parse the token without verification to examine claims
	token, _ := jwt.ParseWithClaims(validToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(primarySecret), nil
	})

	// Get the claims
	claims, ok := token.Claims.(*TokenClaims)
	assert.True(t, ok, "Should be able to parse claims")

	// Verify only expected fields are present in the token
	assert.Equal(t, uint(2), claims.User, "User ID should match")
	assert.Equal(t, jwtIssuer, claims.Issuer, "Issuer should match")
	assert.NotNil(t, claims.IssuedAt, "IssuedAt should be present")
	assert.NotNil(t, claims.NotBefore, "NotBefore should be present")
	assert.NotNil(t, claims.ExpiresAt, "ExpiresAt should be present")

	// Check for potential sensitive data that should not be in the token
	tokStr := validToken
	sensitiveTerms := []string{
		"password", "secret", "hash", "email", "address", "credit", "ssn", "social",
	}

	for _, term := range sensitiveTerms {
		assert.NotContains(t, strings.ToLower(tokStr), strings.ToLower(term),
			"Token should not contain sensitive data: "+term)
	}
}

// TestCreateTokenNoSecret ensures proper error when no secret is set
func TestCreateTokenNoSecret(t *testing.T) {
	// Reset config
	resetJwtTestConfig()

	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()
	user.isPasswordValid = true

	token, err := user.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, e.ErrNoSecret, err, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}
}

// TestTokenTimeSkew tests token validity with time skew
func TestTokenTimeSkew(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()
	user.SetID(2)

	// Create a token with future time (simulating clock skew)
	futureTime := time.Now().Add(time.Hour)

	claims := TokenClaims{
		2,
		jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(futureTime),
			NotBefore: jwt.NewNumericDate(futureTime),
			ExpiresAt: jwt.NewNumericDate(futureTime.Add(time.Hour * 24 * jwtExpireDays)),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header[jwtHeaderKeyID] = 1

	primarySecret, err := GetPrimarySecret()
	assert.NoError(t, err)

	futureToken, err := token.SignedString([]byte(primarySecret))
	if assert.NoError(t, err, "An error was not expected") {
		// Try to parse this token (should fail because NotBefore is in future)
		_, err = jwt.ParseWithClaims(futureToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return validateToken(token, &user)
		})

		// Should get an error about token not being valid yet
		assert.Error(t, err, "Future token should be rejected")
		assert.Contains(t, err.Error(), "not valid yet", "Error should indicate token is not valid yet")
	}
}

// TestTokenWithUnsupportedAlgorithm tests token validation with various unsupported algorithms
func TestTokenWithUnsupportedAlgorithm(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()
	user.SetID(2)

	// We'll manually create a token with alg of "HS256" but then manually change header to test algorithm confusion
	tokenValid := jwt.New(jwt.SigningMethodHS256)
	claims := &TokenClaims{
		2,
		jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * jwtExpireDays)),
		},
	}
	tokenValid.Claims = claims
	tokenValid.Header[jwtHeaderKeyID] = 1

	primarySecret, err := GetPrimarySecret()
	assert.NoError(t, err)

	// Sign a valid token first, to have a good reference
	validTokenString, _ := tokenValid.SignedString([]byte(primarySecret))

	// Now let's create a new token with a different alg header
	parts := strings.Split(validTokenString, ".")
	if len(parts) != 3 {
		t.Fatal("Invalid token format for testing")
	}

	// Decode the header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("Error decoding header: %v", err)
	}

	// Parse header
	var header map[string]interface{}
	if err = json.Unmarshal(headerBytes, &header); err != nil {
		t.Fatalf("Error unmarshaling header: %v", err)
	}

	// Change algorithm to RS256
	header["alg"] = "RS256"

	// Encode header back
	modifiedHeaderBytes, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("Error marshaling header: %v", err)
	}

	// Create new token with modified header
	modifiedHeader := base64.RawURLEncoding.EncodeToString(modifiedHeaderBytes)
	modifiedTokenString := modifiedHeader + "." + parts[1] + "." + parts[2]

	// Try to validate the token with algorithm confusion
	_, err = jwt.ParseWithClaims(modifiedTokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, &user)
	})

	// Should be rejected with an error about unexpected signing method
	assert.Error(t, err, "Token with algorithm confusion should be rejected")
	assert.Contains(t, err.Error(), "unexpected signing method", "Error should indicate unexpected signing method")
}

// TestInvalidTokenClaims tests validation with tampered or invalid claims
func TestInvalidTokenClaims(t *testing.T) {
	// Reset and set up config
	resetJwtTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()
	user.SetID(2)

	now := time.Now()

	// Test just a few key cases that we know will fail validation
	testCases := []struct {
		name   string
		claims jwt.Claims
		errMsg string
	}{
		{
			name: "Invalid Issuer",
			claims: TokenClaims{
				2,
				jwt.RegisteredClaims{
					Issuer:    "invalid-issuer",
					IssuedAt:  jwt.NewNumericDate(now),
					NotBefore: jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
				},
			},
			errMsg: "incorrect issuer",
		},
		{
			name: "Zero User ID",
			claims: TokenClaims{
				0,
				jwt.RegisteredClaims{
					Issuer:    jwtIssuer,
					IssuedAt:  jwt.NewNumericDate(now),
					NotBefore: jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
				},
			},
			errMsg: "invalid user id",
		},
		{
			name: "Expired Token",
			claims: TokenClaims{
				2,
				jwt.RegisteredClaims{
					Issuer:    jwtIssuer,
					IssuedAt:  jwt.NewNumericDate(now.AddDate(0, 0, -100)),
					NotBefore: jwt.NewNumericDate(now.AddDate(0, 0, -100)),
					ExpiresAt: jwt.NewNumericDate(now.AddDate(0, 0, -1)), // Expired yesterday
				},
			},
			errMsg: "token is expired",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create token with test claims
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tc.claims)
			token.Header[jwtHeaderKeyID] = 1

			primarySecret, err := GetPrimarySecret()
			assert.NoError(t, err)

			tokenString, err := token.SignedString([]byte(primarySecret))
			if err != nil {
				t.Fatalf("Failed to sign token: %v", err)
			}

			// Try to validate
			_, err = jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return validateToken(token, &user)
			})

			// We expect all these test cases to fail validation
			assert.Error(t, err, "Token with %s should be rejected", tc.name)
			assert.Contains(t, err.Error(), tc.errMsg,
				"Error for %s should contain '%s'", tc.name, tc.errMsg)
		})
	}
}
