package jwt

import (
	"net/http"

	"github.com/mikaeloduh/expressgo/e"
)

// JWT specific errors
var (
	ErrorTypeJWTExpired              = e.NewError(http.StatusUnauthorized, errorWithMsg("JWT token has expired"))
	ErrorTypeJWTInvalid              = e.NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT token"))
	ErrorTypeJWTInvalidFormat        = e.NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT format"))
	ErrorTypeJWTInvalidSignature     = e.NewError(http.StatusUnauthorized, errorWithMsg("JWT signature is invalid"))
	ErrorTypeJWTInvalidSigningMethod = e.NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT signing method"))
	ErrorTypeJWTMissing              = e.NewError(http.StatusUnauthorized, errorWithMsg("JWT token is missing"))
)

// errorWithMsg is a simple implementation of the error interface
type errorWithMsg string

func (e errorWithMsg) Error() string {
	return string(e)
}
