package errors

import "net/http"

var (
	// JWT specific errors
	ErrorTypeJWTMissing              = NewError(http.StatusUnauthorized, errorWithMsg("JWT token is missing"))
	ErrorTypeJWTInvalidFormat        = NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT format"))
	ErrorTypeJWTInvalid              = NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT token"))
	ErrorTypeJWTExpired              = NewError(http.StatusUnauthorized, errorWithMsg("JWT token has expired"))
	ErrorTypeJWTInvalidSignature     = NewError(http.StatusUnauthorized, errorWithMsg("JWT signature is invalid"))
	ErrorTypeJWTInvalidSigningMethod = NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT signing method"))
)

// errorWithMsg is a simple implementation of the error interface
type errorWithMsg string

func (e errorWithMsg) Error() string {
	return string(e)
}
