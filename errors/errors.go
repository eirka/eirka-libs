package errors

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidParam  = &RequestError{ErrorString: "Bad Request", ErrorCode: http.StatusBadRequest}
	ErrInternalError = &RequestError{ErrorString: "Internal Error", ErrorCode: http.StatusInternalServerError}
	ErrNotFound      = &RequestError{ErrorString: "Request Not Found", ErrorCode: http.StatusNotFound}
	ErrUnauthorized  = &RequestError{ErrorString: "Unauthorized", ErrorCode: http.StatusUnauthorized}
	ErrForbidden     = &RequestError{ErrorString: "Forbidden", ErrorCode: http.StatusForbidden}

	ErrNoIb             = errors.New("Imageboard id required")
	ErrNoThread         = errors.New("Thread id required")
	ErrCommentLong      = errors.New("Comment too long")
	ErrCommentShort     = errors.New("Comment too short")
	ErrNoComment        = errors.New("Comment is required")
	ErrTitleLong        = errors.New("Title too long")
	ErrTitleShort       = errors.New("Title too short")
	ErrNoTitle          = errors.New("Title is required")
	ErrNameEmpty        = errors.New("Name empty")
	ErrNameLong         = errors.New("Name too long")
	ErrNameShort        = errors.New("Name too short")
	ErrNameAlphaNum     = errors.New("Name not alphanumeric")
	ErrPasswordEmpty    = errors.New("Password empty")
	ErrPasswordLong     = errors.New("Password too long")
	ErrPasswordShort    = errors.New("Password too short")
	ErrNoTagId          = errors.New("Tag id required")
	ErrNoTagType        = errors.New("Tag type required")
	ErrTagLong          = errors.New("Tag too long")
	ErrTagShort         = errors.New("Tag too short")
	ErrNoTagName        = errors.New("Tag name required")
	ErrDuplicateTag     = errors.New("Duplicate tag")
	ErrNoImage          = errors.New("Image is required for new threads")
	ErrImageSize        = errors.New("Image size is too large")
	ErrDuplicateImage   = errors.New("Duplicate image")
	ErrNoImageId        = errors.New("Image id required")
	ErrInvalidCookie    = errors.New("Invalid cookie")
	ErrNoCookie         = errors.New("Cookie required")
	ErrInvalidKey       = errors.New("Invalid key")
	ErrNoKey            = errors.New("Antispam key required")
	ErrThreadClosed     = errors.New("Thread is closed")
	ErrIpParse          = errors.New("Input IP cannot be parsed")
	ErrDuplicateName    = errors.New("Name already registered")
	ErrInvalidEmail     = errors.New("Invalid email address")
	ErrEmailSame        = errors.New("Email address the same")
	ErrInvalidUser      = errors.New("User not found")
	ErrInvalidPassword  = errors.New("Password incorrect")
	ErrInvalidSession   = errors.New("Invalid session")
	ErrMaxLogins        = errors.New("Login attempts exceeded")
	ErrUserNotAllowed   = errors.New("Username not allowed")
	ErrFavoriteRemoved  = errors.New("Favorite removed")
	ErrUserNotConfirmed = errors.New("Account not confirmed")
	ErrIpBanned         = errors.New("IP is banned")
	ErrUserBanned       = errors.New("Account banned")
	ErrUserLocked       = errors.New("Account locked")
	ErrUserNotExist     = errors.New("User does not exist")
	ErrNoSecret         = errors.New("No secret key was set")
	ErrInvalidUid       = errors.New("Invalid uid")
	ErrTokenInvalid     = errors.New("Invalid token")
	ErrUserNotValid     = errors.New("User is not valid")
	ErrCsrfNotValid     = errors.New("CSRF token is not valid")
)

type RequestError struct {
	ErrorString string
	ErrorCode   int
}

func (err *RequestError) Code() int {
	return err.ErrorCode
}

func (err *RequestError) Error() string {
	return err.ErrorString
}

func ErrorMessage(error_type *RequestError) (code int, message map[string]interface{}) {
	code = error_type.Code()
	message = map[string]interface{}{"error_message": error_type.Error()}

	return
}
