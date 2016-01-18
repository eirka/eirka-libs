package user

import (
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"time"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

const (
	// jwt claim keys
	jwt_claim_issuer    = "iss"
	jwt_claim_issued    = "iat"
	jwt_claim_notbefore = "nbf"
	jwt_claim_expire    = "exp"
	jwt_claim_user_id   = "user_id"
	// jwt issuer
	jwt_issuer = "pram"
	// jwt expire days
	jwt_expire_days = 90

	// the username validation regex
	username = `^([a-zA-Z0-9]+[\s_-]?)+$`
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

// user struct
type User struct {
	Id              uint
	Name            string
	IsAuthenticated bool
	hash            []byte
	isPasswordValid bool
}

// create a user struct
func DefaultUser() User {
	return User{
		Id:              1,
		IsAuthenticated: false,
	}
}

// sets the user id
func (u *User) SetId(uid uint) {
	u.Id = uid
	return
}

// sets authenticated
func (u *User) SetAuthenticated() {
	// do not set auth for the wrong users
	if u.Id == 0 || u.Id == 1 {
		return
	}

	u.IsAuthenticated = true
	return
}

// check user struct validity
func (u *User) IsValid() bool {

	// this isnt a real user id
	if u.Id == 0 {
		return false
	}

	// the anon account can never be authenticated
	if u.Id == 1 && u.IsAuthenticated {
		return false
	}

	// a user can never be unauthenticated
	if u.Id != 1 && !u.IsAuthenticated {
		return false
	}

	return true
}

// checks if the name is valid
func IsValidName(name string) bool {

	if reservedNameList[strings.ToLower(strings.TrimSpace(name))] {
		return false
	}

	return regexUsername.MatchString(name)
}

func (u *User) FromName(name string) (err error) {

	// password length cant be 0
	if len(name) == 0 {
		return e.ErrUserNotValid
	}

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return
	}

	// get hashed password from database
	err = dbase.QueryRow("select user_id, user_password from users where user_name = ?", name).Scan(&u.Id, &u.hash)
	if err != nil {
		return
	}

	u.SetAuthenticated()

	if !u.IsValid() {
		return e.ErrUserNotValid
	}

	return

}

// check for duplicate name before registering
func CheckDuplicate(name string) (check bool) {

	// password length cant be 0
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

// Creates a JWT token with our claims
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
	if u.Id == 0 || u.Id == 1 {
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
	token.Claims[jwt_claim_issuer] = jwt_issuer
	token.Claims[jwt_claim_issued] = now.Unix()
	token.Claims[jwt_claim_notbefore] = now.Unix()
	token.Claims[jwt_claim_expire] = now.Add(time.Hour * 24 * jwt_expire_days).Unix()
	token.Claims[jwt_claim_user_id] = u.Id

	// Sign and get the complete encoded token as a string
	newtoken, err = token.SignedString([]byte(Secret))
	if err != nil {
		return
	}

	return

}

// get the user info from id
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
    WHERE users.user_id = ?`, ib, u.Id).Scan(&role)
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

	return false

}

// compare password to
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

// hash a given password
func HashPassword(password string) (hash []byte, err error) {

	// password length cant be 0
	if len(password) == 0 {
		err = e.ErrInvalidPassword
		return
	}

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
