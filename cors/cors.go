package cors

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	validSites          = map[string]bool{}
	defaultAllowHeaders = []string{"Origin", "Accept", "Content-Type", "Authorization"}
	defaultAllowMethods []string
)

// CORS will set the headers for Cross-origin resource sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {

		req := c.Request
		origin := req.Header.Get("Origin")

		// Set origin header from sites config
		if isAllowedSite(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Vary", "Origin")

		c.Header("Access-Control-Allow-Credentials", "true")

		if req.Method == "OPTIONS" {

			// Add allowed method header
			c.Header("Access-Control-Allow-Methods", strings.Join(defaultAllowMethods, ","))

			// Add allowed headers header
			c.Header("Access-Control-Allow-Headers", strings.Join(defaultAllowHeaders, ","))

			c.Header("Access-Control-Max-Age", "86400")

			c.AbortWithStatus(http.StatusOK)

			return

		}

		c.Next()

	}
}

// SetDomains will add the valid domains to the map
func SetDomains(domains, methods []string) {
	// add valid sites to map
	for _, site := range domains {
		validSites[site] = true
	}

	// set methods
	defaultAllowMethods = methods

}

// Check if origin is allowed
func isAllowedSite(host string) bool {

	// Get the host from the origin
	parsed, err := url.Parse(host)
	if err != nil {
		return false
	}

	if validSites[strings.ToLower(parsed.Host)] {
		return true
	}

	return false

}
