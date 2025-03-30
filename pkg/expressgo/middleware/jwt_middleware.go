package middleware

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mikaeloduh/expressgo/pkg/expressgo"
	"github.com/mikaeloduh/expressgo/pkg/expressgo/e"
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
// - ErrorTypeJWTInvalidSigningMethod: JWT signing method is invalid
type Options struct {
	Keyfunc    jwt.Keyfunc
	GetHeader  func(r *expressgo.Request) string
	GetClaims  func(r *expressgo.Request) (jwt.MapClaims, bool)
	SetContext func(ctx context.Context, claims jwt.MapClaims) context.Context
}

// JWTAuthMiddleware creates a middleware that validates JWT tokens
func JWTAuthMiddleware(options Options) expressgo.Middleware {
	// Handle default value logic
	if options.GetHeader == nil {
		options.GetHeader = func(r *expressgo.Request) string {
			return r.Header.Get("Authorization")
		}
	}

	if options.GetClaims == nil {
		options.GetClaims = func(r *expressgo.Request) (jwt.MapClaims, bool) {
			return jwt.MapClaims{}, false
		}
	}

	if options.SetContext == nil {
		options.SetContext = func(ctx context.Context, claims jwt.MapClaims) context.Context {
			return WithJWTClaims(ctx, claims)
		}
	}

	return func(w *expressgo.ResponseWriter, r *expressgo.Request, next func()) error {
		// Extract token from Authorization header
		authHeader := options.GetHeader(r)
		if authHeader == "" {
			return e.ErrorTypeJWTMissing
		}

		// Check if the token has Bearer prefix
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return e.ErrorTypeJWTInvalidFormat
		}

		tokenString := authHeader[7:]
		claims := jwt.MapClaims{}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, claims, options.Keyfunc)

		// Handle specific JWT errors
		if err != nil {
			// In jwt v5, we use errors.Is to check for specific errors
			if errors.Is(err, jwt.ErrTokenExpired) {
				return e.ErrorTypeJWTExpired
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				return e.ErrorTypeJWTInvalidSignature
			} else {
				return e.ErrorTypeJWTInvalid
			}
		}

		// Final validation
		if !token.Valid {
			return e.ErrorTypeJWTInvalid
		}

		// Check if custom claims retrieval is provided and has claims
		if customClaims, ok := options.GetClaims(r); ok {
			for k, v := range customClaims {
				claims[k] = v
			}
		}

		// Store claims in request context
		r.Request = r.Request.WithContext(
			options.SetContext(r.Context(), claims),
		)

		next()
		return nil
	}
}
