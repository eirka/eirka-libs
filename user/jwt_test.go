package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	e "github.com/eirka/eirka-libs/errors"
)

func TestCreateToken(t *testing.T) {

	Secret = ""

	user := DefaultUser()

	// secret must be set
	token, err := user.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrNoSecret, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	Secret = "secret"

	// default user state should never get a token
	token, err = user.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	user.SetID(2)

	// a non authed user should never get a token
	token, err = user.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	user.SetAuthenticated()

	// a user that doesnt have a validated password should never get a token
	token, err = user.CreateToken()
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

func TestCreateTokenAnonAuth(t *testing.T) {

	Secret = "secret"

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

	Secret = "secret"

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

	Secret = "secret"

	invalidUser := DefaultUser()
	invalidUser.SetID(0)

	notoken, err := invalidUser.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, notoken, "token should not be returned")
	}
}
