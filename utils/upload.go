package util

import (
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func UploadBase64(ctx *gin.Context, n string) pgtype.Text {
	base64Data := n
	if commaIdx := strings.Index(base64Data, ","); commaIdx != -1 {
		base64Data = base64Data[commaIdx+1:]
	}

	// Decode
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 image"})
	}

	// Save to local folder
	fileName := time.Now().Format("20060102150405") + ".jpg" // you could generate UUID instead
	savePath := filepath.Join("uploads", fileName)

	err = os.WriteFile(savePath, imageData, 0644)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save image"})
	}

	return pgtype.Text{String: savePath, Valid: true}
}