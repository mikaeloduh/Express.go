package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"webframework/framework"
)

// Test secret key
var jwtSecretKey = []byte("jwt-test-secret-key")

// TestJWTAuth tests the JWTAuthMiddleware with the UserQuery handler
func TestJWTAuth(t *testing.T) {
	// Create a JWT helper
	createToken := func(userID string, expiredAt time.Time) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": userID,
			"exp": expiredAt.Unix(),
		})
		
		tokenString, _ := token.SignedString(jwtSecretKey)
		return tokenString
	}

	router := framework.NewRouter()
	router.Use(framework.JSONBodyEncoder)
	router.Use(framework.JWTAuthMiddleware(jwtSecretKey))
	router.Handle("/test-jwt", http.MethodGet, framework.HandlerFunc(func(w *framework.ResponseWriter, r *framework.Request) error {
		// Simple handler that returns a success response
		w.Header().Set("Content-Type", "application/json")
		return w.Encode(map[string]string{"status": "success"})
	}))

	t.Run("test query user with valid JWT token", func(t *testing.T) {
		// Create a valid token
		validToken := createToken("user1", time.Now().Add(time.Hour))
		
		req := httptest.NewRequest(http.MethodGet, "/test-jwt", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		expectedResponse := `{"status":"success"}`

		assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
		assert.JSONEq(t, expectedResponse, rr.Body.String(), "Response body mismatch")
	})

	t.Run("test query user with missing JWT token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-jwt", nil)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
	})

	t.Run("test query user with invalid JWT token format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-jwt", nil)
		req.Header.Set("Authorization", "invalid-token")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
	})

	t.Run("test query user with expired JWT token", func(t *testing.T) {
		// Create an expired token
		expiredToken := createToken("user1", time.Now().Add(-time.Hour))
		
		req := httptest.NewRequest(http.MethodGet, "/test-jwt", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
	})
}
