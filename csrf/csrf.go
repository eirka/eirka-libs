package csrf

import (
	"net/http"
	"time"

	e "github.com/eirka/eirka-libs/errors"
	"github.com/gin-gonic/gin"
)

const (
	// HeaderName is the name of CSRF header
	HeaderName = "X-XSRF-TOKEN"
	// FormFieldName is the name of the form field
	FormFieldName = "csrf_token"
	// CookieName is the name of CSRF cookie
	CookieName = "csrf_token"
	// SessionCookieName the name of the session cookie for angularjs
	SessionCookieName = "XSRF-TOKEN"
)

// skip these methods
var skipMethods = map[string]bool{
	"GET":     true,
	"HEAD":    true,
	"OPTIONS": true,
	"TRACE":   true,
}

// Cookie generates two cookies: a long term csrf token for a user, and a masked session token to verify against
func Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Vary", "Cookie")

		// the token from the users cookie
		var csrfToken []byte

		// get the token from the cookie
		tokenCookie, err := c.Request.Cookie(CookieName)
		if err == nil {
			csrfToken = b64decode(tokenCookie.Value)
		}

		// if the user doesnt have a csrf token create one
		if len(csrfToken) != tokenLength {
			// creates a 32 bit token
			csrfToken = generateToken()

			// set the users csrf token cookie
			csrfCookie := &http.Cookie{
				Name:     CookieName,
				Value:    b64encode(csrfToken),
				Expires:  time.Now().Add(356 * 24 * time.Hour),
				Path:     "/",
				HttpOnly: true,
			}

			// set the csrf token cookie
			http.SetCookie(c.Writer, csrfCookie)

		}

		// generate a session token
		sessionToken := b64encode(maskToken(csrfToken))

		// set the users csrf token tookie
		sessionCookie := &http.Cookie{
			Name:  SessionCookieName,
			Value: sessionToken,
			Path:  "/",
		}

		// set the session cookie
		http.SetCookie(c.Writer, sessionCookie)

		// pass token to controllers
		c.Set("csrf_token", string(sessionToken))

		c.Next()

	}
}

// Verify the sent csrf token
func Verify() gin.HandlerFunc {
	return func(c *gin.Context) {

		// if this is a skippable method
		if skipMethods[c.Request.Method] {
			c.Next()
			return
		}

		// the token from the users cookie
		var csrfToken []byte

		// get the token from the cookie
		tokenCookie, err := c.Request.Cookie(CookieName)
		if err == nil {
			csrfToken = b64decode(tokenCookie.Value)
		}

		var sentToken string

		// Prefer the header over form value
		sentToken = c.Request.Header.Get(HeaderName)

		// Then POST values
		if len(sentToken) == 0 {
			sentToken = c.PostForm(FormFieldName)
		}

		// Then the CSRF session cookie
		if len(sentToken) == 0 {
			sessionCookie, err := c.Request.Cookie(SessionCookieName)
			if err == nil {
				sentToken = sessionCookie.Value
			}
		}

		// sentToken should never be empty at this point so abort
		if len(sentToken) == 0 {
			c.JSON(e.ErrorMessage(e.ErrForbidden))
			c.Error(e.ErrCsrfNotValid).SetMeta("csrf.Verify: No CSRF token found")
			c.Abort()
			return
		}

		// error if there was no csrf token or it isnt verified
		if csrfToken == nil || !verifyToken(csrfToken, b64decode(sentToken)) {
			c.JSON(e.ErrorMessage(e.ErrForbidden))
			c.Error(e.ErrCsrfNotValid).SetMeta("csrf.Verify")
			c.Abort()
			return
		}

		c.Next()

	}
}
