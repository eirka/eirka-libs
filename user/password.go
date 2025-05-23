package user

import (
	"crypto/rand"
	"io"

	"golang.org/x/crypto/bcrypt"

	"github.com/eirka/eirka-libs/config"
	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

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

const (
	// characters for random password generator
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

// will generate a password with random characters
func generateRandomPassword(n int) string {
	// byte slice to hold password
	b := make([]byte, n)

	// Read random bytes
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		// If we can't generate random bytes, panic as this is a security issue
		panic(err)
	}

	// Map random bytes to letterBytes character set
	letterLen := len(letterBytes)
	for i := 0; i < n; i++ {
		b[i] = letterBytes[int(b[i])%letterLen]
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
	if len(hash) == 0 {
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
