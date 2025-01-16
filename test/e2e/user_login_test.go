package e2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"webframework/errors"
	"webframework/framework"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (c *UserController) Login(w *framework.ResponseWriter, r *framework.Request) error {
	var reqData LoginRequest
	if err := r.ParseBodyInto(&reqData); err != nil {
		return err
	}

	if reqData.Password == "" || reqData.Email == "" {
		return errors.NewError(http.StatusBadRequest, fmt.Errorf("Login's format incorrect."))
	}

	user := c.UserService.FindUserByEmail(reqData.Email)

	if user == nil {
		return errors.NewError(http.StatusUnauthorized, fmt.Errorf("User not found."))
	}

	if user.Password != reqData.Password {
		return errors.NewError(http.StatusUnauthorized, fmt.Errorf("Password incorrect."))
	}

	respData := LoginResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")

	return w.Encode(respData)
}

func TestUserLogin(t *testing.T) {
	userController := NewUserController(userService)
	router := framework.NewRouter()
	router.Use(framework.JSONBodyParser)
	router.Use(framework.JSONBodyEncoder)
	router.Handle("/login", http.MethodPost, framework.HandlerFunc(userController.Login))

	t.Run("test login successfully", func(t *testing.T) {
		loginBody := `{"email": "correctEmail@example.com",  "password": "correctPassword"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		expectedResponse := `{"id": 1, "username": "correctName", "email": "correctEmail@example.com"}`

		assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Expected Content-Type application/json")
		assert.JSONEq(t, expectedResponse, rr.Body.String(), "Response body mismatch, hading: %v", rr.Body.String())
	})
}
