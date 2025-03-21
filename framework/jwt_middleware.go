package framework

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	werrors "webframework/errors"
)

// JWTAuthMiddleware creates a new middleware for JWT authentication that validates JWT tokens
// in the Authorization header using the provided secret key.
//
// The middleware performs the following checks:
// - Presence of the Authorization header
// - Proper "Bearer" token format
// - JWT token signature validation
// - JWT token expiration validation
//
// On successful validation, the JWT claims are stored in the request context and can be
// retrieved using the GetJWTClaimsFromContext function.
//
// Error Types:
// - ErrorTypeJWTMissing: Authorization header is missing
// - ErrorTypeJWTInvalidFormat: Token doesn't have the "Bearer" prefix
// - ErrorTypeJWTExpired: Token has expired
// - ErrorTypeJWTInvalidSignature: Token signature is invalid
// - ErrorTypeJWTInvalid: Any other JWT validation error
func JWTAuthMiddleware(secretKey []byte) Middleware {
	return func(w *ResponseWriter, r *Request, next func()) error {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return werrors.ErrorTypeJWTMissing
		}

		// Check if the token has Bearer prefix
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return werrors.ErrorTypeJWTInvalidFormat
		}

		tokenString := authHeader[7:]
		claims := jwt.MapClaims{}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		// Handle specific JWT errors
		if err != nil {
			// In jwt v5, we use errors.Is to check for specific errors
			if errors.Is(err, jwt.ErrTokenExpired) {
				return werrors.ErrorTypeJWTExpired
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				return werrors.ErrorTypeJWTInvalidSignature
			} else {
				return werrors.ErrorTypeJWTInvalid
			}
		}

		// Final validation
		if !token.Valid {
			return werrors.ErrorTypeJWTInvalid
		}

		// Store claims in request context
		r.Request = r.Request.WithContext(
			WithJWTClaims(r.Context(), claims),
		)

		next()
		return nil
	}
}
