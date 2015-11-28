package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"time"

	e "github.com/techjanitor/pram-libs/errors"
)

// Creates a JWT token with our claims
func CreateToken(uid uint) (newtoken string, err error) {

	// error if theres no secret set
	if Secret == "" {
		err = e.ErrNoSecret
		return
	}

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set our claims
	token.Claims["iss"] = "pram"
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(time.Hour * 24 * 90).Unix()
	token.Claims["user_id"] = uid

	// Sign and get the complete encoded token as a string
	newtoken, err = token.SignedString([]byte(Secret))
	if err != nil {
		return
	}

	return

}
