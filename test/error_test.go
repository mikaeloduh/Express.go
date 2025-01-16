package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"webframework/errors"
	"webframework/framework"

	"github.com/stretchr/testify/assert"
)

func TestDefaultErrorHandler(t *testing.T) {
	router := framework.NewRouter()
	router.Handle("/test", http.MethodGet, framework.HandlerFunc(func(w *framework.ResponseWriter, r *framework.Request) error {
		w.Write([]byte("OK"))
		return nil
	}))
	// not register any error handlers (use the default error handler)

	t.Run("test default 404 error handling", func(t *testing.T) {

		// Sent a request to a non-existent path
		req := httptest.NewRequest(http.MethodGet, "/non-existent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Verify the response
		assert.Equal(t, http.StatusNotFound, w.Code, "Expected status code %d, got %d", http.StatusNotFound, w.Code)
		assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"), "Expected Content-Type %q, got %q", "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
		assert.Equal(t, "Cannot find the path \"/non-existent\"", w.Body.String(), "Expected error %q, got %q", "404 page not found", w.Body.String())
	})

	t.Run("test default 405 error handling", func(t *testing.T) {

		// Sent a request with an invalid method
		req := httptest.NewRequest(http.MethodDelete, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Verify the response
		assert.Equal(t, http.StatusMethodNotAllowed, w.Code, "Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
		assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"), "Expected Content-Type %q, got %q", "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
		assert.Equal(t, "Method \"DELETE\" is not allowed on path \"test\"", w.Body.String(), "Expected error %q, got %q", "405 method not allowed", w.Body.String())
	})
}

// The custom error handling function for 404 errors
func JSONNotFoundErrorHandler(err error, w *framework.ResponseWriter, r *framework.Request, next func(error)) {
	if e, ok := err.(*errors.Error); ok {
		if e == errors.ErrorTypeNotFound {
			w.WriteHeader(e.Code)
			response := map[string]interface{}{
				"error":   "404 page not found",
				"path":    r.URL.Path,
				"message": e.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

			return
		}

		next(err)
	}
}

// The custom error handling function for 405 errors
func JSONMethodNotAllowedErrorHandler(err error, w *framework.ResponseWriter, r *framework.Request, next func(error)) {
	if e, ok := err.(*errors.Error); ok {
		if e == errors.ErrorTypeMethodNotAllowed {
			w.WriteHeader(e.Code)
			response := map[string]interface{}{
				"error":   "405 method not allowed",
				"path":    r.URL.Path,
				"method":  r.Method,
				"message": e.Error(),
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

			return
		}
	}

	next(err)
}

func TestCustomErrorHandling(t *testing.T) {
	router := framework.NewRouter()

	// register custom error handlers
	router.RegisterErrorHandler(JSONNotFoundErrorHandler)
	router.RegisterErrorHandler(JSONMethodNotAllowedErrorHandler)

	router.Handle("/test", http.MethodGet, framework.HandlerFunc(func(w *framework.ResponseWriter, r *framework.Request) error {
		w.Write([]byte("OK"))
		return nil
	}))

	tests := []struct {
		name           string
		method         string
		path           string
		expectedCode   int
		expectedError  string
		expectedMethod string
	}{
		{
			name:          "404 error - path not found",
			method:        http.MethodGet,
			path:          "/non-existent",
			expectedCode:  http.StatusNotFound,
			expectedError: "404 page not found",
		},
		{
			name:           "405 error - method not allowed",
			method:         http.MethodPost,
			path:           "/test",
			expectedCode:   http.StatusMethodNotAllowed,
			expectedError:  "405 method not allowed",
			expectedMethod: http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// check status code
			assert.Equal(t, tt.expectedCode, w.Code, "Expected status code %d, got %d", tt.expectedCode, w.Code)

			// Check Content-Type
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"), "Expected Content-Type %q, got %q", "application/json", w.Header().Get("Content-Type"))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err, "Failed to decode response: %v", err)

			assert.Equal(t, tt.expectedError, response["error"], "Expected error %q, got %q", tt.expectedError, response["error"])

			assert.Equal(t, tt.path, response["path"], "Expected path %q, got %q", tt.path, response["path"])

			if tt.expectedMethod != "" {
				method, ok := response["method"].(string)
				assert.True(t, ok, "Expected method to be string")
				assert.Equal(t, tt.expectedMethod, method, "Expected method %q, got %q", tt.expectedMethod, method)
			}
		})
	}
}
