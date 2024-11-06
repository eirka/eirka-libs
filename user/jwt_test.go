package user

import (
	"testing"

	"github.com/stretchr/testify/assert"

	e "github.com/eirka/eirka-libs/errors"
	jwt "github.com/golang-jwt/jwt/v5"
)

func TestMakeToken(t *testing.T) {

	// secret must be set
	token, err := MakeToken("", 2)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrNoSecret, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	// default user state should never get a token
	token, err = MakeToken("secret", 0)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	// a non authed user should never get a token
	token, err = MakeToken("secret", 1)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, token, "Token should be empty")
	}

	token, err = MakeToken("secret", 2)
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "Token should not be empty")
	}

}

func TestMakeTokenValidateOutput(t *testing.T) {

	token, err := MakeToken("secret", 2)
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

func TestCreateTokenBadPassword(t *testing.T) {

	Secret = "secret"

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
