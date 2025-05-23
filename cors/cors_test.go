package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Run tests
	m.Run()
}

func TestSetDomains(t *testing.T) {
	// Reset validSites map before test
	validSites = map[string]bool{}
	defaultAllowMethods = []string{}

	domains := []string{"example.com", "test.com"}
	methods := []string{"GET", "POST", "OPTIONS"}

	SetDomains(domains, methods)

	// Check if domains were added to validSites
	assert.True(t, validSites["example.com"])
	assert.True(t, validSites["test.com"])
	assert.False(t, validSites["unknown.com"])

	// Check if methods were set correctly
	assert.Equal(t, methods, defaultAllowMethods)
}

func TestIsAllowedSite(t *testing.T) {
	// Reset and initialize validSites map
	validSites = map[string]bool{}
	validSites["example.com"] = true
	validSites["test.com"] = true

	tests := []struct {
		name     string
		origin   string
		expected bool
	}{
		{"Valid origin", "http://example.com", true},
		{"Valid origin with HTTPS", "https://example.com", true},
		{"Valid origin with path", "https://example.com/path", true},
		{"Valid origin with port", "http://example.com:8080", false}, // Different host with port
		{"Valid origin uppercase", "http://EXAMPLE.COM", true},       // Should be case insensitive
		{"Invalid origin", "http://unknown.com", false},
		{"Malformed origin", "not-a-url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAllowedSite(tt.origin)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	// Reset and setup
	validSites = map[string]bool{}
	defaultAllowMethods = []string{}

	domains := []string{"example.com"}
	methods := []string{"GET", "POST", "OPTIONS"}

	SetDomains(domains, methods)

	// Create a gin router with the CORS middleware
	router := gin.New()
	router.Use(CORS())

	// Add a test route
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Test cases
	tests := []struct {
		name                string
		method              string
		origin              string
		expectedStatus      int
		expectedOrigin      string
		expectedCredentials string
		expectedMethods     string
		expectedHeaders     string
		expectedMaxAge      string
	}{
		{
			name:                "Valid GET request",
			method:              "GET",
			origin:              "http://example.com",
			expectedStatus:      http.StatusOK,
			expectedOrigin:      "http://example.com",
			expectedCredentials: "true",
			expectedMethods:     "",
			expectedHeaders:     "",
			expectedMaxAge:      "",
		},
		{
			name:                "Valid OPTIONS request",
			method:              "OPTIONS",
			origin:              "http://example.com",
			expectedStatus:      http.StatusOK,
			expectedOrigin:      "http://example.com",
			expectedCredentials: "true",
			expectedMethods:     "GET,POST,OPTIONS",
			expectedHeaders:     "Origin,Accept,Content-Type,Authorization",
			expectedMaxAge:      "86400",
		},
		{
			name:                "Invalid origin GET request",
			method:              "GET",
			origin:              "http://unknown.com",
			expectedStatus:      http.StatusOK,
			expectedOrigin:      "",
			expectedCredentials: "true",
			expectedMethods:     "",
			expectedHeaders:     "",
			expectedMaxAge:      "",
		},
		{
			name:                "Invalid origin OPTIONS request",
			method:              "OPTIONS",
			origin:              "http://unknown.com",
			expectedStatus:      http.StatusOK,
			expectedOrigin:      "",
			expectedCredentials: "true",
			expectedMethods:     "GET,POST,OPTIONS",
			expectedHeaders:     "Origin,Accept,Content-Type,Authorization",
			expectedMaxAge:      "86400",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request
			req, _ := http.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check CORS headers
			assert.Equal(t, tt.expectedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, tt.expectedCredentials, w.Header().Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, tt.expectedMethods, w.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, tt.expectedHeaders, w.Header().Get("Access-Control-Allow-Headers"))
			assert.Equal(t, tt.expectedMaxAge, w.Header().Get("Access-Control-Max-Age"))
			assert.Equal(t, "Origin", w.Header().Get("Vary"))
		})
	}
}
