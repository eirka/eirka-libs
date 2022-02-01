package user

import (
	"fmt"
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

// validateToken checks all the claims in the provided token
func validateToken(token *jwt.Token, user *User) ([]byte, error) {

	// check alg to make sure its hmac
	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	if !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	// get the claims from the token
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, fmt.Errorf("couldnt parse claims")
	}

	// get the issuer from claims
	tokenIssuer := claims.Issuer

	// check the issuer
	if tokenIssuer != jwtIssuer {
		return nil, fmt.Errorf("incorrect issuer")
	}

	// get uid from token
	tokenUID := claims.User

	// set the user id
	user.SetID(uint(tokenUID))
	// set authenticated
	user.SetAuthenticated()

	// check that the user was actually authed
	if !user.IsAuthenticated {
		return nil, fmt.Errorf("user is not authenticated")
	}

	// compare with secret from settings
	return []byte(Secret), nil

}
