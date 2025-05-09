package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
	minAcceptable := int(expectedAvg * 0.5)  // Allow 50% deviation
	maxAcceptable := int(expectedAvg * 1.5)  // Allow 50% deviation
	
	for char, count := range charCount {
		assert.True(t, count >= minAcceptable && count <= maxAcceptable, 
			"Character %c distribution is outside acceptable range: %d occurrences", char, count)
	}
}
