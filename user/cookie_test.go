package user

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateCookie(t *testing.T) {
	// Test with normal token
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNTAwMDAwMDAwfQ.signature"
	cookie := CreateCookie(token)

	// Basic cookie properties
	assert.Equal(t, CookieName, cookie.Name, "Cookie name should match constant")
	assert.Equal(t, token, cookie.Value, "Cookie value should match token")
	assert.Equal(t, "/", cookie.Path, "Cookie path should be root")
	assert.True(t, cookie.HttpOnly, "Cookie should be HttpOnly")

	// Cookie should expire in ~90 days (allow slight variation due to test execution time)
	expectedExpiry := time.Now().Add(90 * 24 * time.Hour)
	timeDiff := expectedExpiry.Sub(cookie.Expires)
	assert.True(t, timeDiff > -time.Second && timeDiff < time.Second,
		"Cookie expiry should be approximately 90 days from now")

	// Test with empty token
	emptyCookie := CreateCookie("")
	assert.Equal(t, "", emptyCookie.Value, "Cookie value should be empty string")

	// Test with very long token
	longToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwiZXhwIjoxNTAwMDAwMDAwfQ." +
		"veryveryveryveryveryveryveryveryveryveryveryveryveryveryveryveryveryverylongsignature"
	longCookie := CreateCookie(longToken)
	assert.Equal(t, longToken, longCookie.Value, "Cookie value should match long token")
}

func TestDeleteCookie(t *testing.T) {
	cookie := DeleteCookie()

	// Basic cookie properties
	assert.Equal(t, CookieName, cookie.Name, "Cookie name should match constant")
	assert.Equal(t, "", cookie.Value, "Cookie value should be empty")
	assert.Equal(t, "/", cookie.Path, "Cookie path should be root")
	assert.True(t, cookie.HttpOnly, "Cookie should be HttpOnly")

	// Cookie should be expired (in the past)
	assert.True(t, cookie.Expires.Before(time.Now()), "Cookie expiry should be in the past")
	assert.Equal(t, -1, cookie.MaxAge, "Cookie MaxAge should be -1")

	// Specific deletion approach: expires approximately 1 year ago
	expectedExpiry := time.Now().AddDate(-1, 0, 0)
	timeDiff := expectedExpiry.Sub(cookie.Expires)
	assert.True(t, timeDiff > -time.Second && timeDiff < time.Second,
		"Cookie expiry should be approximately 1 year ago")
}

func TestCookieUsageInHttpRequest(t *testing.T) {
	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Test adding a session cookie
	token := "test-token"
	req.AddCookie(CreateCookie(token))

	// Verify the cookie was added correctly
	cookie, err := req.Cookie(CookieName)
	assert.NoError(t, err, "Should be able to retrieve cookie")
	assert.Equal(t, token, cookie.Value, "Cookie value should match token")

	// Get the deletion cookie
	deletionCookie := DeleteCookie()

	// Verify the deletion cookie properties
	assert.Equal(t, CookieName, deletionCookie.Name, "Cookie name should match")
	assert.Equal(t, "", deletionCookie.Value, "Deletion cookie value should be empty")
	assert.Equal(t, -1, deletionCookie.MaxAge, "Deletion cookie MaxAge should be -1")
	assert.True(t, deletionCookie.Expires.Before(time.Now()), "Deletion cookie should expire in the past")
}

// TestCookieSecurityFlags tests that cookies have appropriate security flags
func TestCookieSecurityFlags(t *testing.T) {
	cookie := CreateCookie("test-token")

	// Check for security flags
	assert.True(t, cookie.HttpOnly, "Cookie should be HttpOnly to prevent JavaScript access")
	assert.True(t, cookie.Secure, "Cookie should be Secure to prevent transmission over unencrypted connections")
	assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite, "Cookie should use SameSite=Lax mode to prevent CSRF attacks")

	// Also check deletion cookie for the same security flags
	deleteCookie := DeleteCookie()
	assert.True(t, deleteCookie.HttpOnly, "Deletion cookie should be HttpOnly")
	assert.True(t, deleteCookie.Secure, "Deletion cookie should be Secure")
	assert.Equal(t, http.SameSiteLaxMode, deleteCookie.SameSite, "Deletion cookie should use SameSite=Lax mode")
}

// TestCookieDomainAndPath tests domain and path settings
func TestCookieDomainAndPath(t *testing.T) {
	cookie := CreateCookie("test-token")

	// Path should be root
	assert.Equal(t, "/", cookie.Path, "Cookie path should be root")

	// Domain should be empty (which means it will be set to the domain of the request)
	assert.Empty(t, cookie.Domain, "Cookie domain should be empty to default to the current domain")

	// Check deletion cookie as well
	deleteCookie := DeleteCookie()
	assert.Equal(t, "/", deleteCookie.Path, "Deletion cookie path should be root")
	assert.Empty(t, deleteCookie.Domain, "Deletion cookie domain should be empty")
}

// TestCookieExpiration tests the expiration time of cookies
func TestCookieExpiration(t *testing.T) {
	// Create a cookie
	cookie := CreateCookie("test-token")

	// Cookie should expire in ~90 days
	expectedExpiry := time.Now().Add(90 * 24 * time.Hour)
	timeDiff := cookie.Expires.Sub(expectedExpiry)
	assert.True(t, timeDiff > -time.Second*5 && timeDiff < time.Second*5,
		"Cookie expiry should be approximately 90 days from now")

	// Check that MaxAge is not set (should be 0)
	assert.Equal(t, 0, cookie.MaxAge, "Cookie MaxAge should not be set for created cookies")

	// Test deletion cookie
	deleteCookie := DeleteCookie()
	assert.True(t, deleteCookie.Expires.Before(time.Now()), "Deletion cookie should expire in the past")
	assert.Equal(t, -1, deleteCookie.MaxAge, "Deletion cookie MaxAge should be -1")
}

// TestCookieInResponseWriter tests how cookies are set in an HTTP response
func TestCookieInResponseWriter(t *testing.T) {
	// Create a test response recorder
	w := &testResponseWriter{headers: make(http.Header)}

	// Set a cookie
	http.SetCookie(w, CreateCookie("test-token"))

	// Check that the cookie was set in the response headers
	cookies := w.headers.Values("Set-Cookie")
	assert.Equal(t, 1, len(cookies), "There should be 1 Set-Cookie header")

	// The cookie value should contain the token
	assert.Contains(t, cookies[0], "test-token", "Cookie header should contain the token")

	// The cookie should be HttpOnly
	assert.Contains(t, cookies[0], "HttpOnly", "Cookie header should contain HttpOnly flag")

	// Set a deletion cookie
	http.SetCookie(w, DeleteCookie())
	cookies = w.headers.Values("Set-Cookie")
	assert.Equal(t, 2, len(cookies), "There should now be 2 Set-Cookie headers")

	// Second cookie should have Max-Age=0 (browser representation of MaxAge=-1)
	assert.Contains(t, cookies[1], "Max-Age=0", "Deletion cookie should have Max-Age=0 in response")
}

// Mock response writer for testing
type testResponseWriter struct {
	headers http.Header
}

func (w *testResponseWriter) Header() http.Header {
	return w.headers
}

func (w *testResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (w *testResponseWriter) WriteHeader(int) {
}
