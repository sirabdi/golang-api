package util

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("secretpassword")

// GenerateToken generates a JWT token with the user ID as part of the claims
func GenerateAccessToken(userID int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	// Access token valid for 60 seconds as requested
	claims["exp"] = time.Now().Add(time.Second * 60).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// GenerateRefreshToken generates a JWT refresh token.
func GenerateRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	// Refresh token valid for a longer period, e.g., 7 days (or 60 seconds for this example)
	// For this example, we'll also set it to 60 seconds as per your request,
	// but in a real app, this would be much longer than the access token.
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // Example: time.Hour * 24 * 7 for 7 days

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// VerifyToken verifies a token JWT validate 
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
    // Parse the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Check the signing method
        if token.Method != jwt.SigningMethodHS256 {
            return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
        }

        return secretKey, nil
    })

    // Check for errors
    if err != nil {
        return nil, err
    }

    // Validate the token
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("Invalid token")
}