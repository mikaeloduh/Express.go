package framework

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	werrors "webframework/errors"
)

func TestJWTMiddleware(t *testing.T) {
	secretKey := []byte("test-secret-key")
	
	// Helper to create a valid JWT token
	createToken := func(userID string, expiredAt time.Time) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": userID,
			"exp": expiredAt.Unix(),
		})
		
		tokenString, _ := token.SignedString(secretKey)
		return tokenString
	}

	t.Run("valid token", func(t *testing.T) {
		// Create a valid token
		validToken := createToken("user1", time.Now().Add(time.Hour))
		
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		
		rw := httptest.NewRecorder()
		w := NewResponseWriter(rw)
		r := &Request{Request: req}
		
		nextCalled := false
		next := func() {
			nextCalled = true
		}
		
		middleware := JWTAuthMiddleware(secretKey)
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
		w := NewResponseWriter(rw)
		r := &Request{Request: req}
		
		nextCalled := false
		next := func() {
			nextCalled = true
		}
		
		middleware := JWTAuthMiddleware(secretKey)
		err := middleware(w, r, next)
		
		assert.Equal(t, werrors.ErrorTypeJWTMissing, err)
		assert.False(t, nextCalled, "Next function should not be called")
	})
	
	t.Run("invalid token format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "invalid-token")
		
		rw := httptest.NewRecorder()
		w := NewResponseWriter(rw)
		r := &Request{Request: req}
		
		nextCalled := false
		next := func() {
			nextCalled = true
		}
		
		middleware := JWTAuthMiddleware(secretKey)
		err := middleware(w, r, next)
		
		assert.Equal(t, werrors.ErrorTypeJWTInvalidFormat, err)
		assert.False(t, nextCalled, "Next function should not be called")
	})
	
	t.Run("expired token", func(t *testing.T) {
		// Create an expired token
		expiredToken := createToken("user1", time.Now().Add(-time.Hour))
		
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)
		
		rw := httptest.NewRecorder()
		w := NewResponseWriter(rw)
		r := &Request{Request: req}
		
		nextCalled := false
		next := func() {
			nextCalled = true
		}
		
		middleware := JWTAuthMiddleware(secretKey)
		err := middleware(w, r, next)
		
		assert.Equal(t, werrors.ErrorTypeJWTExpired, err)
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
		w := NewResponseWriter(rw)
		r := &Request{Request: req}
		
		nextCalled := false
		next := func() {
			nextCalled = true
		}
		
		middleware := JWTAuthMiddleware(secretKey)
		err := middleware(w, r, next)
		
		assert.Equal(t, werrors.ErrorTypeJWTInvalidSignature, err)
		assert.False(t, nextCalled, "Next function should not be called")
	})
}
