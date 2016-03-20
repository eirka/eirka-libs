package user

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

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

	badnames := []string{
		"cool.dude",
		"cooldude!",
		"admin",
		"Admin",
		"Admin  ",
		"  Admin  ",
		"Mod",
	}

	for _, name := range badnames {
		assert.False(t, IsValidName(name), "Name should not validate")
	}

	goodnames := []string{
		"cooldude2",
		"cooldude",
		"cool dude",
		"cool-dude",
		"way cool dude",
		"way-cool-dude69",
	}

	for _, name := range goodnames {
		assert.True(t, IsValidName(name), "Name should validate")
	}

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

func TestUserPassword(t *testing.T) {

	Secret = "secret"

	user := DefaultUser()
	user.SetId(2)
	user.SetAuthenticated()

	password, err := HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, password, "password should be returned")
	}

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"name", "password"}).AddRow("testaccount", password)

	mock.ExpectQuery("select user_name, user_password from users where user_id").WillReturnRows(rows)

	err = user.Password()
	if assert.NoError(t, err, "An error was not expected") {
		assert.Equal(t, user.Name, "testaccount", "Name should match")
		assert.True(t, user.ComparePassword("testpassword"), "Password should validate")
	}

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestUserBadPassword(t *testing.T) {

	Secret = "secret"

	user := DefaultUser()
	user.SetId(2)
	user.SetAuthenticated()

	password, err := HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, password, "password should be returned")
	}

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"name", "password"}).AddRow("testaccount", password)

	mock.ExpectQuery("select user_name, user_password from users where user_id").WillReturnRows(rows)

	err = user.Password()
	if assert.NoError(t, err, "An error was not expected") {
		assert.Equal(t, user.Name, "testaccount", "Name should match")
		assert.False(t, user.ComparePassword("badpassword"), "Password should not validate")
	}

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestFromName(t *testing.T) {

	Secret = "secret"

	user := DefaultUser()

	password, err := HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, password, "password should be returned")
	}

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"id", "password"}).AddRow(2, password)

	mock.ExpectQuery("select user_id, user_password from users where user_name").WillReturnRows(rows)

	err = user.FromName("testaccount")
	if assert.NoError(t, err, "An error was not expected") {
		assert.Equal(t, user.Id, uint(2), "Id should match")
		assert.True(t, user.IsAuthenticated, "User should be authenticated")
		assert.True(t, user.ComparePassword("testpassword"), "Password should validate")
	}

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestFromNameEmptyName(t *testing.T) {

	Secret = "secret"

	user := DefaultUser()

	err := user.FromName("")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
	}

}

func TestFromNameBadId(t *testing.T) {

	Secret = "secret"

	user := DefaultUser()

	password, err := HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, password, "password should be returned")
	}

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"id", "password"}).AddRow(0, password)

	mock.ExpectQuery("select user_id, user_password from users where user_name").WillReturnRows(rows)

	err = user.FromName("test")
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
	}

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestCheckDuplicateEmpty(t *testing.T) {

	assert.True(t, CheckDuplicate(""), "Should return true")

}

func TestCheckDuplicateGood(t *testing.T) {

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

	mock.ExpectQuery(`select count\(\*\) from users where user_name`).WillReturnRows(rows)

	assert.False(t, CheckDuplicate("test"), "Should not be a duplicate")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestCheckDuplicateBad(t *testing.T) {

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

	mock.ExpectQuery(`select count\(\*\) from users where user_name`).WillReturnRows(rows)

	assert.True(t, CheckDuplicate("test"), "Should be a duplicate")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestIsAuthorizedInvalid(t *testing.T) {

	user := DefaultUser()

	assert.False(t, user.IsAuthorized(0), "Should not be authorized")

	assert.False(t, user.IsAuthorized(1), "Should not be authorized")

	user.SetAuthenticated()

	assert.False(t, user.IsAuthorized(1), "Should not be authorized")

	user.SetId(2)

	assert.False(t, user.IsAuthorized(1), "Should not be authorized")

}

func TestIsAuthorizedDefault(t *testing.T) {

	user := DefaultUser()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"role"}).AddRow(1)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	assert.False(t, user.IsAuthorized(1), "Should not be authorized")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestIsAuthorizedAuth(t *testing.T) {

	user := DefaultUser()
	user.SetId(2)
	user.SetAuthenticated()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"role"}).AddRow(2)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	assert.False(t, user.IsAuthorized(1), "Should not be authorized")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestIsAuthorizedMod(t *testing.T) {

	user := DefaultUser()
	user.SetId(2)
	user.SetAuthenticated()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"role"}).AddRow(3)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	assert.True(t, user.IsAuthorized(1), "Should be authorized")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestIsAuthorizedAdmin(t *testing.T) {

	user := DefaultUser()
	user.SetId(2)
	user.SetAuthenticated()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"role"}).AddRow(4)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	assert.True(t, user.IsAuthorized(1), "Should be authorized")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}
