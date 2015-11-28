package auth

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"

	e "github.com/techjanitor/pram-libs/errors"
)

// holds the hmac secret, is set from main
var Secret string

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
		if err != nil && err != jwt.ErrNoTokenInRequest {
			// if theres some jwt error then return unauth
			c.JSON(e.ErrorMessage(e.ErrUnauthorized))
			c.Error(err)
			c.Abort()
			return
		}

		// process token
		if token != nil {

			// if the token is valid set the data
			if err == nil && token.Valid {

				// get uid from jwt, cast to float
				uid, ok := token.Claims["user_id"].(float64)
				if !ok {
					c.JSON(e.ErrorMessage(e.ErrInternalError))
					c.Error(err)
					c.Abort()
					return
				}

				// set user id in user struct
				user.Id = uint(uid)

				// get the rest of the user info
				err = user.Info()
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error_message": err.Error()})
					c.Error(err)
					c.Abort()
					return
				}

			} else {
				c.JSON(e.ErrorMessage(e.ErrInternalError))
				c.Error(err)
				c.Abort()
				return
			}

		}

		// check if user needed to be authenticated
		if user.IsAuthenticated != authenticated {
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
