package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	// Alias 'db' for clarity when using db.Account
	util "simplebank/utils"

	"github.com/gin-gonic/gin"
)

// Define the LoginRequest struct directly in this file
// if you don't want a separate requests.go file.
// Or ensure you import it if it's in a different package.
type LoginRequest struct {
	Username     string `json:"username" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken     string `json:"refresh_token" binding:"required"`
}

var activeRefreshTokens = sync.Map{}

// Function for logging in
func Login(c *gin.Context) {
	var req LoginRequest // Bind to the new request struct

	// Check user credentials and generate a JWT token
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
		return
	}

	// Now you access the fields from 'req'
	// Replace this logic with real authentication against your database
	if req.Username == "daruma" && req.PasswordHash == "admin123!" {
		// Assuming a user ID of 1 for 'daruma' for demonstration.
		// In a real app, you'd fetch the actual user ID from the database
		// after successful credential verification.
		userID := int64(1) // Placeholder userID

		refreshToken, err := util.GenerateRefreshToken(userID) // Cast int64 to uint if your GenerateToken expects uint
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
			return
		}

		token, err := util.GenerateAccessToken(userID) // Cast int64 to uint if your GenerateToken expects uint
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"refresh_token": refreshToken,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// In your handlers package (e.g., handlers/login.go or handlers/refresh.go)

// RefreshToken handles requests to obtain a new access token using a valid refresh token.
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// --- Step 1: Check if the refresh token exists in our active store ---
	storedExp, found := activeRefreshTokens.Load(req.RefreshToken)
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found or already used"})
		return
	}

	// --- Step 2: Check if the stored token has expired (even before JWT verification) ---
	if exp, ok := storedExp.(int64); ok {
		if time.Now().Unix() > exp {
			activeRefreshTokens.Delete(req.RefreshToken) // Remove expired token from store
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired in store"})
			return
		}
	} else {
		// This case should ideally not happen if storing correctly
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid expiration data for refresh token"})
		return
	}


	// --- Step 3: Verify the refresh token's signature and claims ---
	claims, err := util.VerifyToken(req.RefreshToken)
	if err != nil {
		// If the refresh token is invalid or expired, return unauthorized
		// The 'details' field here is crucial for debugging!
		activeRefreshTokens.Delete(req.RefreshToken) // Invalidate token if verification fails
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token during JWT verification", "details": err.Error()})
		return
	}

	// --- Step 4: Invalidate the old refresh token immediately after successful verification ---
	activeRefreshTokens.Delete(req.RefreshToken)

	// Extract user ID from claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found or invalid type in refresh token claims"})
		return
	}
	userID := int64(userIDFloat)

	// --- Step 5: Generate NEW access and refresh tokens ---
	newAccessToken, err := util.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating new access token"})
		return
	}

	newRefreshToken, err := util.GenerateRefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating new refresh token"})
		return
	}

	// --- Step 6: Store the NEW refresh token ---
	newRefreshClaims, _ := util.VerifyToken(newRefreshToken) // Get claims for new token
	if exp, ok := newRefreshClaims["exp"].(float64); ok {
		activeRefreshTokens.Store(newRefreshToken, int64(exp)) // Store new token and its expiry
	} else {
		fmt.Println("Warning: Could not get expiration from new refresh token claims.")
	}


	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// Function for registering a new user (for demonstration purposes)
// func Register(c *gin.Context) {
//     // You would typically create a RegisterRequest struct here as well
//     var account db.Account // Still using db.Account for output/database insert

//     if err := c.ShouldBindJSON(&account); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data", "details": err.Error()})
//         return
//     }

//     // Remember to securely hash passwords before storing them
//     // Example: account.PasswordHash = util.HashPassword(account.PasswordHash)

//     // In a real scenario, you'd insert 'account' into your database
//     account.ID = 1 // Just for demonstration purposes
//     c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "username": account.Username})
// }