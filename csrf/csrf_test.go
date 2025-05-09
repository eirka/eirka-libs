package csrf

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	sessionCookie *http.Cookie
	csrfCookie    *http.Cookie
	sessionToken  string
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

func performCsrfSessionCookieRequest(r http.Handler, method, path string, cookie bool) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)

	if cookie {
		req.AddCookie(sessionCookie)
	}

	req.AddCookie(csrfCookie)
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
	})

	router.POST("/reply", func(c *gin.Context) {
		c.String(200, "OK")
	})

	first := performRequest(router, "GET", "/")

	assert.Equal(t, first.Code, 200, "HTTP request code should match")

	second := performRequest(router, "POST", "/reply")

	assert.Equal(t, second.Code, 403, "HTTP request code should match")

}

func TestCsrfCookie(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// route issues csrf cookies
	router.Use(Cookie())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
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

	sCookie, err := request.Cookie(SessionCookieName)
	if assert.NoError(t, err, "An error was not expected") {
		assert.Contains(t, sCookie.String(), SessionCookieName, "Response must include session cookie")
		sessionToken = sCookie.Value
		sessionCookie = sCookie
	}

}

func TestCsrfVerifyHeader(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// posts need to be verified
	router.Use(Verify())

	router.POST("/reply", func(c *gin.Context) {
		c.String(200, "OK")
	})

	badrequest := performCsrfHeaderRequest(router, "POST", "/reply", "badtoken")

	assert.Equal(t, badrequest.Code, 403, "HTTP request code should match")

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
	})

	badrequest := performCsrfFormRequest(router, "POST", "/reply", "badtoken")

	assert.Equal(t, badrequest.Code, 403, "HTTP request code should match")

	goodrequest := performCsrfFormRequest(router, "POST", "/reply", sessionToken)

	assert.Equal(t, goodrequest.Code, 200, "HTTP request code should match")

}

func TestCsrfVerifySessionCookiePost(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// posts need to be verified
	router.Use(Verify())

	router.POST("/reply", func(c *gin.Context) {
		c.String(200, "OK")
	})

	badrequest := performCsrfSessionCookieRequest(router, "POST", "/reply", false)

	assert.Equal(t, badrequest.Code, 403, "HTTP request code should match")

	goodrequest := performCsrfSessionCookieRequest(router, "POST", "/reply", true)

	assert.Equal(t, goodrequest.Code, 200, "HTTP request code should match")

}

// Helper function for testing various HTTP methods with CSRF verification
func performCsrfMethod(t *testing.T, method string, shouldPass bool) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(Verify())
	router.Handle(method, "/endpoint", func(c *gin.Context) {
		c.String(200, "OK")
	})

	var resp *httptest.ResponseRecorder
	if shouldPass {
		// For GET/HEAD/OPTIONS/TRACE methods, we don't need the CSRF token
		resp = performRequest(router, method, "/endpoint")
	} else {
		// For other methods without CSRF token should fail
		resp = performRequest(router, method, "/endpoint")
	}

	if shouldPass {
		assert.Equal(t, 200, resp.Code, "HTTP request code should be 200 for method "+method)
	} else {
		if method == "GET" || method == "HEAD" || method == "OPTIONS" || method == "TRACE" {
			assert.Equal(t, 200, resp.Code, "HTTP request code should be 200 for skipped method "+method)
		} else {
			assert.Equal(t, 403, resp.Code, "HTTP request code should be 403 for method "+method+" without valid CSRF")
		}
	}
}

func TestCsrfMethodVerification(t *testing.T) {
	// Test methods that should be skipped
	performCsrfMethod(t, "GET", true)
	performCsrfMethod(t, "HEAD", true)
	performCsrfMethod(t, "OPTIONS", true)
	performCsrfMethod(t, "TRACE", true)

	// Test methods that should be verified
	performCsrfMethod(t, "POST", false)
	performCsrfMethod(t, "PUT", false)
	performCsrfMethod(t, "DELETE", false)
	performCsrfMethod(t, "PATCH", false)
}

// Test verifying an invalid cookie
func TestCsrfInvalidCookie(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(Verify())

	router.POST("/reply", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Create a fake request with an invalid cookie
	req, _ := http.NewRequest("POST", "/reply", nil)
	invalidCookie := &http.Cookie{
		Name:  CookieName,
		Value: "invalidvalue",
	}
	req.AddCookie(invalidCookie)
	req.Header.Set(HeaderName, sessionToken)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code, "HTTP request code should be 403 with invalid cookie")
}

// Test multiple CSRF cookies in the same request
func TestCsrfMultipleCookies(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	// Generate CSRF cookies
	router.Use(Cookie())
	
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// First request to get cookies
	first := performRequest(router, "GET", "/test")
	assert.Equal(t, 200, first.Code, "HTTP request code should match")

	// Second request with cookies already set
	req, _ := http.NewRequest("GET", "/test", nil)
	req.AddCookie(csrfCookie)
	req.AddCookie(sessionCookie)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "HTTP request code should match")

	// Verify session cookie is always regenerated
	sessionCookieFound := false
	for _, cookie := range w.HeaderMap["Set-Cookie"] {
		if strings.Contains(cookie, SessionCookieName) {
			sessionCookieFound = true
			break
		}
	}
	assert.True(t, sessionCookieFound, "A session cookie should always be generated")
}

// Test cookie with invalid length
func TestCsrfInvalidCookieLength(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(Cookie())
	
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Create a request with an invalid length cookie
	req, _ := http.NewRequest("GET", "/test", nil)
	invalidCookie := &http.Cookie{
		Name:  CookieName,
		Value: "tooShort", // Not a valid base64 encoded token
	}
	req.AddCookie(invalidCookie)
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code, "HTTP request code should match")

	// A new cookie should be generated
	cookieFound := false
	for _, cookie := range w.HeaderMap["Set-Cookie"] {
		if cookie != invalidCookie.String() && cookie != sessionCookie.String() {
			cookieFound = true
			break
		}
	}
	assert.True(t, cookieFound, "A new cookie should be generated when an invalid length cookie is provided")
}

// Test for empty token
func TestCsrfEmptyToken(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(Verify())

	router.POST("/reply", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Create a request with no token
	req, _ := http.NewRequest("POST", "/reply", nil)
	req.AddCookie(csrfCookie)
	// Deliberately not setting any CSRF tokens
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code, "HTTP request code should be 403 with no CSRF token")
}

// Test for b64encode and b64decode functions
func TestB64EncodeDecode(t *testing.T) {
	testData := []byte("TestDataForEncoding12345")
	
	// Test encode
	encoded := b64encode(testData)
	assert.NotEmpty(t, encoded, "Encoded string should not be empty")
	
	// Test decode
	decoded := b64decode(encoded)
	assert.Equal(t, testData, decoded, "Decoded data should match original")
	
	// Test decode with invalid base64
	invalidDecoded := b64decode("this is not valid base64!!!!!")
	assert.Nil(t, invalidDecoded, "Invalid base64 should return nil")
}