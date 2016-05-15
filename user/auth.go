package user

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	e "github.com/eirka/eirka-libs/errors"
)

// Secret holds the hmac secret, is set from main
var Secret string

// Auth is a gin middleware that checks for session cookie and
// handles permissions
func Auth(authenticated bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// error if theres no secret set
		if Secret == "" {
			c.JSON(e.ErrorMessage(e.ErrInternalError))
			c.Error(e.ErrNoSecret).SetMeta("auth.Auth")
			c.Abort()
			return
		}

		// set default anonymous user
		user := DefaultUser()

		// try and get the jwt cookie from the request
		cookie, err := c.Request.Cookie(CookieName)
		// parse jwt token if its there
		if err != http.ErrNoCookie {
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return validateToken(token, &user)
			})
			// if theres some jwt error other than no token in request or the token is
			// invalid then return unauth
			// the client side should delete any saved JWT tokens on unauth error
			if err != nil || !token.Valid {
				c.JSON(e.ErrorMessage(e.ErrUnauthorized))
				c.Error(err).SetMeta("user.Auth")
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

// validateToken checks all the claims in the provided token
func validateToken(token *jwt.Token, user *User) ([]byte, error) {

	// check alg to make sure its hmac
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	// get the issuer from claims
	tokenIssuer, ok := token.Claims[jwtClaimIssuer].(string)
	if !ok {
		return nil, fmt.Errorf("Couldnt find issuer")
	}

	// check the issuer
	if tokenIssuer != jwtIssuer {
		return nil, fmt.Errorf("Incorrect issuer")
	}

	// get uid from token
	tokenUID, ok := token.Claims[jwtClaimUserID].(float64)
	if !ok {
		return nil, fmt.Errorf("Couldnt find user id")
	}

	// set the user id
	user.SetID(uint(tokenUID))
	// set authenticated
	user.SetAuthenticated()

	// check that the user was actually authed
	if !user.IsAuthenticated {
		return nil, fmt.Errorf("User is not authenticated")
	}

	// compare with secret from settings
	return []byte(Secret), nil

}
