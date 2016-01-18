package csrf

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	csrfCookie   *http.Cookie
	sessionToken string
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performCsrfHeaderRequest(r http.Handler, method, path, token string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	req.AddCookie(csrfCookie)
	req.Header.Set(HeaderName, token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performCsrfFormRequest(r http.Handler, method, path, token string) *httptest.ResponseRecorder {
	var b bytes.Buffer

	mw := multipart.NewWriter(&b)
	mw.WriteField(FormFieldName, token)
	mw.Close()

	req, _ := http.NewRequest(method, path, &b)
	req.AddCookie(csrfCookie)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCsrfVerify(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// posts need to be verified
	router.Use(Verify())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	router.POST("/reply", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	first := performRequest(router, "GET", "/")

	assert.Equal(t, first.Code, 200, "HTTP request code should match")

	second := performRequest(router, "POST", "/reply")

	assert.Equal(t, second.Code, 401, "HTTP request code should match")

}

func TestCsrfCookie(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route issues csrf cookies
	router.Use(Cookie())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	first := performRequest(router, "GET", "/")

	assert.Equal(t, first.Code, 200, "HTTP request code should match")

	assert.Contains(t, first.HeaderMap["Vary"], "Cookie", "Response must include Vary: Cookie header")

	header := http.Header{}

	for _, cookie := range first.HeaderMap["Set-Cookie"] {
		header.Add("Cookie", cookie)
	}

	request := http.Request{Header: header}

	userCookie, err := request.Cookie(CookieName)
	if assert.NoError(t, err, "An error was not expected") {
		assert.Contains(t, userCookie.String(), CookieName, "Response must include user cookie")
		csrfCookie = userCookie
	}

	sessionCookie, err := request.Cookie(SessionName)
	if assert.NoError(t, err, "An error was not expected") {
		assert.Contains(t, sessionCookie.String(), SessionName, "Response must include session cookie")
		sessionToken = sessionCookie.Value
	}

}

func TestCsrfVerifyHeader(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// posts need to be verified
	router.Use(Verify())

	router.POST("/reply", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	badrequest := performCsrfHeaderRequest(router, "POST", "/reply", "badtoken")

	assert.Equal(t, badrequest.Code, 401, "HTTP request code should match")

	goodrequest := performCsrfHeaderRequest(router, "POST", "/reply", sessionToken)

	assert.Equal(t, goodrequest.Code, 200, "HTTP request code should match")

}

func TestCsrfVerifyForm(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// posts need to be verified
	router.Use(Verify())

	router.POST("/reply", func(c *gin.Context) {
		c.String(200, "OK")
		return
	})

	badrequest := performCsrfFormRequest(router, "POST", "/reply", "badtoken")

	assert.Equal(t, badrequest.Code, 401, "HTTP request code should match")

	goodrequest := performCsrfFormRequest(router, "POST", "/reply", sessionToken)

	assert.Equal(t, goodrequest.Code, 200, "HTTP request code should match")

}
