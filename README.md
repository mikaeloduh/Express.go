# Express.go

A lightweight and flexible web framework for Go, inspired by Express.js. Express.go provides a robust set of features for building web applications and APIs with an elegant and intuitive interface.

## Features

- **Middleware Support**: Flexible middleware system for request/response handling
- **Routing**: Simple and intuitive routing system
- **Error Handling**: Built-in error handling middleware
- **Body Parsing**: Support for JSON and XML request/response parsing
- **Dependency Injection**: Built-in dependency injection container with different scopes
- **JWT Authentication**: Comprehensive JWT authentication middleware with granular error handling
- **Rate Limiting**: Built-in rate limiting middleware
- **Extensible**: Easy to extend with custom middleware and handlers

## Installation

```bash
go get webframework
```

## Quick Start

```go
package main

import "webframework/framework"

func main() {
    router := framework.NewRouter()
    
    // Add middleware
    router.Use(framework.JSONBodyParser)
    
    // Define routes
    router.Handle("/", "GET", func(w *framework.ResponseWriter, r *framework.Request) error {
        return w.Encode(map[string]string{"message": "Hello, World!"})
    })
    
    // Start server
    http.ListenAndServe(":3000", router)
}
```

## Documentation

### Router

The router is the core of Express.go. It handles HTTP requests and routes them to the appropriate handlers.

```go
router := framework.NewRouter()
router.Handle("/path", "GET", handlerFunc)
```

### Middleware

Middleware functions can be used to modify requests/responses and perform actions before/after handlers:

```go
router.Use(framework.JSONBodyParser)
router.Use(YourCustomMiddleware)
```

### Error Handling

Express.go provides built-in error handling:

```go
router.RegisterErrorHandler(framework.DefaultNotFoundErrorHandler)
router.RegisterErrorHandler(framework.DefaultMethodNotAllowedErrorHandler)
```

### Dependency Injection

The framework includes a dependency injection container with support for different scopes:

- Singleton
- Prototype
- Request

```go
container := framework.NewContainer()
container.Register("service", factory, framework.SingletonScopeStrategy{})
```

### JWT Authentication

Express.go provides a comprehensive JWT authentication middleware that handles token validation and error reporting:

```go
// Create JWT middleware with your secret key
secretKey := []byte("your-secret-key")
router.Use(framework.JWTAuthMiddleware(secretKey))

// Access JWT claims in your handler
router.Handle("/protected", "GET", func(w *framework.ResponseWriter, r *framework.Request) error {
    claims, ok := framework.GetJWTClaimsFromContext(r.Context())
    if !ok {
        return errors.New("JWT claims not found in context")
    }
    
    // Access claims data
    userID, _ := claims["user_id"].(string)
    
    return w.Encode(map[string]string{"message": "Protected resource", "user_id": userID})
})
```

The middleware handles various JWT validation scenarios with specific error types:

- Missing token
- Invalid token format
- Expired token
- Invalid signature
- Malformed token

Each error type has a specific error message that can be used for client-side error handling.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
