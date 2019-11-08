package user

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
		jwt.StandardClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(time.Hour * 24 * jwtExpireDays).Unix(),
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

	claims := jwt.StandardClaims{
		Issuer:    jwtIssuer,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(time.Hour * 24 * jwtExpireDays).Unix(),
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
		jwt.StandardClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(time.Hour * 24 * jwtExpireDays).Unix(),
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
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(time.Hour * 24 * jwtExpireDays).Unix(),
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
		jwt.StandardClaims{
			Issuer:    "derp",
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(time.Hour * 24 * jwtExpireDays).Unix(),
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
		jwt.StandardClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.AddDate(0, 1, 0).Unix(),
			ExpiresAt: now.Add(time.Hour * 24 * jwtExpireDays).Unix(),
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
		jwt.StandardClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.AddDate(0, -1, 0).Unix(),
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
