package api

import (
	"database/sql"
	"net/http"
	util "simplebank/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username     	string `json:"username" binding:"required"`
	PasswordHash 	string `json:"password" binding:"required"`
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

	token, err := util.GenerateJWT(user.Username)// Assuming utils.GenerateJWTToken exists
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)// Assuming utils.GenerateJWTToken exists
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
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