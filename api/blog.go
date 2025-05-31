package api

import (
	"net/http"
	db "simplebank/db/sqlc"

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


func (server *Server) createBlog(ctx *gin.Context) {
	var req createBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateBlogsParams{
		Title: req.Title,
		DescriptionArticle: req.DescriptionArticle,
		CategoryID: req.CategoryID,
		AccountID: req.AccountID,
		BannerImage: pgtype.Text{String: req.BannerImage, Valid: req.BannerImage != ""},
	}

	blog, err := server.store.CreateBlogs(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, blog)
}
