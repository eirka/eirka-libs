package auth

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	e "github.com/eirka/eirka-libs/errors"
)

// holds the hmac secret, is set from main
var Secret string

// user struct
type User struct {
	Id              uint
	IsAuthenticated bool
}

// checks for session cookie and handles permissions
func Auth(authenticated bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// error if theres no secret set
		if Secret == "" {
			c.JSON(e.ErrorMessage(e.ErrInternalError))
			c.Error(e.ErrNoSecret)
			c.Abort()
			return
		}

		// set default anonymous user
		user := User{
			Id:              1,
			IsAuthenticated: false,
		}

		// parse jwt token if its there
		token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {

			// check alg to make sure its hmac
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// compare with secret from settings
			return []byte(Secret), nil

		})
		// if theres some jwt error other than no token in request then return unauth
		if err != nil && err != jwt.ErrNoTokenInRequest {
			c.JSON(e.ErrorMessage(e.ErrUnauthorized))
			c.Error(err)
			c.Abort()
			return
		}

		// if there is a token and its not valid
		if token != nil && err == nil && !token.Valid {
			c.JSON(e.ErrorMessage(e.ErrUnauthorized))
			c.Error(e.ErrTokenInvalid)
			c.Abort()
			return
		}

		// process token if its there and valid
		if token != nil && err == nil && token.Valid {

			// get uid from jwt, cast to float
			jwt_uid, ok := token.Claims[user_id_claim].(float64)
			if !ok {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(e.ErrInternalError)
				c.Abort()
				return
			}

			// cast to uint
			uid := uint(jwt_uid)

			// these are invalid uids
			if uid == 0 || uid == 1 {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(e.ErrInvalidUid)
				c.Abort()
				return
			}

			// set user id in user struct and isauthenticated to true
			user.Id = uid
			user.IsAuthenticated = true

		}

		// check if user needed to be authenticated
		if authenticated && !user.IsAuthenticated {
			c.JSON(e.ErrorMessage(e.ErrUnauthorized))
			c.Error(e.ErrUnauthorized)
			c.Abort()
			return
		}

		// set user data
		c.Set("userdata", user)

		c.Next()

	}

}
