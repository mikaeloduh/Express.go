package framework

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a new JWT token with the provided claims and signs it with the given secret key
func GenerateJWT(claims jwt.MapClaims, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// GenerateUserJWT creates a JWT for a user with standard claims
func GenerateUserJWT(userID string, expiresIn time.Duration, secretKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,                                 // Subject (user ID)
		"iat": time.Now().Unix(),                      // Issued At
		"exp": time.Now().Add(expiresIn).Unix(),       // Expiration Time
	}
	
	return GenerateJWT(claims, secretKey)
}
