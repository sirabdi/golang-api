package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("secretpassword")

// GenerateToken generates a JWT token with the user ID as part of the claims
func GenerateAccessToken(userID int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["type"] = "access"
	// Access token valid for 60 seconds as requested
	claims["exp"] = time.Now().Add(time.Second * 60).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// GenerateRefreshToken generates a JWT refresh token.
func GenerateRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["type"] = "refresh"
	// Refresh token valid for a longer period, e.g., 7 days (or 60 seconds for this example)
	// For this example, we'll also set it to 60 seconds as per your request,
	// but in a real app, this would be much longer than the access token.
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // Example: time.Hour * 24 * 7 for 7 days

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// VerifyToken verifies a token JWT validate 
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	fmt.Println("tokenString ", tokenString)

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte("secretpassword"), nil
    })

	fmt.Println("token ", token)
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}