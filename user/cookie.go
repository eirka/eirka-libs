package user

import (
	"net/http"
	"time"
)

const (
	// CookieName is the name of the jwt session cookie
	CookieName = "session_jwt"
)

// CreateCookie will make a cookie for the JWT
func CreateCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    token,
		Expires:  time.Now().Add(90 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

// DeleteCookie will delete the JWT cookie
func DeleteCookie() *http.Cookie {
	return &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Expires:  time.Now().AddDate(-1, 0, 0),
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
