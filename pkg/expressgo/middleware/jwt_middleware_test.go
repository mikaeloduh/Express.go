package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mikaeloduh/expressgo/pkg/expressgo"
	"github.com/mikaeloduh/expressgo/pkg/expressgo/e"
)

func TestJWTMiddleware(t *testing.T) {
	secretKey := []byte("test-secret-key")
	options := Options{
		Keyfunc: func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, e.ErrorTypeJWTInvalidFormat
			}
			return secretKey, nil
		},
	}

	// Helper to create a valid JWT token
	createToken := func(userID string, expiredAt time.Time) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": userID,
			"exp": expiredAt.Unix(),
		})

		tokenString, _ := token.SignedString(secretKey)
		return tokenString
	}

	t.Run("custom options", func(t *testing.T) {
		// Create a valid token
		validToken := createToken("user1", time.Now().Add(time.Hour))
		customClaims := jwt.MapClaims{"role": "admin"}

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Custom-Auth", "Bearer "+validToken) // Use custom header

		rw := httptest.NewRecorder()
		w := expressgo.NewResponseWriter(rw)
		r := &expressgo.Request{Request: req}

		nextCalled := false
		next := func() {
			nextCalled = true
		}

		// Create custom options
		options := Options{
			Keyfunc: func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			},
			GetHeader: func(r *expressgo.Request) string {
				return r.Header.Get("X-Custom-Auth") // Use custom header
			},
			GetClaims: func(r *expressgo.Request) (jwt.MapClaims, bool) {
				return customClaims, true // Add custom claims
			},
			SetContext: func(ctx context.Context, claims jwt.MapClaims) context.Context {
				return WithJWTClaims(ctx, claims)
			},
		}

		middleware := JWTAuthMiddleware(options)
		err := middleware(w, r, next)

		assert.NoError(t, err)
		assert.True(t, nextCalled, "Next function should be called")

		// Assert the claims were set in the request context
		tokenClaims, ok := GetJWTClaimsFromContext(r.Context())
		assert.True(t, ok, "JWT claims should be set in the request context")
		assert.Equal(t, "user1", tokenClaims["sub"])
		assert.Equal(t, "admin", tokenClaims["role"]) // Check custom claim was merged
	})

	t.Run("valid token", func(t *testing.T) {
		// Create a valid token
		validToken := createToken("user1", time.Now().Add(time.Hour))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rw := httptest.NewRecorder()
		w := expressgo.NewResponseWriter(rw)
		r := &expressgo.Request{Request: req}

		nextCalled := false
		next := func() {
			nextCalled = true
		}

		middleware := JWTAuthMiddleware(options)
		err := middleware(w, r, next)

		assert.NoError(t, err)
		assert.True(t, nextCalled, "Next function should be called")

		// Assert the claims were set in the request context
		tokenClaims, ok := GetJWTClaimsFromContext(r.Context())
		assert.True(t, ok, "JWT claims should be set in the request context")
		assert.Equal(t, "user1", tokenClaims["sub"])
	})

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		rw := httptest.NewRecorder()
		w := expressgo.NewResponseWriter(rw)
		r := &expressgo.Request{Request: req}

		nextCalled := false
		next := func() {
			nextCalled = true
		}

		middleware := JWTAuthMiddleware(options)
		err := middleware(w, r, next)

		assert.Equal(t, e.ErrorTypeJWTMissing, err)
		assert.False(t, nextCalled, "Next function should not be called")
	})

	t.Run("invalid token format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "invalid-token")

		rw := httptest.NewRecorder()
		w := expressgo.NewResponseWriter(rw)
		r := &expressgo.Request{Request: req}

		nextCalled := false
		next := func() {
			nextCalled = true
		}

		middleware := JWTAuthMiddleware(options)
		err := middleware(w, r, next)

		assert.Equal(t, e.ErrorTypeJWTInvalidFormat, err)
		assert.False(t, nextCalled, "Next function should not be called")
	})

	t.Run("expired token", func(t *testing.T) {
		// Create an expired token
		expiredToken := createToken("user1", time.Now().Add(-time.Hour))

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)

		rw := httptest.NewRecorder()
		w := expressgo.NewResponseWriter(rw)
		r := &expressgo.Request{Request: req}

		nextCalled := false
		next := func() {
			nextCalled = true
		}

		middleware := JWTAuthMiddleware(options)
		err := middleware(w, r, next)

		assert.Equal(t, e.ErrorTypeJWTExpired, err)
		assert.False(t, nextCalled, "Next function should not be called")
	})

	t.Run("invalid signature", func(t *testing.T) {
		// Create a token with a different secret key
		differentKey := []byte("different-secret-key")
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "user1",
			"exp": time.Now().Add(time.Hour).Unix(),
		})

		tokenString, _ := token.SignedString(differentKey)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		rw := httptest.NewRecorder()
		w := expressgo.NewResponseWriter(rw)
		r := &expressgo.Request{Request: req}

		nextCalled := false
		next := func() {
			nextCalled = true
		}

		middleware := JWTAuthMiddleware(options)
		err := middleware(w, r, next)

		assert.Equal(t, e.ErrorTypeJWTInvalidSignature, err)
		assert.False(t, nextCalled, "Next function should not be called")
	})
}

func TestJWTMiddlewareWithInvalidOptions(t *testing.T) {
	// Test with invalid options
	emptyOptions := Options{}
	middleware := JWTAuthMiddleware(emptyOptions)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.signature")

	rw := httptest.NewRecorder()
	w := expressgo.NewResponseWriter(rw)
	r := &expressgo.Request{Request: req}

	next := func() {}

	// When Keyfunc is nil, should return an error
	err := middleware(w, r, next)
	assert.Error(t, err, "Should return an error with invalid options")
}
