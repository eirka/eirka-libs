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

	password := generateRandomPassword(20)

	assert.True(t, len(password) == 20, "Password should be 20 chars")

	for i := 1; i <= 1000; i++ {
		password1 := generateRandomPassword(20)
		password2 := generateRandomPassword(20)
		assert.NotEqual(t, password1, password2, "Passwords should not be equal")
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
