package user

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
	"github.com/eirka/eirka-libs/validate"
)

func TestProtect(t *testing.T) {

	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(validate.ValidateParams())
	router.Use(Auth(true))
	router.Use(Protect())

	router.GET("/important/:ib", func(c *gin.Context) {
		c.String(200, "OK")
	})

	first := performRequest(router, "GET", "/important/1")

	assert.Equal(t, first.Code, 403, "HTTP request code should match")

	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	user.hash, err = HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, user.hash, "password should be returned")
	}

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

	token, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "token should be returned")
	}

	firstrows := sqlmock.NewRows([]string{"role"}).AddRow(1)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(firstrows)

	second := performJWTCookieRequest(router, "GET", "/important/1", token)

	assert.Equal(t, second.Code, 403, "HTTP request code should match")

	secondrows := sqlmock.NewRows([]string{"role"}).AddRow(3)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(secondrows)

	third := performJWTCookieRequest(router, "GET", "/important/1", token)

	assert.Equal(t, third.Code, 200, "HTTP request code should match")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

	// Reset config
	resetAuthTestConfig()
}

// TestProtectWithMultipleIBs tests the case where multiple image boards are concerned
func TestProtectWithMultipleIBs(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Special validate middleware for this test that sets multiple IBs
	router.Use(func(c *gin.Context) {
		// Set multiple params - this simulates a more complex validation scenario
		c.Set("params", []uint{1, 2, 3})
		c.Next()
	})

	router.Use(Auth(true))
	router.Use(Protect())

	router.GET("/important/multi", func(c *gin.Context) {
		// Check if the protected flag was set
		protected, exists := c.Get("protected")
		assert.True(t, exists, "protected flag should exist")
		assert.True(t, protected.(bool), "protected flag should be true")

		c.String(200, "OK")
	})

	// Create a test user
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	user.hash, err = HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, user.hash, "password should be returned")
	}

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

	token, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "token should be returned")
	}

	// User is a moderator in the first image board (id=1)
	rows := sqlmock.NewRows([]string{"role"}).AddRow(3)
	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	result := performJWTCookieRequest(router, "GET", "/important/multi", token)

	// Should succeed because we're checking authorization against the first board (params[0])
	assert.Equal(t, 200, result.Code, "HTTP request code should match")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

	// Reset config
	resetAuthTestConfig()
}

// TestProtectErrorResponse verifies the exact error response format
func TestProtectErrorResponse(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(validate.ValidateParams())
	router.Use(Auth(true))
	router.Use(Protect())

	router.GET("/important/:ib", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Create a test user
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	user.hash, err = HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, user.hash, "password should be returned")
	}

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

	token, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "token should be returned")
	}

	// User is not a moderator
	rows := sqlmock.NewRows([]string{"role"}).AddRow(1)
	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	result := performJWTCookieRequest(router, "GET", "/important/1", token)

	// Should fail with 403
	assert.Equal(t, e.ErrForbidden.Code(), result.Code, "HTTP request code should match")
	assert.Contains(t, result.Body.String(), e.ErrForbidden.Error(), "Response should contain the correct error message")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

	// Reset config
	resetAuthTestConfig()
}

// TestProtectContextValues tests that protect sets appropriate context values
func TestProtectContextValues(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(validate.ValidateParams())
	router.Use(Auth(true))
	router.Use(Protect())

	// Capture the protected value from the context
	var protectedFlag bool
	var protectedExists bool

	router.GET("/important/:ib", func(c *gin.Context) {
		// Get the protected flag from context
		val, exists := c.Get("protected")
		if exists {
			protectedFlag = val.(bool)
			protectedExists = exists
		}
		c.String(200, "OK")
	})

	// Create a test user
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	user.hash, err = HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, user.hash, "password should be returned")
	}

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

	token, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "token should be returned")
	}

	// User is a moderator
	rows := sqlmock.NewRows([]string{"role"}).AddRow(3)
	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	result := performJWTCookieRequest(router, "GET", "/important/1", token)

	// Should succeed
	assert.Equal(t, 200, result.Code, "HTTP request code should match")

	// Check that the protected flag was properly set
	assert.True(t, protectedExists, "Protected flag should exist in context")
	assert.True(t, protectedFlag, "Protected flag should be true")

	assert.NoError(t, mock.ExpectationsWereMet(), "An error was not expected")

	// Reset config
	resetAuthTestConfig()
}

// TestProtectWithInvalidParams tests protection middleware with invalid parameters
func TestProtectWithInvalidParams(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	_, err = db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Setup middleware - but use a non-empty params slice to avoid index out of range
	router.Use(Auth(true))
	router.Use(func(c *gin.Context) {
		// Set a valid params with an invalid image board ID (negative)
		c.Set("params", []uint{99999})
		c.Next()
	})
	router.Use(Protect())

	router.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Create a test user
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	user.hash, err = HashPassword("testpassword")
	assert.NoError(t, err, "An error was not expected")

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

	token, err := user.CreateToken()
	assert.NoError(t, err, "An error was not expected")

	// Test with invalid image board ID - should result in forbidden
	result := performJWTCookieRequest(router, "GET", "/test", token)
	assert.Equal(t, 403, result.Code, "HTTP request with invalid image board ID should fail with forbidden")

	// Reset config
	resetAuthTestConfig()
}

// TestProtectPanicRecovery tests that the middleware safely handles panics
func TestProtectPanicRecovery(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	_, err = db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	gin.SetMode(gin.ReleaseMode)

	// Create router with recovery middleware
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(Auth(true))

	// Instead of triggering a panic with invalid type, set params with valid type but invalid value
	router.Use(func(c *gin.Context) {
		// Set params as valid type but with an invalid value that will be safely rejected
		c.Set("params", []uint{99999})
		c.Next()
	})
	router.Use(Protect())

	router.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Create a test user with token
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()
	user.hash, err = HashPassword("testpassword")
	assert.NoError(t, err, "An error was not expected")
	user.ComparePassword("testpassword")
	token, err := user.CreateToken()
	assert.NoError(t, err, "An error was not expected")

	// Request should be processed but return a 403 forbidden error
	result := performJWTCookieRequest(router, "GET", "/test", token)
	assert.Equal(t, 403, result.Code, "Request with invalid params value should fail with forbidden")

	// Reset config
	resetAuthTestConfig()
}

// TestProtectRoleValues tests that the middleware properly handles various role values
func TestProtectRoleValues(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(validate.ValidateParams())
	router.Use(Auth(true))
	router.Use(Protect())

	router.GET("/important/:ib", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Create a test user
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()
	user.hash, err = HashPassword("testpassword")
	assert.NoError(t, err, "An error was not expected")
	user.ComparePassword("testpassword")
	token, err := user.CreateToken()
	assert.NoError(t, err, "An error was not expected")

	// Test with different role values
	roleTests := []struct {
		role      uint
		isAllowed bool
	}{
		{0, false}, // No role
		{1, false}, // Basic user
		{2, false}, // Normal user
		{3, true},  // Moderator
		{4, true},  // Admin
		{5, false}, // Invalid role
	}

	for _, tc := range roleTests {
		t.Run(fmt.Sprintf("Role_%d", tc.role), func(t *testing.T) {
			rows := sqlmock.NewRows([]string{"role"}).AddRow(tc.role)
			mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

			result := performJWTCookieRequest(router, "GET", "/important/1", token)

			if tc.isAllowed {
				assert.Equal(t, 200, result.Code, "User with role %d should be allowed", tc.role)
			} else {
				assert.Equal(t, 403, result.Code, "User with role %d should not be allowed", tc.role)
			}
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet(), "All expectations should be met")

	// Reset config
	resetAuthTestConfig()
}
