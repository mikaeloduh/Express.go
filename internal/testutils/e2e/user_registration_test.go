package e2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mikaeloduh/expressgo"
	"github.com/mikaeloduh/expressgo/e"
	"github.com/mikaeloduh/expressgo/middleware/body_parser"
)

type RegisterRequest struct {
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
	Password string `json:"password" xml:"password"`
}

type RegisterResponse struct {
	Id       uint64 `json:"id" xml:"id"`
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
}

func (c *UserController) Register(w *expressgo.ResponseWriter, r *expressgo.Request) error {
	var reqData RegisterRequest
	if err := r.ParseBodyInto(&reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	if reqData.Username == "" || reqData.Email == "" || reqData.Password == "" {
		return e.NewError(http.StatusBadRequest, fmt.Errorf("Registration's format incorrect."))
	}

	if c.UserService.FindUserByEmail(reqData.Email) != nil {
		return e.NewError(http.StatusBadRequest, fmt.Errorf("Duplicate email"))
	}

	user, _ := c.UserService.CreateUser(reqData.Username, reqData.Email, reqData.Password)

	respData := RegisterResponse{
		Id:       user.Id,
		Username: user.Username,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")

	return w.Encode(respData)
}

func TestRegisterHandlerJSON(t *testing.T) {
	userController := NewUserController(userService)
	router := expressgo.NewRouter()
	router.Use(body_parser.JSONBodyParser)
	router.Use(expressgo.JSONBodyEncoder)
	router.Use(body_parser.XMLBodyParser)
	router.Use(expressgo.XMLBodyEncoder)
	router.Handle("/register", http.MethodPost, expressgo.HandlerFunc(userController.Register))

	t.Run("test register user successfully", func(t *testing.T) {
		jsonBody := `{"username": "John Doe", "email": "jdoe@example.com", "password": "abc"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		expectedResponse := `{"username": "John Doe", "email": "jdoe@example.com", "id": 2}`

		assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Expected Content-Type application/json")
		assert.JSONEq(t, expectedResponse, rr.Body.String(), "Response body mismatch")
	})

	t.Run("register fail: email exists", func(t *testing.T) {
		jsonBody := `{"username": "John Doe", "email": "jdoe@example.com", "password": "abc"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status BadRequest")
		assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"), "Expected Content-Type text/plain")
		assert.Equal(t, "Duplicate email", rr.Body.String(), "Response body mismatch")
	})

	t.Run("register fail: invalid format", func(t *testing.T) {
		jsonBody := `{"incorrectrequest": "this is a incorrect test request"}`
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status BadRequest")
		assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"), "Expected Content-Type text/plain")
		assert.Equal(t, "Registration's format incorrect.", rr.Body.String(), "Response body mismatch")
	})
}
