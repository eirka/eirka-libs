package user

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	// Helper function for min value
	"golang.org/x/exp/constraints"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

func TestPassword(t *testing.T) {

	_, err := HashPassword("")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrPasswordEmpty, "Error should match")
	}

	_, err = HashPassword("heh")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrPasswordShort, "Error should match")
	}

	_, err = HashPassword("k8dyuCqJfW9v5iFUeeS4YOeiuk5Wee6Q9tZvWFHqE10ftzhaxVzxlKzx4n7CcBpRcgtaX9dZ2lBIRrsvgqXPPvmjNpIgnrums2Xtst8FsZkpZo61u3ChCs7MEO1DGy4Qa")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrPasswordLong, "Error should match")
	}

	password, err := HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, password, "password should be returned")
	}

	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	assert.False(t, user.ComparePassword("atestpassword"), "Password should not validate")

	user.hash = password

	assert.False(t, user.ComparePassword(""), "Password should not validate")

	assert.False(t, user.ComparePassword("wrongpassword"), "Password should not validate")

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

}

func TestGeneratePassword(t *testing.T) {
	// Test that the password is the correct length
	password := generateRandomPassword(20)
	assert.True(t, len(password) == 20, "Password should be 20 chars")

	// Test that we get different passwords each time (entropy test)
	for i := 1; i <= 1000; i++ {
		password1 := generateRandomPassword(20)
		password2 := generateRandomPassword(20)
		assert.NotEqual(t, password1, password2, "Passwords should not be equal")
	}

	// Check character set compliance
	for i := 0; i < 100; i++ {
		pass := generateRandomPassword(100) // Generate a long password to increase chances of coverage
		for _, char := range pass {
			assert.Contains(t, letterBytes, string(char), "Password contains invalid character")
		}
	}

	// Test with different lengths
	lengths := []int{1, 5, 10, 32, 64, 128}
	for _, length := range lengths {
		pass := generateRandomPassword(length)
		assert.Equal(t, length, len(pass), "Password length doesn't match requested length")
	}
}

func TestRandomPassword(t *testing.T) {
	for i := 1; i <= 10; i++ {
		password, hash, err := RandomPassword()
		if assert.NoError(t, err, "An error was not expected") {
			assert.NotNil(t, hash, "hash should be returned")
			assert.NotEmpty(t, password, "password should be returned")
		}

		user := DefaultUser()
		user.SetID(2)
		user.SetAuthenticated()
		user.hash = hash

		assert.True(t, user.ComparePassword(password), "Password should validate")
	}
}

// TestGeneratePasswordZeroLength tests the behavior when generating a password of zero length
func TestGeneratePasswordZeroLength(t *testing.T) {
	// Should return an empty string but not crash
	pass := generateRandomPassword(0)
	assert.Equal(t, "", pass, "Zero-length password should be empty string")
}

// Helper for min
func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// TestPasswordDistribution tests that our password generator creates a relatively even distribution
// of characters from the letterBytes set
func TestPasswordDistribution(t *testing.T) {
	// This is a statistical test to check if our random generation is reasonably distributed
	// We're not looking for perfect distribution, but want to ensure there's no obvious bias
	charCount := make(map[rune]int)

	// Initialize the map with all possible characters
	for _, char := range letterBytes {
		charCount[char] = 0
	}

	// Generate a very long password to get a good sample
	longPass := generateRandomPassword(10000)

	// Count occurrences of each character
	for _, char := range longPass {
		charCount[char]++
	}

	// With 10000 characters and 62 possible values, each character should
	// appear approximately 161 times (10000/62). We'll allow a reasonable deviation.
	expectedAvg := 10000.0 / float64(len(letterBytes))
	minAcceptable := int(expectedAvg * 0.5) // Allow 50% deviation
	maxAcceptable := int(expectedAvg * 1.5) // Allow 50% deviation

	for char, count := range charCount {
		assert.True(t, count >= minAcceptable && count <= maxAcceptable,
			"Character %c distribution is outside acceptable range: %d occurrences", char, count)
	}
}

// TestComparePasswordWithEmptyHash verifies behavior when the user hash is empty
func TestComparePasswordWithEmptyHash(t *testing.T) {
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	// Explicitly set hash to empty
	user.hash = []byte{}

	// This should return false since hash is empty
	assert.False(t, user.ComparePassword("testpassword"), "Password comparison with empty hash should fail")
	assert.False(t, user.isPasswordValid, "isPasswordValid should be false with empty hash")
}

// TestBcryptCost verifies that we're using a sufficiently strong bcrypt cost
func TestBcryptCost(t *testing.T) {
	// Best practice recommends at least cost 10 for security
	// DefaultCost is 10, but we should verify that's what we're using

	password, err := HashPassword("testingcost")
	assert.NoError(t, err, "An error was not expected")

	// Extract the cost parameter from the hash
	cost, err := bcrypt.Cost(password)
	assert.NoError(t, err, "An error was not expected")

	// Verify cost is at least 10
	assert.GreaterOrEqual(t, cost, 10, "Bcrypt cost should be at least 10 for security")
}

// TestPasswordComplexity tests how we handle password complexity
// Note: This is a test to verify what the current behavior is.
// If a password complexity policy is desired, the code would need to be updated.
func TestPasswordComplexity(t *testing.T) {
	// Currently, the code only checks length, not complexity
	// These simple passwords should pass since they meet the length requirement
	validPasswords := []string{
		"password123",
		"123456789012",
		strings.Repeat("a", config.Settings.Limits.PasswordMinLength),
	}

	for _, pass := range validPasswords {
		hash, err := HashPassword(pass)
		assert.NoError(t, err, "Password should be accepted: "+pass)
		assert.NotEmpty(t, hash, "Hash should be generated")
	}
}

// TestRandomPasswordSecurity verifies security properties of random passwords
func TestRandomPasswordSecurity(t *testing.T) {
	// Check that RandomPassword generates passwords with sufficient entropy

	// Generate multiple random passwords and check they're all different
	passwords := make(map[string]bool)
	iterations := 100

	for i := 0; i < iterations; i++ {
		pass, hash, err := RandomPassword()
		assert.NoError(t, err, "RandomPassword should not error")
		assert.NotEmpty(t, pass, "Password should not be empty")
		assert.NotEmpty(t, hash, "Hash should not be empty")

		// Check for duplicates
		assert.False(t, passwords[pass], "Generated password should be unique")
		passwords[pass] = true

		// Verify the password has expected length (20 characters)
		assert.Equal(t, 20, len(pass), "Random password should be 20 characters")
	}
}

// TestBcryptCompareWithEmptyPassword verifies behavior of ComparePassword with empty password
func TestBcryptCompareWithEmptyPassword(t *testing.T) {
	// Create a hash for a valid password
	validPassword := "testpassword"
	hash, err := HashPassword(validPassword)
	assert.NoError(t, err, "HashPassword should not error")

	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()
	user.hash = hash

	// This should return false
	assert.False(t, user.ComparePassword(""), "Empty password should not validate")
	assert.False(t, user.isPasswordValid, "isPasswordValid should be false after failed comparison")
}

// TestHashPasswordWithExtremeLength tests the behavior with very long passwords
func TestHashPasswordWithExtremeLength(t *testing.T) {
	// Bcrypt has an internal limit around 72 bytes
	// Let's test passwords right at the limits

	// Just below config max length but also below bcrypt's 72 byte limit
	// Note: bcrypt has a 72 byte limit, so we need to ensure our test doesn't exceed that
	maxLength := min(71, config.Settings.Limits.PasswordMaxLength-1)
	almostTooLong := strings.Repeat("A", maxLength)
	hash, err := HashPassword(almostTooLong)
	assert.NoError(t, err, "Password just under limit should be accepted")
	assert.NotEmpty(t, hash, "Hash should be generated")

	// Over config max length
	tooLong := strings.Repeat("A", config.Settings.Limits.PasswordMaxLength+1)
	hash, err = HashPassword(tooLong)
	assert.Equal(t, e.ErrPasswordLong, err, "Password over limit should return ErrPasswordLong")
	assert.Empty(t, hash, "Hash should not be generated for too-long password")
}

// TestConstantTimeCompare ensures that password comparison doesn't leak timing information
func TestConstantTimeCompare(t *testing.T) {
	// Create a user with a password
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	password := "correctPassword123"
	hash, err := HashPassword(password)
	assert.NoError(t, err, "Hashing password should not error")
	user.hash = hash

	// Run timing tests
	iterations := 100

	// Valid comparison
	validStart := time.Now()
	for i := 0; i < iterations; i++ {
		user.ComparePassword(password)
	}
	validDuration := time.Since(validStart)
	validAvg := validDuration.Nanoseconds() / int64(iterations)

	// Invalid comparison with same length
	wrongPassword := "incorrectPassword1"
	invalidStart := time.Now()
	for i := 0; i < iterations; i++ {
		user.ComparePassword(wrongPassword)
	}
	invalidDuration := time.Since(invalidStart)
	invalidAvg := invalidDuration.Nanoseconds() / int64(iterations)

	// The timing shouldn't be too different for valid vs invalid
	// This is a loose test since exact timing depends on system load
	// We're mainly checking that we use bcrypt which is designed to be constant-time
	ratio := float64(validAvg) / float64(invalidAvg)
	assert.True(t, ratio > 0.5 && ratio < 2.0,
		"Valid and invalid password comparison times should be within reasonable range")

	// Most importantly, checking an empty password should be rejected quickly
	assert.False(t, user.ComparePassword(""), "Empty password should be rejected")
}

// TestSecureRandomness tests that the random password generator uses a secure source of randomness
func TestSecureRandomness(t *testing.T) {
	// Generate multiple passwords and ensure they're unique and unpredictable
	passwordCount := 1000
	passwordLength := 20
	passwords := make(map[string]bool, passwordCount)

	// Generate a large set of passwords
	for i := 0; i < passwordCount; i++ {
		password := generateRandomPassword(passwordLength)
		assert.Len(t, password, passwordLength, "Password should be correct length")

		// Should never generate duplicates with secure randomness
		assert.False(t, passwords[password], "Password should be unique: %s", password)
		passwords[password] = true
	}

	// Test character distribution - all character classes should be represented
	// For crypto-secure RNG, we should see all character classes well-represented
	alphabet := make(map[byte]int)
	for i := 0; i < len(letterBytes); i++ {
		alphabet[letterBytes[i]] = 0
	}

	// Count occurrences of each character across all passwords
	for pass := range passwords {
		for i := 0; i < len(pass); i++ {
			alphabet[pass[i]]++
		}
	}

	// All characters should be used
	for i := 0; i < len(letterBytes); i++ {
		count := alphabet[letterBytes[i]]
		assert.True(t, count > 0, "Character %c should be used at least once", letterBytes[i])
	}
}

// TestUpdatePasswordSecurity tests various security aspects of updating passwords
func TestUpdatePasswordSecurity(t *testing.T) {
	// Test with a minimal valid password
	minPassword := strings.Repeat("a", config.Settings.Limits.PasswordMinLength)
	hash, err := HashPassword(minPassword)
	assert.NoError(t, err, "Should accept minimum length password")

	// Test with various invalid UIDs
	invalidUserIDs := []uint{0, 1}
	for _, uid := range invalidUserIDs {
		err = UpdatePassword(hash, uid)
		assert.Equal(t, e.ErrUserNotValid, err, "Should reject invalid user IDs: %d", uid)
	}

	// Test with NULL hash
	err = UpdatePassword(nil, 2)
	assert.Equal(t, e.ErrInvalidPassword, err, "Should reject NULL hash")

	// Test with empty hash
	err = UpdatePassword([]byte{}, 2)
	assert.Equal(t, e.ErrInvalidPassword, err, "Should reject empty hash")

	// Create a mock DB to simulate a database failure
	mock, err := db.NewTestDb()
	assert.NoError(t, err, "Creating test DB should not error")

	// Mock a database error
	mock.ExpectExec("UPDATE users SET user_password").
		WithArgs(hash, 2).
		WillReturnError(e.ErrInternalError)

	// Should handle DB errors
	err = UpdatePassword(hash, 2)
	assert.Error(t, err, "Should handle database errors")

	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")
}
