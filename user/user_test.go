package user

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/eirka/eirka-libs/config"
	//"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

func init() {
	config.Settings.Limits.PasswordMinLength = 8
	config.Settings.Limits.PasswordMaxLength = 128
}

func TestDefaultUser(t *testing.T) {

	user := DefaultUser()

	assert.Equal(t, uint(1), user.Id, "default user id should be 1")

	assert.False(t, user.IsAuthenticated, "default user should not be authenticated")

}

func TestSetId(t *testing.T) {

	user := DefaultUser()

	user.SetId(2)

	assert.Equal(t, uint(2), user.Id, "user id should be 2")

}

func TestSetAuthenticated(t *testing.T) {

	user := DefaultUser()

	user.SetAuthenticated()

	assert.False(t, user.IsAuthenticated, "User should be not authorized")

	user.SetId(2)

	user.SetAuthenticated()

	assert.True(t, user.IsAuthenticated, "User should be authorized")

	assert.True(t, user.IsValid(), "Authed non-anon user should be valid")

}

func TestIsValid(t *testing.T) {

	user := DefaultUser()

	assert.True(t, user.IsValid(), "DefaultUser should be valid")

	user.SetId(2)

	assert.False(t, user.IsValid(), "Unauthenticated non-anon should be invalid")

	user.SetAuthenticated()

	assert.True(t, user.IsValid(), "Authed non-anon user should be valid")

	user.SetId(0)

	assert.False(t, user.IsValid(), "User zero should be invalid")

	user.SetId(1)

	assert.False(t, user.IsValid(), "An authenticated anon user should be invalid")

}

func TestIsValidName(t *testing.T) {

	assert.True(t, IsValidName("cooldude2"), "Name should validate")

	assert.True(t, IsValidName("cool dude"), "Name should validate")

	assert.True(t, IsValidName("cool_dude"), "Name should validate")

	assert.True(t, IsValidName("cool-dude"), "Name should validate")

	assert.False(t, IsValidName("cool.dude"), "Name should not validate")

	assert.False(t, IsValidName("cooldude!"), "Name should not validate")

	assert.False(t, IsValidName("admin"), "Name should not validate")

	assert.False(t, IsValidName("Admin"), "Name should not validate")

	assert.False(t, IsValidName("Admin "), "Name should not validate")

	assert.False(t, IsValidName(" Admin  "), "Name should not validate")

}

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
	user.SetId(2)
	user.SetAuthenticated()

	assert.False(t, user.ComparePassword("atestpassword"), "Password should not validate")

	user.hash = password

	assert.False(t, user.ComparePassword(""), "Password should not validate")

	assert.False(t, user.ComparePassword("wrongpassword"), "Password should not validate")

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

}

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

	user.SetId(2)

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
	invalidUser.SetId(1)
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
	invalidUser.SetId(0)
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
	invalidUser.SetId(0)

	notoken, err := invalidUser.CreateToken()
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
		assert.Empty(t, notoken, "token should not be returned")
	}
}
