package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

func TestDefaultUser(t *testing.T) {

	user := DefaultUser()

	assert.Equal(t, uint(1), user.ID, "default user id should be 1")

	assert.False(t, user.IsAuthenticated, "default user should not be authenticated")

}

func TestSetId(t *testing.T) {

	user := DefaultUser()

	user.SetID(2)

	assert.Equal(t, uint(2), user.ID, "user id should be 2")

}

func TestSetAuthenticated(t *testing.T) {

	user := DefaultUser()

	user.SetAuthenticated()

	assert.False(t, user.IsAuthenticated, "User should be not authorized")

	user.SetID(2)

	user.SetAuthenticated()

	assert.True(t, user.IsAuthenticated, "User should be authorized")

	assert.True(t, user.IsValid(), "Authed non-anon user should be valid")

}

func TestIsValid(t *testing.T) {

	user := DefaultUser()

	assert.True(t, user.IsValid(), "DefaultUser should be valid")

	user.SetID(2)

	assert.False(t, user.IsValid(), "Unauthenticated non-anon should be invalid")

	user.SetAuthenticated()

	assert.True(t, user.IsValid(), "Authed non-anon user should be valid")

	user.SetID(0)

	assert.False(t, user.IsValid(), "User zero should be invalid")

	user.SetID(1)

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

func TestUserPassword(t *testing.T) {

	Secret = "secret"

	user := DefaultUser()
	user.SetID(2)
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
	user.SetID(2)
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
		assert.Equal(t, user.ID, uint(2), "Id should match")
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

	user.SetID(2)

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
	user.SetID(2)
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
	user.SetID(2)
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
	user.SetID(2)
	user.SetAuthenticated()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"role"}).AddRow(4)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	assert.True(t, user.IsAuthorized(1), "Should be authorized")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestUpdatePassword(t *testing.T) {

	var err error

	_, hash, err := RandomPassword()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, hash, "hash should be returned")
	}

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	mock.ExpectExec("UPDATE users SET user_password").
		WithArgs(hash, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdatePassword(hash, 2)
	assert.NoError(t, err, "An error was not expected")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

}

func TestUpdatePasswordBadHash(t *testing.T) {

	var err error

	err = UpdatePassword(nil, 2)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrInvalidPassword, "Error should match")
	}

	err = UpdatePassword([]byte{}, 2)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrInvalidPassword, "Error should match")
	}

}

func TestUpdatePasswordBadUser(t *testing.T) {

	var err error

	_, hash, err := RandomPassword()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, hash, "hash should be returned")
	}

	err = UpdatePassword(hash, 1)
	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err, e.ErrUserNotValid, "Error should match")
	}

}
