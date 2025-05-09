package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/eirka/eirka-libs/config"
	e "github.com/eirka/eirka-libs/errors"
)

// Auth is a gin middleware that checks for session cookie and
// handles permissions
func Auth(authenticated bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if secrets are properly configured
		if !IsInitialized() {
			c.JSON(e.ErrorMessage(e.ErrInternalError))
			c.Error(e.ErrNoSecret).SetMeta("auth.Auth.NoSecret")
			c.Abort()
			return
		}

		// set default anonymous user
		user := DefaultUser()

		// try and get the jwt cookie from the request
		cookie, err := c.Request.Cookie(CookieName)

		// parse jwt token if its there
		if err != http.ErrNoCookie {
			// Get all active secrets
			secrets, err := GetSecrets()
			if err != nil {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(err).SetMeta("auth.Auth.GetSecrets")
				c.Abort()
				return
			}

			// First try with new secret (always the first one)
			parseFunc := func(token *jwt.Token) (interface{}, error) {
				return validateToken(token, &user)
			}

			token, parseErr := jwt.ParseWithClaims(cookie.Value, &TokenClaims{}, parseFunc)

			// If token validation failed and old secret is available, try with it
			if parseErr != nil && len(secrets) > 1 {
				// Try with old secret directly
				secondaryFunc := func(token *jwt.Token) (interface{}, error) {
					// Validate algorithm
					_, ok := token.Method.(*jwt.SigningMethodHMAC)
					if !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}

					// Get claims
					claims, ok := token.Claims.(*TokenClaims)
					if !ok {
						return nil, fmt.Errorf("couldn't parse claims")
					}

					// Validate claims
					if claims.Issuer != jwtIssuer {
						return nil, fmt.Errorf("incorrect issuer")
					}

					if claims.User == 0 || claims.User == 1 {
						return nil, fmt.Errorf("invalid user id")
					}

					// Set user info
					user.SetID(claims.User)
					user.SetAuthenticated()

					if !user.IsAuthenticated {
						return nil, fmt.Errorf("user is not authenticated")
					}

					// Return old secret for validation
					return []byte(config.Settings.Session.OldSecret), nil
				}

				// Try parsing with old secret
				token, parseErr = jwt.ParseWithClaims(cookie.Value, &TokenClaims{}, secondaryFunc)
			}

			// If still invalid after all attempts
			if parseErr != nil || !token.Valid {
				// delete the cookie
				http.SetCookie(c.Writer, DeleteCookie())
				c.JSON(e.ErrorMessage(e.ErrUnauthorized))
				c.Error(parseErr).SetMeta("user.Auth")
				c.Abort()
				return
			}
		}

		// check if user needed to be authenticated
		// this needs to be like this for routes that dont need auth
		// if we just check equality then logged in users wont be able
		// to view anon pages ;P
		if authenticated && !user.IsAuthenticated {
			c.JSON(e.ErrorMessage(e.ErrForbidden))
			c.Error(e.ErrForbidden).SetMeta("user.Auth")
			c.Abort()
			return
		}

		// set user data for controllers
		c.Set("userdata", user)

		c.Next()
	}
}
