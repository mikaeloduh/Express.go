package expressgo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Handler functions remain the same
func homeHandler(w *ResponseWriter, r *Request) error {
	fmt.Fprintf(w, "Welcome to the homepage!")
	return nil
}

func helloHandler(w *ResponseWriter, r *Request) error {
	fmt.Fprintf(w, "Hello, World!")
	return nil
}

func getUserHandler(w *ResponseWriter, r *Request) error {
	fmt.Fprintf(w, "Retrieve user information")
	return nil
}

func postUserHandler(w *ResponseWriter, r *Request) error {
	fmt.Fprintf(w, "Create a new user")
	return nil
}

func userProfileHandler(w *ResponseWriter, r *Request) error {
	fmt.Fprintf(w, "User profile page")
	return nil
}

// TestRouting verifies that the routing is correctly set up and that each route returns the expected response.
func TestRouting(t *testing.T) {
	// Create a new Router
	route := NewRouter()

	// Register handlers with exact path and method matching
	route.Handle("/", http.MethodGet, HandlerFunc(homeHandler))
	route.Handle("/hello", http.MethodGet, HandlerFunc(helloHandler))
	route.Handle("/user", http.MethodGet, HandlerFunc(getUserHandler))
	route.Handle("/user", http.MethodPost, HandlerFunc(postUserHandler))
	route.Handle("/user/profile", http.MethodGet, HandlerFunc(userProfileHandler))

	// Start a new test server using the custom Router
	ts := httptest.NewServer(route)
	defer ts.Close()

	// Define test cases
	tests := []struct {
		method       string
		path         string
		expectedCode int
		expectedBody string
	}{
		{"GET", "/", http.StatusOK, "Welcome to the homepage!"},
		{"GET", "/hello", http.StatusOK, "Hello, World!"},
		{"GET", "/user", http.StatusOK, "Retrieve user information"},
		{"POST", "/user", http.StatusOK, "Create a new user"},
		{"PUT", "/user", http.StatusMethodNotAllowed, "Method \"PUT\" is not allowed on path \"user\""},
		{"GET", "/user/profile", http.StatusOK, "User profile page"},
		{"GET", "/usersomething", http.StatusNotFound, "Cannot find the path \"/usersomething\""},
		{"GET", "/userextra", http.StatusNotFound, "Cannot find the path \"/userextra\""},
	}

	for _, tc := range tests {
		// Create a new HTTP request
		req, err := http.NewRequest(tc.method, ts.URL+tc.path, nil)
		assert.NoError(t, err)

		// Send the request
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		// Check the status code
		assert.Equal(t, tc.expectedCode, resp.StatusCode, "Unexpected status code for path %s", tc.path)

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		resp.Body.Close()

		// Check the response body
		assert.Equal(t, tc.expectedBody, string(body), "Unexpected response body for path %s", tc.path)
	}
}

// Create a test middleware that checks for a specific header
func testMiddleware(w *ResponseWriter, r *Request, next func()) error {
	if r.Header.Get("X-Test") != "test-value" {
		return fmt.Errorf("missing or invalid X-Test header")
	}
	next()
	return nil
}

// TestMiddleware verifies that middleware functions are correctly applied
func TestMiddleware(t *testing.T) {
	// Create a new Router
	route := NewRouter()

	// Register the middleware
	route.Use(testMiddleware)

	// Register a simple handler
	route.Handle("/test", http.MethodGet, HandlerFunc(func(w *ResponseWriter, r *Request) error {
		fmt.Fprintf(w, "test passed")
		return nil
	}))

	// Start a new test server
	ts := httptest.NewServer(route)
	defer ts.Close()

	// Test case 1: Request without required header should fail
	t.Run("without header", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/test", nil)
		assert.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	// Test case 2: Request with required header should succeed
	t.Run("with header", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/test", nil)
		assert.NoError(t, err)
		req.Header.Set("X-Test", "test-value")

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "test passed", string(body))
	})
}
