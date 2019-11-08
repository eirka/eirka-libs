package validate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/eirka/eirka-libs/config"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestValidateParams(t *testing.T) {

	config.Settings.Limits.ParamMaxSize = 10

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// posts need to be verified
	router.Use(ValidateParams())

	router.GET("/index/:id", func(c *gin.Context) {
		c.String(200, "OK")
	})

	first := performRequest(router, "GET", "/index/test")

	assert.Equal(t, first.Code, 400, "HTTP request code should match")

	second := performRequest(router, "GET", "/index/12")

	assert.Equal(t, second.Code, 400, "HTTP request code should match")

	third := performRequest(router, "GET", "/index/1")

	assert.Equal(t, third.Code, 200, "HTTP request code should match")

}
