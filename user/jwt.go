package user

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	e "github.com/eirka/eirka-libs/errors"
)

const (
	// jwt header keys
	jwtHeaderKeyID = "kid"
	// jwt claim keys
	jwtClaimIssuer    = "iss"
	jwtClaimIssued    = "iat"
	jwtClaimNotBefore = "nbf"
	jwtClaimExpire    = "exp"
	jwtClaimUserID    = "user_id"
	// jwt issuer
	jwtIssuer = "pram"
	// jwt expire days
	jwtExpireDays = 90
)

// CreateToken will make a JWT token associated with a user
func (u *User) CreateToken() (newtoken string, err error) {

	// check user struct validity
	if !u.IsValid() {
		err = e.ErrUserNotValid
		return
	}

	// a token should never be created
	if !u.IsAuthenticated {
		err = e.ErrUserNotValid
		return
	}

	// check if password was valid
	if !u.isPasswordValid {
		err = e.ErrInvalidPassword
		return
	}

	return MakeToken(Secret, u.ID)

}

// MakeToken will create a JWT token
func MakeToken(secret string, uid uint) (newtoken string, err error) {

	// error if theres no secret set
	if secret == "" {
		err = e.ErrNoSecret
		return
	}

	// a token should never be created for these users
	if uid == 0 || uid == 1 {
		err = e.ErrUserNotValid
		return
	}

	// the current timestamp
	now := time.Now()

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	// Set our claims
	token.Claims[jwtClaimIssuer] = jwtIssuer
	token.Claims[jwtClaimIssued] = now.Unix()
	token.Claims[jwtClaimNotBefore] = now.Unix()
	token.Claims[jwtClaimExpire] = now.Add(time.Hour * 24 * jwtExpireDays).Unix()
	token.Claims[jwtClaimUserID] = uid

	return token.SignedString([]byte(secret))

}
