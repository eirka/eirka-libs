package user

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"

	"github.com/eirka/eirka-libs/db"
	"github.com/eirka/eirka-libs/validate"
)

func TestProtect(t *testing.T) {

	var err error

	Secret = "secret"

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(validate.ValidateParams())
	router.Use(Auth(true))
	router.Use(Protect())

	router.GET("/important/:ib", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	first := performRequest(router, "GET", "/important/1")

	assert.Equal(t, first.Code, 401, "HTTP request code should match")

	user := DefaultUser()
	user.SetId(2)
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

	second := performJwtHeaderRequest(router, "GET", "/important/1", token)

	assert.Equal(t, second.Code, 403, "HTTP request code should match")

	secondrows := sqlmock.NewRows([]string{"role"}).AddRow(3)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(secondrows)

	third := performJwtHeaderRequest(router, "GET", "/important/1", token)

	assert.Equal(t, third.Code, 200, "HTTP request code should match")

}
