package user

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"

	local "github.com/eirka/eirka-post/config"
)

func init() {

	// Database connection settings
	dbase := db.Database{

		User:           local.Settings.Database.User,
		Password:       local.Settings.Database.Password,
		Proto:          local.Settings.Database.Proto,
		Host:           local.Settings.Database.Host,
		Database:       local.Settings.Database.Database,
		MaxIdle:        local.Settings.Database.MaxIdle,
		MaxConnections: local.Settings.Database.MaxConnections,
	}

	// Set up DB connection
	dbase.NewDb()

	// Get limits and stuff from database
	config.GetDatabaseSettings()
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performJwtHeaderRequest(r http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performJwtFormRequest(r http.Handler, method, path, token string) *httptest.ResponseRecorder {
	var b bytes.Buffer

	mw := multipart.NewWriter(&b)
	mw.WriteField("access_token", token)
	mw.Close()

	req, _ := http.NewRequest(method, path, &b)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthSecret(t *testing.T) {

	Secret = ""

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is open
	router.Use(Auth(false))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	first := performRequest(router, "GET", "/")

	assert.Equal(t, first.Code, 500, "HTTP request code should match")

	Secret = "secret"

	second := performRequest(router, "GET", "/")

	assert.Equal(t, second.Code, 200, "HTTP request code should match")

}

func TestAuthHeaderToken(t *testing.T) {

	var err error

	Secret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	first := performRequest(router, "GET", "/")

	assert.Equal(t, first.Code, 401, "HTTP request code should match")

	user := DefaultUser()
	user.SetId(2)
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

	second := performJwtHeaderRequest(router, "GET", "/", badtoken)

	assert.Equal(t, second.Code, 401, "HTTP request code should match")

	goodtoken, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, goodtoken, "token should be returned")
	}

	third := performJwtHeaderRequest(router, "GET", "/", goodtoken)

	assert.Equal(t, third.Code, 200, "HTTP request code should match")

}

func TestAuthFormToken(t *testing.T) {

	var err error

	Secret = "secret"

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route is protected
	router.Use(Auth(true))

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	user := DefaultUser()
	user.SetId(2)
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

	first := performJwtFormRequest(router, "GET", "/", badtoken)

	assert.Equal(t, first.Code, 401, "HTTP request code should match")

	goodtoken, err := user.CreateToken()
	if assert.NoError(t, err, "An error was not expected") {
		assert.NotEmpty(t, goodtoken, "token should be returned")
	}

	second := performJwtFormRequest(router, "GET", "/", goodtoken)

	assert.Equal(t, second.Code, 200, "HTTP request code should match")

}
