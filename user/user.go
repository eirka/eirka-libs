package user

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

const (
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

	// the username validation regex
	username = `^([a-zA-Z0-9]+[\s_-]?)+$`

	// characters for random password generator
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

var regexUsername = regexp.MustCompile(username)

// reserved name list
var reservedNameList = map[string]bool{
	"admin":          true,
	"administrator":  true,
	"administration": true,
	"support":        true,
	"mod":            true,
	"moderator":      true,
	"anon":           true,
	"anonymous":      true,
	"root":           true,
	"webmaster":      true,
	"username":       true,
	"user":           true,
}

// Authenticator defines the methods for authentication
type Authenticator interface {
	IsValid() bool
	IsAuthorized(ib uint) bool
	SetID(uid uint)
	SetAuthenticated()
	Password() (err error)
	ComparePassword(password string) bool
	FromName(name string) (err error)
	CreateToken() (newtoken string, err error)
}

// User data struct
type User struct {
	ID              uint
	Name            string
	IsAuthenticated bool
	hash            []byte
	isPasswordValid bool
}

var _ = Authenticator(&User{})

// DefaultUser creates an anonymous user struct
func DefaultUser() User {
	return User{
		ID:              1,
		IsAuthenticated: false,
	}
}

// SetID sets the user id
func (u *User) SetID(uid uint) {
	u.ID = uid
	return
}

// SetAuthenticated sets a user as authenticated
func (u *User) SetAuthenticated() {
	// do not set auth for the wrong users
	if u.ID == 0 || u.ID == 1 {
		return
	}

	u.IsAuthenticated = true
	return
}

// IsValid will check user struct validity
func (u *User) IsValid() bool {

	// this isnt a real user id
	if u.ID == 0 {
		return false
	}

	// the anon account can never be authenticated
	if u.ID == 1 && u.IsAuthenticated {
		return false
	}

	// a user can never be unauthenticated
	if u.ID != 1 && !u.IsAuthenticated {
		return false
	}

	return true
}

// IsValidName checks if the name is valid
func IsValidName(name string) bool {

	if reservedNameList[strings.ToLower(strings.TrimSpace(name))] {
		return false
	}

	return regexUsername.MatchString(name)
}

// Password will get the password and name from the database for an instantiated user
func (u *User) Password() (err error) {

	// check user struct validity
	if !u.IsValid() {
		return e.ErrUserNotValid
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return
	}

	// get hashed password from database
	err = dbase.QueryRow("select user_name, user_password from users where user_id = ?", u.ID).Scan(&u.Name, &u.hash)
	if err != nil {
		return
	}

	return
}

// FromName will get the password and user id from the database for a user name
func (u *User) FromName(name string) (err error) {

	// name cant be empty
	if len(name) == 0 {
		return e.ErrUserNotValid
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return
	}

	// get hashed password from database
	err = dbase.QueryRow("select user_id, user_password from users where user_name = ?", name).Scan(&u.ID, &u.hash)
	if err != nil {
		return
	}

	u.SetAuthenticated()

	if !u.IsValid() {
		return e.ErrUserNotValid
	}

	return

}

// CheckDuplicate will check for duplicate name before registering
func CheckDuplicate(name string) (check bool) {

	// name cant be empty
	if len(name) == 0 {
		return true
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return true
	}

	// this will return true if there is a user
	err = dbase.QueryRow("select count(*) from users where user_name = ?", name).Scan(&check)
	if err != nil {
		return true
	}

	return

}

// CreateToken will make a JWT token with our claims
func (u *User) CreateToken() (newtoken string, err error) {

	// error if theres no secret set
	if Secret == "" {
		err = e.ErrNoSecret
		return
	}

	// check user struct validity
	if !u.IsValid() {
		err = e.ErrUserNotValid
		return
	}

	// a token should never be created
	if u.ID == 0 || u.ID == 1 {
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

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// the current time
	now := time.Now()

	// Set our claims
	token.Claims[jwtClaimIssuer] = jwtIssuer
	token.Claims[jwtClaimIssued] = now.Unix()
	token.Claims[jwtClaimNotBefore] = now.Unix()
	token.Claims[jwtClaimExpire] = now.Add(time.Hour * 24 * jwtExpireDays).Unix()
	token.Claims[jwtClaimUserID] = u.ID

	// Sign and get the complete encoded token as a string
	newtoken, err = token.SignedString([]byte(Secret))
	if err != nil {
		return
	}

	return

}

// IsAuthorized will get the perms and role info from the userid
func (u *User) IsAuthorized(ib uint) bool {

	var err error

	if !u.IsValid() {
		return false
	}

	// check for invalid stuff
	if ib == 0 {
		return false
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return false
	}

	// holds our role
	var role uint

	// get data from users table
	err = dbase.QueryRow(`SELECT COALESCE((SELECT MAX(role_id) FROM user_ib_role_map WHERE user_ib_role_map.user_id = users.user_id AND ib_id = ?),user_role_map.role_id) as role
    FROM users
    INNER JOIN user_role_map ON (user_role_map.user_id = users.user_id)
    WHERE users.user_id = ?`, ib, u.ID).Scan(&role)
	if err != nil {
		return false
	}

	switch role {
	case 3:
		return true
	case 4:
		return true
	default:
		return false
	}

}

// ComparePassword will compare the supplied password to the hash from the database
func (u *User) ComparePassword(password string) bool {

	// password length cant be 0
	if len(password) == 0 {
		return false
	}

	// if the hash wasnt populated
	if len(u.hash) == 0 {
		return false
	}

	// compare the stored hash with the provided password
	err := bcrypt.CompareHashAndPassword(u.hash, []byte(password))

	// we only want jwt tokens to be created after a valid password has been given
	u.isPasswordValid = err == nil

	return u.isPasswordValid

}

// HashPassword will create a bcrypt hash from the given password
func HashPassword(password string) (hash []byte, err error) {

	// check that password has correct lengths
	if len(password) == 0 {
		err = e.ErrPasswordEmpty
		return
	} else if len(password) < config.Settings.Limits.PasswordMinLength {
		err = e.ErrPasswordShort
		return
	} else if len(password) > config.Settings.Limits.PasswordMaxLength {
		err = e.ErrPasswordLong
		return
	}

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

}

// RandomPassword will generate a random password for password resets
func RandomPassword() (password string, hash []byte, err error) {

	password = generateRandomPassword(20)

	hash, err = HashPassword(password)

	return

}

// will generate a password with random characters
func generateRandomPassword(n int) string {

	// random source
	src := rand.NewSource(time.Now().UnixNano())

	// byte slice to hold password
	b := make([]byte, n)

	// range over byte slice and fill with random letters
	for i := range b {
		b[i] = letterBytes[src.Int63()%int64(len(letterBytes))]
	}

	return string(b)

}

// UpdatePassword will update the user password hash in database
func UpdatePassword(hash []byte, uid uint) (err error) {

	// name cant be empty
	if uid == 0 || uid == 1 {
		return e.ErrUserNotValid
	}

	// hash cant be empty
	if hash == nil || len(hash) == 0 {
		return e.ErrInvalidPassword
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return
	}

	_, err = dbase.Exec("UPDATE users SET user_password = ? WHERE user_id = ?", hash, uid)
	if err != nil {
		return
	}

	return

}
