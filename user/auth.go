package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"

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
			token, err := jwt.ParseWithClaims(cookie.Value, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return validateToken(token, &user)
			})
			// if theres some jwt error other than no token in request or the token is
			// invalid then return unauth
			// the client side should delete any saved JWT tokens on unauth error
			if err != nil || !token.Valid {
				// delete the cookie
				http.SetCookie(c.Writer, DeleteCookie())
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
