package user

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
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

// TestAuthenticatorInterface verifies that the User struct implements the Authenticator interface
func TestAuthenticatorInterface(t *testing.T) {
	// This is a compile-time check that User implements Authenticator
	var _ Authenticator = &User{}

	// Create a user and call interface methods to verify they work
	user := DefaultUser()

	// Call methods from the interface
	assert.True(t, user.IsValid(), "Default user should be valid")
	assert.False(t, user.IsAuthorized(1), "Default user should not be authorized")

	// Methods that modify the user
	user.SetID(2)
	user.SetAuthenticated()

	// Verify the changes took effect
	assert.Equal(t, uint(2), user.ID, "User ID should be 2")
	assert.True(t, user.IsAuthenticated, "User should be authenticated")
}

// TestIsValidNameWithEmptyString tests validation of empty strings
func TestIsValidNameWithEmptyString(t *testing.T) {
	assert.False(t, IsValidName(""), "Empty string should not be a valid name")
}

// TestReservedNameCaseSensitivity tests if the reserved name check is case insensitive
func TestReservedNameCaseSensitivity(t *testing.T) {
	// Test with various cases of reserved names
	reservedVariations := []string{
		"Admin",
		"ADMIN",
		"admin",
		"aDmIn",
		"Moderator",
		"MOD",
		"mod",
		"Mod",
	}

	for _, name := range reservedVariations {
		assert.False(t, IsValidName(name), "Reserved name variation should not be valid: "+name)
	}
}

// TestPasswordDatabaseError tests error handling when the database query fails
func TestPasswordDatabaseError(t *testing.T) {
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected when setting up mock")

	// Mock a database error
	mock.ExpectQuery("select user_name, user_password from users where user_id").
		WillReturnError(errors.New("database error"))

	err = user.Password()
	assert.Error(t, err, "An error was expected from Password() with database error")
	assert.Contains(t, err.Error(), "database error", "Error should contain the database error message")

	assert.NoError(t, mock.ExpectationsWereMet(), "All expected mock calls should be made")
}

// TestPasswordInvalidUser tests error handling when user is invalid
func TestPasswordInvalidUser(t *testing.T) {
	user := DefaultUser()
	// Make the user invalid
	user.SetID(2)
	// Not setting authenticated, making the user invalid

	err := user.Password()
	assert.Error(t, err, "An error was expected from Password() with invalid user")
	assert.Equal(t, e.ErrUserNotValid, err, "Error should be ErrUserNotValid")
}

// TestFromNameDatabaseError tests error handling when the database query fails in FromName
func TestFromNameDatabaseError(t *testing.T) {
	user := DefaultUser()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected when setting up mock")

	// Mock a database error
	mock.ExpectQuery("select user_id, user_password from users where user_name").
		WillReturnError(errors.New("database error"))

	err = user.FromName("testuser")
	assert.Error(t, err, "An error was expected from FromName() with database error")
	assert.Contains(t, err.Error(), "database error", "Error should contain the database error message")

	assert.NoError(t, mock.ExpectationsWereMet(), "All expected mock calls should be made")
}

// TestFromNameNoRows tests error handling when no rows are returned (user not found)
func TestFromNameNoRows(t *testing.T) {
	user := DefaultUser()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected when setting up mock")

	// Mock a no rows error
	mock.ExpectQuery("select user_id, user_password from users where user_name").
		WillReturnError(sql.ErrNoRows)

	err = user.FromName("nonexistentuser")
	assert.Error(t, err, "An error was expected from FromName() with no rows")
	assert.Equal(t, sql.ErrNoRows, err, "Error should be sql.ErrNoRows")

	assert.NoError(t, mock.ExpectationsWereMet(), "All expected mock calls should be made")
}

// TestCheckDuplicateDatabaseError tests error handling when the database query fails in CheckDuplicate
func TestCheckDuplicateDatabaseError(t *testing.T) {
	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected when setting up mock")

	// Mock a database error
	mock.ExpectQuery("select count\\(\\*\\) from users where user_name").
		WillReturnError(errors.New("database error"))

	// Should return true (duplicate) on error
	assert.True(t, CheckDuplicate("testuser"), "CheckDuplicate should return true on database error")

	assert.NoError(t, mock.ExpectationsWereMet(), "All expected mock calls should be made")
}

// TestIsAuthorizedDatabaseError tests error handling when the database query fails in IsAuthorized
func TestIsAuthorizedDatabaseError(t *testing.T) {
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected when setting up mock")

	// Mock a database error
	mock.ExpectQuery("SELECT COALESCE").
		WillReturnError(errors.New("database error"))

	// Should return false on error
	assert.False(t, user.IsAuthorized(1), "IsAuthorized should return false on database error")

	assert.NoError(t, mock.ExpectationsWereMet(), "All expected mock calls should be made")
}

// TestIsAuthorizedRoleValues tests authorization with different role values
func TestIsAuthorizedRoleValues(t *testing.T) {
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected when setting up mock")

	// Test with various role values
	roleTests := []struct {
		role     uint
		expected bool
	}{
		{0, false},
		{1, false},
		{2, false},
		{3, true},  // Moderator
		{4, true},  // Admin
		{5, false}, // Any other value
	}

	for _, test := range roleTests {
		rows := sqlmock.NewRows([]string{"role"}).AddRow(test.role)
		mock.ExpectQuery("SELECT COALESCE").WillReturnRows(rows)

		result := user.IsAuthorized(1)
		assert.Equal(t, test.expected, result, "IsAuthorized with role %d should return %t", test.role, test.expected)
	}

	assert.NoError(t, mock.ExpectationsWereMet(), "All expected mock calls should be made")
}

// TestUpdatePasswordDatabaseError tests error handling when the database query fails in UpdatePassword
func TestUpdatePasswordDatabaseError(t *testing.T) {
	_, hash, err := RandomPassword()
	assert.NoError(t, err, "RandomPassword should not error")

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected when setting up mock")

	// Mock a database error
	mock.ExpectExec("UPDATE users SET user_password").
		WithArgs(hash, 2).
		WillReturnError(errors.New("database error"))

	err = UpdatePassword(hash, 2)
	assert.Error(t, err, "An error was expected from UpdatePassword with database error")
	assert.Contains(t, err.Error(), "database error", "Error should contain the database error message")

	assert.NoError(t, mock.ExpectationsWereMet(), "All expected mock calls should be made")
}

// TestSetAuthenticatedWithUserZero tests setting authentication on user 0
func TestSetAuthenticatedWithUserZero(t *testing.T) {
	user := DefaultUser()
	user.SetID(0)
	user.SetAuthenticated()

	assert.False(t, user.IsAuthenticated, "User with ID 0 should not be authenticated")
}

// TestUsernameSanitization tests that username validation properly sanitizes inputs
func TestUsernameSanitization(t *testing.T) {
	// Test various username formats
	validUsernames := []string{
		"alice",
		"bob123",
		"charlie_brown",
		"david-84",
		"eve with spaces",
		"UPPERCASE",
		"Mixed Case 123",
	}

	for _, name := range validUsernames {
		assert.True(t, IsValidName(name), "Username should be valid: %s", name)
	}

	// Test injection attempts
	invalidUsernames := []string{
		"<script>alert(1)</script>",
		"'; DROP TABLE users; --",
		"admin/**/",
		"mod/**/",
		"admin' --",
		"admin\u0000", // null byte
		"admin\n",     // newline
		"\tadmin",     // tab
		"admin;",      // semicolon
		"admin' OR 1=1",
	}

	for _, name := range invalidUsernames {
		assert.False(t, IsValidName(name), "Username should be invalid: %s", name)
	}
}

// TestReservedNamesCaseInsensitive tests that reserved names are rejected regardless of case
func TestReservedNamesCaseInsensitive(t *testing.T) {
	// Test all reserved names with different case variations
	for reservedName := range reservedNameList {
		// Try different case variations
		variations := []string{
			strings.ToLower(reservedName),
			strings.ToUpper(reservedName),
			strings.Title(reservedName),
			// Mix of cases
			func() string {
				mixed := []rune(reservedName)
				for i := range mixed {
					if i%2 == 0 {
						mixed[i] = []rune(strings.ToUpper(string(mixed[i])))[0]
					}
				}
				return string(mixed)
			}(),
			// With spaces
			"  " + reservedName + "  ",
		}

		for _, variant := range variations {
			assert.False(t, IsValidName(variant), "Reserved name should be rejected regardless of case: %s", variant)
		}
	}
}

// TestIsAuthorizedCombinations tests various combinations of validity and authorization
func TestIsAuthorizedCombinations(t *testing.T) {
	testCases := []struct {
		name      string
		setupUser func() User
		setupMock func(mock sqlmock.Sqlmock) sqlmock.Sqlmock
		ibID      uint
		expected  bool
	}{
		{
			name: "Invalid User",
			setupUser: func() User {
				u := DefaultUser()
				u.SetID(0) // Invalid ID
				return u
			},
			setupMock: func(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
				return mock // No expectations needed
			},
			ibID:     1,
			expected: false,
		},
		{
			name: "Invalid Board ID",
			setupUser: func() User {
				u := DefaultUser()
				u.SetID(2)
				u.SetAuthenticated()
				return u
			},
			setupMock: func(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
				return mock // No expectations needed
			},
			ibID:     0, // Invalid board ID
			expected: false,
		},
		{
			name: "DB Error",
			setupUser: func() User {
				u := DefaultUser()
				u.SetID(2)
				u.SetAuthenticated()
				return u
			},
			setupMock: func(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
				mock.ExpectQuery(`SELECT COALESCE`).WillReturnError(sql.ErrConnDone)
				return mock
			},
			ibID:     1,
			expected: false,
		},
		{
			name: "Valid User Moderator Role",
			setupUser: func() User {
				u := DefaultUser()
				u.SetID(2)
				u.SetAuthenticated()
				return u
			},
			setupMock: func(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
				rows := sqlmock.NewRows([]string{"role"}).AddRow(3) // Moderator
				mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)
				return mock
			},
			ibID:     1,
			expected: true,
		},
		{
			name: "Valid User Admin Role",
			setupUser: func() User {
				u := DefaultUser()
				u.SetID(2)
				u.SetAuthenticated()
				return u
			},
			setupMock: func(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
				rows := sqlmock.NewRows([]string{"role"}).AddRow(4) // Admin
				mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)
				return mock
			},
			ibID:     1,
			expected: true,
		},
		{
			name: "Valid User Invalid Role",
			setupUser: func() User {
				u := DefaultUser()
				u.SetID(2)
				u.SetAuthenticated()
				return u
			},
			setupMock: func(mock sqlmock.Sqlmock) sqlmock.Sqlmock {
				rows := sqlmock.NewRows([]string{"role"}).AddRow(99) // Invalid role
				mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)
				return mock
			},
			ibID:     1,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock, err := db.NewTestDb()
			assert.NoError(t, err, "Creating test DB should not error")

			user := tc.setupUser()
			mock = tc.setupMock(mock)

			result := user.IsAuthorized(tc.ibID)
			assert.Equal(t, tc.expected, result, "IsAuthorized should return expected result")

			assert.NoError(t, mock.ExpectationsWereMet(), "All mock expectations should be met")
		})
	}
}

// TestConcurrentPasswordChecks tests that concurrent password checks work correctly
func TestConcurrentPasswordChecks(t *testing.T) {
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	var err error
	user.hash, err = HashPassword("correctpassword")
	assert.NoError(t, err, "Hashing password should not error")

	// Number of concurrent checks
	concurrency := 10

	// Perform concurrent checks for correct password
	var wg sync.WaitGroup
	results := make([]bool, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = user.ComparePassword("correctpassword")
		}(i)
	}

	wg.Wait()

	// All checks should succeed
	for i, result := range results {
		assert.True(t, result, "Password check %d should succeed", i)
	}

	// Reset and check incorrect passwords
	for i := 0; i < concurrency; i++ {
		results[i] = false
	}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = user.ComparePassword("wrongpassword")
		}(i)
	}

	wg.Wait()

	// All checks should fail
	for i, result := range results {
		assert.False(t, result, "Password check %d should fail", i)
	}
}

// TestCheckDuplicatePerformance tests that CheckDuplicate can handle multiple concurrent checks
func TestCheckDuplicatePerformance(t *testing.T) {
	mock, err := db.NewTestDb()
	assert.NoError(t, err, "Creating test DB should not error")

	// Set up expectations for multiple calls
	for i := 0; i < 10; i++ {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery(`select count\(\*\) from users where user_name`).WillReturnRows(rows)
	}

	// Run multiple checks concurrently
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()
			result := CheckDuplicate(username)
			assert.False(t, result, "Username should not be a duplicate")
		}(fmt.Sprintf("testuser%d", i))
	}

	wg.Wait()

	assert.NoError(t, mock.ExpectationsWereMet(), "All mock expectations should be met")
}
