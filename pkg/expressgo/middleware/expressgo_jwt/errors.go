package expressgo_jwt

import (
	"net/http"

	"github.com/mikaeloduh/expressgo/pkg/expressgo/e"
)

// JWT specific errors
var (
	ErrorTypeJWTMissing              = e.NewError(http.StatusUnauthorized, errorWithMsg("JWT token is missing"))
	ErrorTypeJWTInvalidFormat        = e.NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT format"))
	ErrorTypeJWTInvalid              = e.NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT token"))
	ErrorTypeJWTExpired              = e.NewError(http.StatusUnauthorized, errorWithMsg("JWT token has expired"))
	ErrorTypeJWTInvalidSignature     = e.NewError(http.StatusUnauthorized, errorWithMsg("JWT signature is invalid"))
	ErrorTypeJWTInvalidSigningMethod = e.NewError(http.StatusUnauthorized, errorWithMsg("Invalid JWT signing method"))
)

// errorWithMsg is a simple implementation of the error interface
type errorWithMsg string

func (e errorWithMsg) Error() string {
	return string(e)
}
