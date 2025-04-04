# JWT Middleware

The Express.go framework provides flexible JWT authentication middleware with support for various customization options.

## Basic Usage

The simplest usage is to create middleware directly with a secret key:

```go
// Create JWT middleware with the provided secret key
secretKey := []byte("your-secret-key")
jwtMiddleware := framework.JWTAuthMiddleware(framework.NewJWTOptions(secretKey))

// Use the middleware in a route
app.GET("/protected", myHandler, jwtMiddleware)
```

## Using Custom Options

For more flexibility, you can use the `Options` struct for customization:

```go
// Create custom options
options := framework.Options{
    // Custom key validation function
    Keyfunc: func(token *jwt.Token) (interface{}, error) {
        // Custom key retrieval logic
        return mySecretKey, nil
    },
    // Custom token retrieval method
    GetHeader: func(r *expressgo.Request) string {
        // For example, get token from custom header or query parameter
        return r.Header.Get("X-API-Token")
    },
    // Add custom claims
    GetClaims: func(r *expressgo.Request) (jwt.MapClaims, bool) {
        return jwt.MapClaims{"custom": "value"}, true
    },
    // Custom context setting
    SetContext: func(ctx context.Context, claims jwt.MapClaims) context.Context {
        // Custom context processing
        return myCustomContext(ctx, claims)
    },
}

// Create middleware with custom options
jwtMiddleware := framework.JWTAuthMiddleware(options)
```

## Options Struct

The `Options` struct defines all configurable options for the JWT middleware:

```go
type Options struct {
    // Keyfunc is a function used to validate JWT signature
    Keyfunc jwt.Keyfunc

    // GetHeader is a function to get the authentication header from the request
    GetHeader func(r *expressgo.Request) string

    // GetClaims is a function to get additional claims, which will be merged with claims from the JWT token
    GetClaims func(r *expressgo.Request) (jwt.MapClaims, bool)

    // SetContext is a function to customize setting JWT claims to the context
    SetContext func(ctx context.Context, claims jwt.MapClaims) context.Context
}
```

## Default Options

If you only provide some options, the remaining ones will use default values:

```go
// Get default options with your secret key
options := framework.NewJWTOptions(secretKey)
// Then customize only what you need
options.GetHeader = myCustomHeaderFunc

// Create middleware with the modified options
jwtMiddleware := framework.JWTAuthMiddleware(options)
```

## Error Handling

The JWT middleware returns specific error types that can be handled differently:

- `ErrorTypeJWTMissing`: Authentication header is missing
- `ErrorTypeJWTInvalidFormat`: Token format is incorrect (missing "Bearer" prefix)
- `ErrorTypeJWTExpired`: Token has expired
- `ErrorTypeJWTInvalidSignature`: Token signature is invalid
- `ErrorTypeJWTInvalid`: Other JWT validation errors
- `ErrorTypeJWTInvalidSigningMethod`: Invalid JWT signing method

You can identify these error types in your error handling middleware and respond accordingly.

## Examples

For a complete example, please refer to `examples/jwt_middleware_example.go`.
