package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"webframework/errors"
	"webframework/framework"
)

type UserQueryResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func UserQuery(w *framework.ResponseWriter, r *framework.Request) error {
	res := UserQueryResponse{
		Username: "correctName",
		Email:    "q4o5D@example.com",
	}

	w.Header().Set("Content-Type", "application/json")

	return w.Encode(res)
}

func AuthMiddleware(w *framework.ResponseWriter, r *framework.Request, next func()) error {
	token := r.Header.Get("Authorization")
	if token != "secret" {
		return errors.ErrorTypeUnauthorized
	}

	next()
	return nil
}

func TestUserQuery(t *testing.T) {

	router := framework.NewRouter()
	router.Use(framework.JSONBodyEncoder)
	router.Use(AuthMiddleware)
	router.Handle("/query", http.MethodGet, framework.HandlerFunc(UserQuery))

	t.Run("test query user successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req.Header.Set("Authorization", "secret")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		expectedResponse := `{"username":"correctName","email":"q4o5D@example.com"}`

		assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
		assert.JSONEq(t, expectedResponse, rr.Body.String(), "Response body mismatch")
	})

	t.Run("test query user with invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/query", nil)
		req.Header.Set("Authorization", "invalid")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code, "Expected status Unauthorized")
	})
}
