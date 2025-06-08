package api

import (
	"database/sql"
	"net/http"
	db "simplebank/db/sqlc"
	util "simplebank/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username     	string `json:"username" binding:"required"`
	PasswordHash 	string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) loginAuth(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetAccountByUsername(ctx, req.Username)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
        return
    }

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.PasswordHash)); err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

	token, err := util.GenerateAccessToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	arg := db.CreateRefreshTokenParams{
		AccountID: user.ID,
		TokenRefresh: refreshToken,
	}

	if _, err := server.store.CreateRefreshToken(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Respond with success message and token
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"access_token":  	token,
		"refresh_token": 	refreshToken,
		"user": gin.H{ 
			"id":       	user.ID,
			"username": 	user.Username,
			"role":     	user.Role,
		},
	})
}

func (server *Server) refreshAuth(ctx *gin.Context) {
	var req RefreshRequest
	
	// ISSUE 1: You need to bind the request body first
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the refresh token
	claims, err := util.VerifyToken(req.RefreshToken)
	if err != nil {
		// ISSUE 2: This should be 401 Unauthorized, not 500 Internal Server Error
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return // ISSUE 3: Missing return statement
	}

	// Check if it's actually a refresh token
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		// ISSUE 4: You're using 'err' but it might be nil here
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Provided token is not a refresh token"})
		return // ISSUE 5: Missing return statement
	}

	// Extract user ID
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		// ISSUE 6: Again using 'err' which might be nil
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id in refresh token"})
		return // ISSUE 7: Missing return statement
	}

	userID := int64(userIDFloat)

	refreshTokenDetail, err := server.store.GetRefreshToken(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	if refreshTokenDetail.TokenRefresh == req.RefreshToken {
		// Valid refresh token - proceed with cleanup and regeneration
		err := server.store.DeleteRefreshToken(ctx, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		// Token mismatch - reject the request or ignore deletion
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}


	// Generate new access token
	token, err := util.GenerateAccessToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token", "details": err.Error()})
		return
	}

	// Generate new refresh token
	refreshToken, err := util.GenerateRefreshToken(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token", "details": err.Error()})
		return
	}

	arg := db.CreateRefreshTokenParams{
		AccountID: userID,
		TokenRefresh: refreshToken,
	}

	if _, err := server.store.CreateRefreshToken(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}
