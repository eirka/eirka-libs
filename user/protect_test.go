package user

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

func TestProtect(t *testing.T) {

	var err error

	Secret = "secret"

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"role"}).AddRow(1)

	mock.ExpectQuery(`SELECT COALESCE`).WillReturnRows(rows)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(validate.ValidateParams())
	router.Use(user.Auth(true))
	router.Use(user.Protect())

	router.GET("/important", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	first := performRequest(router, "GET", "/important")

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
		assert.NotEmpty(t, badtoken, "token should be returned")
	}

	second := performJwtHeaderRequest(router, "GET", "/important", token)

}
