package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	db "simplebank/db/sqlc"
	util "simplebank/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type createBlogRequest struct {
	Title    string `json:"title" binding:"required"`
	DescriptionArticle string `json:"description_article" binding:"required"`
	CategoryID int64 `json:"category_id" binding:"required"`
	AccountID int64 `json:"account_id" binding:"required"`
	BannerImage string `json:"banner_image"`
}

type updateBlogRequest struct {
	Title    string `json:"title"`
	DescriptionArticle string `json:"description_article"`
	CategoryID int64 `json:"category_id" binding:"required"`
	BannerImage string `json:"banner_image"`
}

type getBlogRequest struct {
	ID    int64 `uri:"id" binding:"required,min=1"`
}

type listBlogsRequest struct {
	PageID    int64 `form:"page_id" binding:"required,min=1"`
	PageSize    int64 `form:"page_size" binding:"required,min=5,max=10"`
}


func (server *Server) createBlog(ctx *gin.Context) {
	var req createBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check Blog Categories
	blogCategories, err := server.store.GetBlogCategories(ctx, req.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check Account Categories
	account, err := server.store.GetAccount(ctx, req.AccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	uploadedPath := util.UploadBase64(ctx, req.BannerImage)

	bannerImagePgText := pgtype.Text{
		String: uploadedPath,
		Valid:  true, 
	}

	arg := db.CreateBlogsParams{
		Title: req.Title,
		DescriptionArticle: req.DescriptionArticle,
		CategoryID: blogCategories.ID,
		AccountID: account.ID,
		BannerImage: bannerImagePgText,
	}

	blog, err := server.store.CreateBlogs(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

func (server *Server) getBlog(ctx *gin.Context) {
	var req getBlogRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	blog, err := server.store.GetBlogs(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, blog)
}

func (server *Server) listBlogs(ctx *gin.Context) {
	var req listBlogsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListBlogssParams{
		Limit: int32(req.PageSize),
		Offset: int32((req.PageID - 1) * req.PageSize),
	}

	blogs, err := server.store.ListBlogss(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, blogs)
}

func (server *Server) updateBlog(ctx *gin.Context) {
	var uriReq getBlogRequest
	var jsonReq updateBlogRequest

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Check Blog Categories
	blogCategories, err := server.store.GetBlogCategories(ctx, jsonReq.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	existingBlog, err := server.store.GetBlogs(ctx, uriReq.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// --- 2. Process the new BannerImage (if provided) ---
	var newBannerImagePgText pgtype.Text

	if jsonReq.BannerImage != "" { // Check if a new base64 image was provided in the request
		uploadedPath := util.UploadBase64(ctx, jsonReq.BannerImage)
		if uploadedPath == "" {
			// Error already handled by util.UploadBase64 via ctx.JSON
			return
		}

		// Create pgtype.Text for the new image path
		newBannerImagePgText = pgtype.Text{
			String: uploadedPath,
			Valid:  true, // It's a valid path now
		}

		// --- 3. Optional: Delete the old image file ---
		// Only delete if there was an existing banner image and it's different from the new one
        // and if new one is successfully uploaded
		if existingBlog.BannerImage.Valid && existingBlog.BannerImage.String != uploadedPath {
			oldImagePath := existingBlog.BannerImage.String
			// Construct the full path to the old file
            // Assuming oldImagePath is something like "uploads/filename.jpg"
            // And your current working directory is where "uploads" folder resides
			fullOldFilePath := oldImagePath // If BannerImage stores full path (e.g. "uploads/image.jpg")

            // On Windows, if oldImagePath was "uploads\filename.jpg"
            // you might need to convert it to OS-specific path first if it's not already
            // If it's already "uploads/filename.jpg" (URL-friendly), os.Remove will usually work
            // but filepath.FromSlash is safer for OS-specific paths.
            fullOldFilePath = filepath.FromSlash(fullOldFilePath)


			// Attempt to remove the file
			err := os.Remove(fullOldFilePath)
			if err != nil {
				// Log the error but don't stop the update.
				// This might happen if the file doesn't exist, permissions issue, etc.
				// It's generally not critical enough to fail the entire blog update.
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old image file: " + err.Error()})
				// You might choose to return here if file deletion is critical
				// but often it's a soft error.
                // For this example, we continue with the update.
                // You might just log this error without sending to user
                // log.Printf("Failed to delete old image %s: %v", fullOldFilePath, err)
			}
		}
	} else {
		// If no new BannerImage was provided in the request, retain the existing one
        // or set it to NULL if you want to allow clearing it.
        // For now, we'll retain the existing one.
		newBannerImagePgText = existingBlog.BannerImage
	}

	arg := db.UpdateBlogsParams{
		ID:      			uriReq.ID,
		Title: 				jsonReq.Title,
		DescriptionArticle: jsonReq.DescriptionArticle,
		CategoryID: 		blogCategories.ID,
		AccountID: 			existingBlog.AccountID,	
		BannerImage:    	newBannerImagePgText,
	}

	err = server.store.UpdateBlogs(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog updated!"})
}

func (server *Server) deleteBlog(ctx *gin.Context) {
	var req getBlogRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	existingBlog, err := server.store.GetBlogs(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// const uploadDir = "./uploads" 

	if existingBlog.BannerImage.Valid {
		oldImagePath := existingBlog.BannerImage.String
		// Ensure oldImagePath is a full path.
		// If existingBlog.BannerImage.String is relative (e.g., "blog_images/image.jpg"),
		// you need to prepend your base upload directory.
		// If it's already an absolute path, this won't change anything.
		fullOldFilePath := filepath.Join(oldImagePath)
		fullOldFilePath = filepath.FromSlash(fullOldFilePath) // Still good for cross-OS compatibility

		fmt.Println("Attempting to delete file:", fullOldFilePath) // Crucial for debugging

		err := os.Remove(fullOldFilePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Warning: Image file '%s' already missing during deletion attempt.\n", fullOldFilePath)
			} else {
				// This is the error you need to examine closely
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old image file: " + err.Error()})
				return
			}
		}
	}

	err = server.store.DeleteBlogs(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog deleted!"})
}