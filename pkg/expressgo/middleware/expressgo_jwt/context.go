package expressgo_jwt

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaimsContextKey is the key type for storing JWT claims in the request context
type JWTClaimsContextKey string

const (
	// JWTClaimsContextKeyValue is the key for storing JWT claims in the request context
	JWTClaimsContextKeyValue JWTClaimsContextKey = "jwt_claims"
)

// WithJWTClaims stores JWT claims in the context
// This function is used internally by the JWTAuthMiddleware to store validated claims
// that can later be retrieved by route handlers
func WithJWTClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, JWTClaimsContextKeyValue, claims)
}

// GetJWTClaimsFromContext extracts JWT claims from the context
//
// This function is used by route handlers to access the JWT claims after the
// JWTAuthMiddleware has successfully validated the token and stored the claims.
//
// Returns:
//   - claims: The JWT claims as a map
//   - ok: Boolean indicating whether claims were found in the context
//
// Example usage:
//
//	func MyHandler(w *framework.ResponseWriter, r *framework.Request) error {
//	    claims, ok := framework.GetJWTClaimsFromContext(r.Context())
//	    if !ok {
//	        return errors.New("JWT claims not found in context")
//	    }
//
//	    // Access claims
//	    userID, _ := claims["user_id"].(string)
//	    return nil
//	}
func GetJWTClaimsFromContext(ctx context.Context) (jwt.MapClaims, bool) {
	claims, ok := ctx.Value(JWTClaimsContextKeyValue).(jwt.MapClaims)
	return claims, ok
}
