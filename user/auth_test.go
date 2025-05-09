package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	e "github.com/eirka/eirka-libs/errors"
)

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

func TestAuthSecret(t *testing.T) {

	var err error

	Secret = ""

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is open
	router.Use(Auth(false))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	first := performRequest(router, "GET", "/")

	assert.Equal(t, first.Code, 500, "HTTP request code should match")

	Secret = "secret"

	second := performRequest(router, "GET", "/")

	assert.Equal(t, second.Code, 200, "HTTP request code should match")

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

}

func TestAuthToken(t *testing.T) {

	var err error

	Secret = "secret"

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

	Secret = "changed"

	second := performJWTCookieRequest(router, "GET", "/", badtoken)

	assert.Equal(t, second.Code, 401, "HTTP request code should match")

	third := performJWTCookieRequest(router, "GET", "/", "auhwfuiwaehf")

	assert.Equal(t, third.Code, 401, "HTTP request code should match")

	fourth := performJWTCookieRequest(router, "GET", "/", "")

	assert.Equal(t, fourth.Code, 401, "HTTP request code should match")

	goodtoken, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, goodtoken, "token should be returned")
	}

	fifth := performJWTCookieRequest(router, "GET", "/", goodtoken)

	assert.Equal(t, fifth.Code, 200, "HTTP request code should match")

}

func TestAuthValidateToken(t *testing.T) {

	Secret = "secret"

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

}

func TestAuthValidateTokenNoUser(t *testing.T) {

	Secret = "secret"

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

}

func TestAuthValidateTokenBadUser(t *testing.T) {

	Secret = "secret"

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

	Secret = "secret"

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

	Secret = "secret"

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

	Secret = "secret"

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

	Secret = "secret"

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

	Secret = "secret"

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
	Secret = "secret"

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
	Secret = "secret"

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

		tokenString, err := token.SignedString([]byte(Secret))
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
	Secret = "secret"

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
	Secret = ""
	resultNoSecret := performRequest(router, "GET", "/")
	assert.Equal(t, e.ErrInternalError.Code(), resultNoSecret.Code, "HTTP request code should match e.ErrInternalError.Code()")
	assert.Contains(t, resultNoSecret.Body.String(), e.ErrInternalError.Error(), "Response should contain the correct error message")
}

// TestMalformedJWTCookies tests how Auth middleware handles malformed JWT cookies
func TestMalformedJWTCookies(t *testing.T) {
	Secret = "secret"

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
