package user

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	e "github.com/eirka/eirka-libs/errors"
)

const (
	// jwt header keys
	jwtHeaderKeyID = "kid"
	// jwt issuer
	jwtIssuer = "pram"
	// jwt expire days
	jwtExpireDays = 90
)

// TokenClaims holds the custom and standard claims for the JWT token
type TokenClaims struct {
	User uint `json:"user_id"`
	jwt.StandardClaims
}

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

	claims := TokenClaims{
		uid,
		jwt.StandardClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(time.Hour * 24 * jwtExpireDays).Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// set our header info
	token.Header[jwtHeaderKeyID] = 1

	return token.SignedString([]byte(secret))

}
