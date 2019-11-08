package user

import (
	"regexp"
	"strings"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

const (
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
}

// SetAuthenticated sets a user as authenticated
func (u *User) SetAuthenticated() {
	// do not set auth for the wrong users
	if u.ID == 0 || u.ID == 1 {
		return
	}

	u.IsAuthenticated = true
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
