package user

import (
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	e "github.com/eirka/eirka-libs/errors"
)

// holds the hmac secret, is set from main
var Secret string

// checks for session cookie and handles permissions
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

		// parse jwt token if its there
		token, err := jwt.ParseFromRequest(c.Request, func(token *jwt.Token) (interface{}, error) {

			// check alg to make sure its hmac
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// get the issuer from claims
			token_issuer, ok := token.Claims[jwt_claim_issuer].(string)
			if !ok {
				return nil, fmt.Errorf("Couldnt find issuer")
			}

			// check the issuer
			if token_issuer != jwt_issuer {
				return nil, fmt.Errorf("Incorrect issuer")
			}

			// get uid from token
			token_uid, ok := token.Claims[jwt_claim_user_id].(float64)
			if !ok {
				return nil, fmt.Errorf("Couldnt find user id")
			}

			// set the user id
			user.SetId(uint(token_uid))
			// set authenticated
			user.SetAuthenticated()

			// check that the generated user is valid
			if !user.IsValid() {
				return nil, fmt.Errorf("Generated invalid user")
			}

			// compare with secret from settings
			return []byte(Secret), nil

		})
		// if theres some jwt error other than no token in request or the token is invalid then return unauth
		if err != nil && err != jwt.ErrNoTokenInRequest || token != nil && !token.Valid {
			c.JSON(e.ErrorMessage(e.ErrUnauthorized))
			c.Error(err).SetMeta("auth.Auth")
			c.Abort()
			return
		}

		// check if user needed to be authenticated
		if authenticated && !user.IsAuthenticated {
			c.JSON(e.ErrorMessage(e.ErrUnauthorized))
			c.Error(e.ErrUnauthorized).SetMeta("auth.Auth")
			c.Abort()
			return
		}

		// set user data
		c.Set("userdata", user)

		c.Next()

	}

}
