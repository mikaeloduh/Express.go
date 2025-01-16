package errors

import "net/http"

var (
	ErrorTypeBadRequest       = NewError(http.StatusBadRequest, nil)       // 400
	ErrorTypeUnauthorized     = NewError(http.StatusUnauthorized, nil)     // 401
	ErrorTypeForbidden        = NewError(http.StatusForbidden, nil)        // 403
	ErrorTypeNotFound         = NewError(http.StatusNotFound, nil)         // 404
	ErrorTypeMethodNotAllowed = NewError(http.StatusMethodNotAllowed, nil) // 405

	ErrorTypeInternalServerError = NewError(http.StatusInternalServerError, nil) // 500
)

// Error
type Error struct {
	Code int
	Err  error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return http.StatusText(e.Code)
}

// NewError
func NewError(code int, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}
