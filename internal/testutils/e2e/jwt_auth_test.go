package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mikaeloduh/expressgo/pkg/expressgo"
	"github.com/mikaeloduh/expressgo/pkg/expressgo/e"
	"github.com/mikaeloduh/expressgo/pkg/expressgo/middleware/expressgo_jwt"
)

// Test secret key
var jwtSecretKey = []byte("jwt-test-secret-key")

// TestJWTAuth tests the JWTAuthMiddleware with the UserQuery handler
func TestJWTAuth(t *testing.T) {

	jwtOptions := expressgo_jwt.Options{
		Keyfunc: func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, e.ErrorTypeJWTInvalidSigningMethod
			}
			return jwtSecretKey, nil
		},
	}

	router := expressgo.NewRouter()
	router.Use(expressgo.JSONBodyEncoder)
	router.Use(expressgo_jwt.JWTAuthMiddleware(jwtOptions))
	router.Handle("/test-jwt", http.MethodGet, expressgo.HandlerFunc(func(w *expressgo.ResponseWriter, r *expressgo.Request) error {
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

// Create a JWT helper
func createToken(userID string, expiredAt time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": expiredAt.Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecretKey)
	return tokenString
}
