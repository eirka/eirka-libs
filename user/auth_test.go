package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/eirka/eirka-libs/config"
	e "github.com/eirka/eirka-libs/errors"
)

func init() {
	// Enable test mode for secret validation
	SetTestMode(true)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performJWTCookieRequest(r http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.AddCookie(CreateCookie(token))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// resetAuthTestConfig resets the config for auth tests
func resetAuthTestConfig() {
	// Reset secrets
	config.Settings.Session.NewSecret = ""
	config.Settings.Session.OldSecret = ""
}

func TestAuthSecret(t *testing.T) {
	var err error

	// Reset config
	resetAuthTestConfig()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is open
	router.Use(Auth(false))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Test with no secret set
	first := performRequest(router, "GET", "/")
	assert.Equal(t, first.Code, 500, "HTTP request code should match")

	// Set secret in config
	config.Settings.Session.NewSecret = "secret"

	second := performRequest(router, "GET", "/")
	assert.Equal(t, second.Code, 200, "HTTP request code should match")

	// Test with user authentication
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

	third := performJWTCookieRequest(router, "GET", "/", token)
	assert.Equal(t, third.Code, 200, "HTTP request code should match")

	// Reset config
	resetAuthTestConfig()
}

func TestAuthToken(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	first := performRequest(router, "GET", "/")
	assert.Equal(t, first.Code, 403, "HTTP request code should match")

	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	user.hash, err = HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, user.hash, "password should be returned")
	}

	assert.True(t, user.ComparePassword("testpassword"), "Password should validate")

	badtoken, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, badtoken, "token should be returned")
	}

	// Test token validation with secret rotation
	// Change to a new secret and keep old one for rotation
	config.Settings.Session.OldSecret = "secret"
	config.Settings.Session.NewSecret = "changed"

	// Test with a token signed with the old secret (should still work)
	second := performJWTCookieRequest(router, "GET", "/", badtoken)
	assert.Equal(t, second.Code, 200, "HTTP request code should match - old token should work during rotation")

	// Test with malformed tokens
	third := performJWTCookieRequest(router, "GET", "/", "auhwfuiwaehf")
	assert.Equal(t, third.Code, 401, "HTTP request code should match")

	fourth := performJWTCookieRequest(router, "GET", "/", "")
	assert.Equal(t, fourth.Code, 401, "HTTP request code should match")

	// Create token with new secret
	goodtoken, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, goodtoken, "token should be returned")
	}

	fifth := performJWTCookieRequest(router, "GET", "/", goodtoken)
	assert.Equal(t, fifth.Code, 200, "HTTP request code should match")

	// Clear old secret to end rotation
	config.Settings.Session.OldSecret = ""

	// Old token should now fail
	sixth := performJWTCookieRequest(router, "GET", "/", badtoken)
	assert.Equal(t, sixth.Code, 401, "HTTP request code should match - old token should fail after rotation")

	// New token should still work
	seventh := performJWTCookieRequest(router, "GET", "/", goodtoken)
	assert.Equal(t, seventh.Code, 200, "HTTP request code should match - new token should work after rotation")

	// Reset config
	resetAuthTestConfig()
}

func TestAuthValidateToken(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()

	// the current timestamp
	now := time.Now()

	claims := TokenClaims{
		2,
		jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	tkn, err := token.SignedString([]byte("secret"))
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, tkn, "Token should be returned")
	}

	out, err := jwt.ParseWithClaims(tkn, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, &user)
	})
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, out, "Token should be returned")
	}

	// Reset config
	resetAuthTestConfig()
}

func TestAuthValidateTokenNoUser(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()

	// the current timestamp
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    jwtIssuer,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	tkn, err := token.SignedString([]byte("secret"))
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, tkn, "Token should be returned")
	}

	_, err = jwt.ParseWithClaims(tkn, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, &user)
	})
	assert.Error(t, err, "An error was expected")

	// Reset config
	resetAuthTestConfig()
}

func TestAuthValidateTokenBadUser(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()

	// the current timestamp
	now := time.Now()

	claims := TokenClaims{
		1,
		jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	tkn, err := token.SignedString([]byte("secret"))
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, tkn, "Token should be returned")
	}

	_, err = jwt.ParseWithClaims(tkn, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, &user)
	})
	assert.Error(t, err, "An error was expected")
}

func TestAuthValidateTokenNoIssuer(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()

	// the current timestamp
	now := time.Now()

	claims := TokenClaims{
		2,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	tkn, err := token.SignedString([]byte("secret"))
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, tkn, "Token should be returned")
	}

	_, err = jwt.ParseWithClaims(tkn, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, &user)
	})
	assert.Error(t, err, "An error was expected")
}

func TestAuthValidateTokenBadIssuer(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	user := DefaultUser()

	// the current timestamp
	now := time.Now()

	claims := TokenClaims{
		2,
		jwt.RegisteredClaims{
			Issuer:    "derp",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	tkn, err := token.SignedString([]byte("secret"))
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, tkn, "Token should be returned")
	}

	_, err = jwt.ParseWithClaims(tkn, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return validateToken(token, &user)
	})
	assert.Error(t, err, "An error was expected")
}

func TestAuthTokenBadNBF(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// the current timestamp
	now := time.Now()

	claims := TokenClaims{
		2,
		jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now.AddDate(0, 1, 0)),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	tkn, err := token.SignedString([]byte("secret"))
	assert.NoError(t, err, "An error was not expected")

	req := performJWTCookieRequest(router, "GET", "/", tkn)

	assert.Equal(t, req.Code, 401, "HTTP request code should match")
}

func TestAuthTokenExpired(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// the current timestamp
	now := time.Now()

	claims := TokenClaims{
		2,
		jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.AddDate(0, -1, 0)),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	tkn, err := token.SignedString([]byte("secret"))
	assert.NoError(t, err, "An error was not expected")

	req := performJWTCookieRequest(router, "GET", "/", tkn)

	assert.Equal(t, req.Code, 401, "HTTP request code should match")
}

// TestAuthAuthenticatedUserAccessingNonAuthRoute tests that an authenticated user
// can access a route that doesn't require authentication
func TestAuthAuthenticatedUserAccessingNonAuthRoute(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is open (doesn't require authentication)
	router.Use(Auth(false))

	// This will capture the user data from context
	var contextUser User
	router.GET("/", func(c *gin.Context) {
		// Get the user data from context
		userdata, exists := c.Get("userdata")
		if exists {
			contextUser = userdata.(User)
		}
		c.String(200, "OK")
	})

	// Create an authenticated user
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()

	user.hash, err = HashPassword("testpassword")
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotNil(t, user.hash, "password should be returned")
	}

	// Set password as valid
	user.isPasswordValid = true

	// Create a token
	token, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, token, "token should be returned")
	}

	// Perform request with token
	result := performJWTCookieRequest(router, "GET", "/", token)

	// Verify request was successful
	assert.Equal(t, 200, result.Code, "HTTP request code should match")

	// Verify user data was set correctly in context
	assert.True(t, contextUser.IsAuthenticated, "User should be authenticated")
	assert.Equal(t, uint(2), contextUser.ID, "User ID should match")
}

// TestAuthInvalidSigningMethod tests token validation failure due to incorrect signing method
func TestAuthInvalidSigningMethod(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// The current timestamp
	now := time.Now()

	claims := TokenClaims{
		2,
		jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
		},
	}

	// Create token with different signing method
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	token.Header[jwtHeaderKeyID] = 1

	// Sign token with none method
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, tokenString, "Token should be returned")
	}

	// Perform request with invalid token
	result := performJWTCookieRequest(router, "GET", "/", tokenString)

	// Should fail with 401 unauthorized
	assert.Equal(t, 401, result.Code, "HTTP request code should match")
}

// TestAuthInvalidUserID tests token validation failure due to invalid user ID
func TestAuthInvalidUserID(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Test with restricted user IDs
	testUserIDs := []uint{0, 1}

	for _, uid := range testUserIDs {
		// The current timestamp
		now := time.Now()

		claims := TokenClaims{
			uid,
			jwt.RegisteredClaims{
				Issuer:    jwtIssuer,
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * jwtExpireDays)),
			},
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		token.Header[jwtHeaderKeyID] = 1

		primarySecret, err := GetPrimarySecret()
		assert.NoError(t, err)

		tokenString, err := token.SignedString([]byte(primarySecret))
		if assert.NoError(t, err, "An error was not expected") {
			assert.NotEmpty(t, tokenString, "Token should be returned")
		}

		// Perform request with token containing invalid user ID
		result := performJWTCookieRequest(router, "GET", "/", tokenString)

		// Should fail with 401 unauthorized
		assert.Equal(t, 401, result.Code, "HTTP request code should match for user ID "+string(rune(uid)))
	}
}

// TestAuthErrorResponse tests that the correct error code is returned in response
func TestAuthErrorResponse(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Test for no token/cookie
	resultNoToken := performRequest(router, "GET", "/")
	assert.Equal(t, e.ErrForbidden.Code(), resultNoToken.Code, "HTTP request code should match e.ErrForbidden.Code()")
	assert.Contains(t, resultNoToken.Body.String(), e.ErrForbidden.Error(), "Response should contain the correct error message")

	// Test for invalid token
	resultInvalidToken := performJWTCookieRequest(router, "GET", "/", "invalid-token")
	assert.Equal(t, e.ErrUnauthorized.Code(), resultInvalidToken.Code, "HTTP request code should match e.ErrUnauthorized.Code()")
	assert.Contains(t, resultInvalidToken.Body.String(), e.ErrUnauthorized.Error(), "Response should contain the correct error message")

	// Test for missing secret
	resetAuthTestConfig()
	resultNoSecret := performRequest(router, "GET", "/")
	assert.Equal(t, e.ErrInternalError.Code(), resultNoSecret.Code, "HTTP request code should match e.ErrInternalError.Code()")
	assert.Contains(t, resultNoSecret.Body.String(), e.ErrInternalError.Error(), "Response should contain the correct error message")
}

// TestMalformedJWTCookies tests how Auth middleware handles malformed JWT cookies
func TestMalformedJWTCookies(t *testing.T) {
	// Reset and set up config
	resetAuthTestConfig()
	config.Settings.Session.NewSecret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Test cases with different malformed tokens
	malformedTokens := []string{
		"",                                     // Empty token
		"abc",                                  // Too short
		"header.payload",                       // Missing signature
		"header.payload.signature",             // Not base64 encoded parts
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", // Just header
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0", // Missing signature
		"...", // Empty components
	}

	for _, token := range malformedTokens {
		result := performJWTCookieRequest(router, "GET", "/", token)
		assert.Equal(t, e.ErrUnauthorized.Code(), result.Code,
			"Malformed token should return unauthorized: "+token)
		assert.Contains(t, result.Body.String(), e.ErrUnauthorized.Error(),
			"Response should contain correct error message for: "+token)
	}
}

// TestSecretRotation tests token validation during a secret rotation
func TestSecretRotation(t *testing.T) {
	var err error

	// Reset and set up config
	resetAuthTestConfig()

	// Initialize with a primary secret
	config.Settings.Session.NewSecret = "original-secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(Auth(true))
	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Create a user and generate a token with the original secret
	user := DefaultUser()
	user.SetID(2)
	user.SetAuthenticated()
	user.hash, err = HashPassword("testpassword")
	assert.NoError(t, err)
	user.ComparePassword("testpassword")
	oldToken, err := user.CreateToken()
	assert.NoError(t, err)

	// Verify the token works
	result := performJWTCookieRequest(router, "GET", "/", oldToken)
	assert.Equal(t, 200, result.Code, "Token should be valid")

	// Rotate to a new secret
	config.Settings.Session.OldSecret = "original-secret"
	config.Settings.Session.NewSecret = "new-secret-value"

	// The old token should still work because of our rotation support
	result = performJWTCookieRequest(router, "GET", "/", oldToken)
	assert.Equal(t, 200, result.Code, "Old token should still be valid during rotation")

	// Generate a new token with the new secret
	newToken, err := user.CreateToken()
	assert.NoError(t, err)

	// The new token should work
	result = performJWTCookieRequest(router, "GET", "/", newToken)
	assert.Equal(t, 200, result.Code, "New token should be valid")

	// Clear the secondary secret (simulate end of rotation period)
	config.Settings.Session.OldSecret = ""

	// New token still works
	result = performJWTCookieRequest(router, "GET", "/", newToken)
	assert.Equal(t, 200, result.Code, "New token should still be valid after rotation completed")

	// Old token should now fail
	result = performJWTCookieRequest(router, "GET", "/", oldToken)
	assert.Equal(t, 401, result.Code, "Old token should be invalid after rotation completed")

	// Reset config
	resetAuthTestConfig()
}
