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

type UserQueryResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func UserQueryHandler(w *expressgo.Response, _ *expressgo.Request) error {
	res := UserQueryResponse{
		Username: "correctName",
		Email:    "q4o5D@example.com",
	}

	w.Header().Set("Content-Type", "application/json")

	return w.Encode(res)
}

func TestUserQuery(t *testing.T) {
	// JWT secret key for authentication
	var authSecretKey = []byte("auth-secret-key")

	var jwtOptions = jwtmw.Options{
		Keyfunc: func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwtmw.ErrorTypeJWTInvalidSigningMethod
			}
			return authSecretKey, nil
		},
	}

	router := expressgo.NewRouter()
	router.Use(expressgo.JSONBodyEncoder)
	router.Use(jwtmw.JWTAuthMiddleware(jwtOptions))
	router.Handle("/query", http.MethodGet, expressgo.HandlerFunc(UserQueryHandler))

	t.Run("test query user successfully with JWT", func(t *testing.T) {
		// Generate a valid JWT token
		validToken := generateTestJWT("user123", time.Hour, authSecretKey)

		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		expectedResponse := `{"username":"correctName","email":"q4o5D@example.com"}`

		assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
		assert.JSONEq(t, expectedResponse, rr.Body.String(), "Response body mismatch")
	})

	t.Run("test query user with missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		// No Authorization header

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
		// Check response body for specific error message
		assert.Contains(t, rr.Body.String(), "JWT token is missing", "Expected JWT missing error message")
	})

	t.Run("test query user with invalid token format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req.Header.Set("Authorization", "invalid-token")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
		// Check response body for specific error message
		assert.Contains(t, rr.Body.String(), "Invalid JWT format", "Expected invalid JWT format error message")
	})

	t.Run("test query user with expired token", func(t *testing.T) {
		// Generate an expired JWT token
		expiredToken := generateTestJWT("user123", -time.Hour, authSecretKey)

		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req.Header.Set("Authorization", "Bearer "+expiredToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
		// Check response body for specific error message
		assert.Contains(t, rr.Body.String(), "JWT token has expired", "Expected JWT expired error message")
	})

	t.Run("test query user with invalid signature", func(t *testing.T) {
		// Generate a token signed with a different key
		differentSecret := []byte("different-secret-key")
		claims := jwt.MapClaims{
			"sub": "user123",
			"exp": time.Now().Add(time.Hour).Unix(),
			"iat": time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		invalidToken, _ := token.SignedString(differentSecret)

		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req.Header.Set("Authorization", "Bearer "+invalidToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
		// Check response body for specific error message
		assert.Contains(t, rr.Body.String(), "JWT signature is invalid", "Expected invalid signature error message")
	})

	t.Run("test query user with malformed JWT token", func(t *testing.T) {
		// Use a malformed token (missing parts/invalid structure)
		malformedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI.this-is-invalid.and-incomplete"

		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req.Header.Set("Authorization", "Bearer "+malformedToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
		// Check response body for specific error message
		assert.Contains(t, rr.Body.String(), "Invalid JWT token", "Expected invalid JWT error message")
	})
}

// Helper to generate a JWT token for testing
func generateTestJWT(userId string, expiresIn time.Duration, authSecretKey []byte) string {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(expiresIn).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(authSecretKey)
	return tokenString
}
