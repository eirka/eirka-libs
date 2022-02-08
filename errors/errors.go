package errors

import (
	"errors"
	"net/http"
)

var (
	// ErrInvalidParam either means the given route parameter was wrong, like a non uint, or too long
	ErrInvalidParam  = &RequestError{ErrorString: "bad request", ErrorCode: http.StatusBadRequest}
	ErrInternalError = &RequestError{ErrorString: "internal error", ErrorCode: http.StatusInternalServerError}
	ErrNotFound      = &RequestError{ErrorString: "request not found", ErrorCode: http.StatusNotFound}
	// ErrUnauthorized means the user could not be validated and any JWT tokens on client side should be removed
	ErrUnauthorized = &RequestError{ErrorString: "unauthorized", ErrorCode: http.StatusUnauthorized}
	// ErrForbidden is either anon accessing a route that requires auth, or an authed user without the correct permissions
	ErrForbidden = &RequestError{ErrorString: "forbidden", ErrorCode: http.StatusForbidden}

	ErrNoIb             = errors.New("imageboard id required")
	ErrNoThread         = errors.New("thread id required")
	ErrCommentLong      = errors.New("comment too long")
	ErrCommentShort     = errors.New("comment too short")
	ErrNoComment        = errors.New("comment is required")
	ErrTitleLong        = errors.New("title too long")
	ErrTitleShort       = errors.New("title too short")
	ErrNoTitle          = errors.New("title is required")
	ErrNameEmpty        = errors.New("name empty")
	ErrNameLong         = errors.New("name too long")
	ErrNameShort        = errors.New("name too short")
	ErrNameAlphaNum     = errors.New("name not alphanumeric")
	ErrPasswordEmpty    = errors.New("password empty")
	ErrPasswordLong     = errors.New("password too long")
	ErrPasswordShort    = errors.New("password too short")
	ErrNoTagID          = errors.New("tag id required")
	ErrNoTagType        = errors.New("tag type required")
	ErrTagLong          = errors.New("tag too long")
	ErrTagShort         = errors.New("tag too short")
	ErrNoTagName        = errors.New("tag name required")
	ErrDuplicateTag     = errors.New("duplicate tag")
	ErrNoImage          = errors.New("image is required for new threads")
	ErrImageSize        = errors.New("image size is too large")
	ErrDuplicateImage   = errors.New("duplicate image")
	ErrNoImageID        = errors.New("image id required")
	ErrInvalidCookie    = errors.New("invalid cookie")
	ErrNoCookie         = errors.New("cookie required")
	ErrInvalidKey       = errors.New("invalid key")
	ErrNoKey            = errors.New("antispam key required")
	ErrThreadClosed     = errors.New("thread is closed")
	ErrIPParse          = errors.New("input ip cannot be parsed")
	ErrDuplicateName    = errors.New("name already registered")
	ErrInvalidEmail     = errors.New("invalid email address")
	ErrEmailSame        = errors.New("email address the same")
	ErrInvalidUser      = errors.New("user not found")
	ErrInvalidPassword  = errors.New("password incorrect")
	ErrInvalidSession   = errors.New("invalid session")
	ErrMaxLogins        = errors.New("login attempts exceeded")
	ErrUserNotAllowed   = errors.New("username not allowed")
	ErrFavoriteRemoved  = errors.New("favorite removed")
	ErrUserNotConfirmed = errors.New("account not confirmed")
	ErrIPBanned         = errors.New("ip is banned")
	ErrUserBanned       = errors.New("account banned")
	ErrUserLocked       = errors.New("account locked")
	ErrUserNotExist     = errors.New("user does not exist")
	ErrNoSecret         = errors.New("no secret key was set")
	ErrInvalidUID       = errors.New("invalid uid")
	ErrTokenInvalid     = errors.New("invalid token")
	ErrUserNotValid     = errors.New("user is not valid")
	ErrCsrfNotValid     = errors.New("csrf token is not valid")
	ErrBlacklist        = errors.New("ip is on blacklist")
)

// RequestError holds the message string and http code
type RequestError struct {
	ErrorString string
	ErrorCode   int
}

// Code returns the http error code
func (err *RequestError) Code() int {
	return err.ErrorCode
}

func (err *RequestError) Error() string {
	return err.ErrorString
}

// ErrorMessage returns the code and message for Gins JSON helpers
func ErrorMessage(errorType *RequestError) (code int, message map[string]interface{}) {
	return errorType.Code(), map[string]interface{}{"error_message": errorType.Error()}
}
