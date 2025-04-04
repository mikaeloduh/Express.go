package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mikaeloduh/expressgo"
	jwtmw "github.com/mikaeloduh/expressgo/middleware/jwt"
)

// Test secret key
var jwtSecretKey = []byte("jwt-test-secret-key")

// TestJWTAuth tests the JWTAuthMiddleware with the UserQuery handler
func TestJWTAuth(t *testing.T) {

	jwtOptions := jwtmw.Options{
		Keyfunc: func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwtmw.ErrorTypeJWTInvalidSigningMethod
			}
			return jwtSecretKey, nil
		},
	}

	router := expressgo.NewRouter()
	router.Use(expressgo.JSONBodyEncoder)
	router.Use(jwtmw.AuthMiddleware(jwtOptions))
	router.Handle("/test-jwt", http.MethodGet, expressgo.HandlerFunc(func(req *expressgo.Request, res *expressgo.Response) error {
		// Simple handler that returns a success response
		res.Header().Set("Content-Type", "application/json")
		return res.Encode(map[string]string{"status": "success"})
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
