package errors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRequestErrorImplementsError ensures that RequestError implements the Error interface
func TestRequestErrorImplementsError(t *testing.T) {
	var _ error = &RequestError{}
}

// TestRequestErrorCode tests the Code method of the RequestError struct
func TestRequestErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      *RequestError
		expected int
	}{
		{
			name:     "Bad Request",
			err:      &RequestError{ErrorString: "test error", ErrorCode: http.StatusBadRequest},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Internal Server Error",
			err:      &RequestError{ErrorString: "test error", ErrorCode: http.StatusInternalServerError},
			expected: http.StatusInternalServerError,
		},
		{
			name:     "Not Found",
			err:      &RequestError{ErrorString: "test error", ErrorCode: http.StatusNotFound},
			expected: http.StatusNotFound,
		},
		{
			name:     "Unauthorized",
			err:      &RequestError{ErrorString: "test error", ErrorCode: http.StatusUnauthorized},
			expected: http.StatusUnauthorized,
		},
		{
			name:     "Forbidden",
			err:      &RequestError{ErrorString: "test error", ErrorCode: http.StatusForbidden},
			expected: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := tt.err.Code()
			assert.Equal(t, tt.expected, code)
		})
	}
}

// TestRequestErrorError tests the Error method of the RequestError struct
func TestRequestErrorError(t *testing.T) {
	tests := []struct {
		name     string
		err      *RequestError
		expected string
	}{
		{
			name:     "Simple Error Message",
			err:      &RequestError{ErrorString: "simple error", ErrorCode: http.StatusBadRequest},
			expected: "simple error",
		},
		{
			name:     "Empty Error Message",
			err:      &RequestError{ErrorString: "", ErrorCode: http.StatusBadRequest},
			expected: "",
		},
		{
			name:     "Long Error Message",
			err:      &RequestError{ErrorString: "this is a very long error message that provides lots of details about what went wrong", ErrorCode: http.StatusBadRequest},
			expected: "this is a very long error message that provides lots of details about what went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := tt.err.Error()
			assert.Equal(t, tt.expected, message)
		})
	}
}

// TestErrorMessage tests the ErrorMessage function
func TestErrorMessage(t *testing.T) {
	tests := []struct {
		name           string
		err            *RequestError
		expectedCode   int
		expectedErrMsg string
	}{
		{
			name:           "Bad Request",
			err:            &RequestError{ErrorString: "test error", ErrorCode: http.StatusBadRequest},
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "test error",
		},
		{
			name:           "Internal Server Error",
			err:            &RequestError{ErrorString: "internal server error", ErrorCode: http.StatusInternalServerError},
			expectedCode:   http.StatusInternalServerError,
			expectedErrMsg: "internal server error",
		},
		{
			name:           "Not Found",
			err:            &RequestError{ErrorString: "not found", ErrorCode: http.StatusNotFound},
			expectedCode:   http.StatusNotFound,
			expectedErrMsg: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, message := ErrorMessage(tt.err)
			
			assert.Equal(t, tt.expectedCode, code)
			assert.Equal(t, tt.expectedErrMsg, message["error_message"])
		})
	}
}

// TestPredefinedErrors tests the predefined error values
func TestPredefinedErrors(t *testing.T) {
	// Test predefined RequestError instances
	assert.Equal(t, "bad request", ErrInvalidParam.Error())
	assert.Equal(t, http.StatusBadRequest, ErrInvalidParam.Code())
	
	assert.Equal(t, "internal error", ErrInternalError.Error())
	assert.Equal(t, http.StatusInternalServerError, ErrInternalError.Code())
	
	assert.Equal(t, "request not found", ErrNotFound.Error())
	assert.Equal(t, http.StatusNotFound, ErrNotFound.Code())
	
	assert.Equal(t, "unauthorized", ErrUnauthorized.Error())
	assert.Equal(t, http.StatusUnauthorized, ErrUnauthorized.Code())
	
	assert.Equal(t, "forbidden", ErrForbidden.Error())
	assert.Equal(t, http.StatusForbidden, ErrForbidden.Code())
	
	// Test a few standard errors
	assert.Equal(t, "imageboard id required", ErrNoIb.Error())
	assert.Equal(t, "thread id required", ErrNoThread.Error())
	assert.Equal(t, "comment too long", ErrCommentLong.Error())
	assert.Equal(t, "invalid token", ErrTokenInvalid.Error())
	assert.Equal(t, "csrf token is not valid", ErrCsrfNotValid.Error())
}