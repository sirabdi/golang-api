package api

import (
	"database/sql"
	"net/http"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type createBlogCategoriesRequest struct {
	Name    string `json:"name" binding:"required"`
	Status 	bool   `json:"status" binding:"required"`
}

type updateBlogCategoriesRequest struct {
	Name    string `json:"name"`
	Status 	bool `json:"status"`
}

type listBlogsCategoriesRequest struct {
	PageID    int64 `form:"page_id" binding:"required,min=1"`
	PageSize    int64 `form:"page_size" binding:"required,min=5,max=10"`
}

type getBlogsCategoriesRequest struct {
	ID    int64 `uri:"id" binding:"required,min=1"`
}


func (server *Server) createBlogCategories(ctx *gin.Context) {
	var req createBlogCategoriesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateBlogCategoriesParams{
		Name: req.Name,
		Status: pgtype.Bool{Bool: req.Status, Valid: true},
	}

	blog, err := server.store.CreateBlogCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, blog)
}

func (server *Server) listBlogCategories(ctx *gin.Context) {
	var req listBlogsCategoriesRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListBlogCategoriesParams{
		Limit: int32(req.PageSize),
		Offset: int32((req.PageID - 1) * req.PageSize),
	}

	blogCategories, err := server.store.ListBlogCategories(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	ctx.JSON(http.StatusOK, blogCategories)
}

func (server *Server) updateBlogCategories(ctx *gin.Context) {
	var uriReq getBlogsCategoriesRequest
	var jsonReq updateBlogCategoriesRequest

	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindJSON(&jsonReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateBlogCategoriesParams{
		ID:      			uriReq.ID,
		Name: 				jsonReq.Name,
		Status: 			pgtype.Bool{ Bool: jsonReq.Status, Valid: true},
	}

	err := server.store.UpdateBlogCategories(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog Category updated!"})
}

func (server *Server) deleteBlogCategories(ctx *gin.Context) {
	var req getBlogRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	existingBlogCategory, err := server.store.GetBlogCategories(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.DeleteBlogCategories(ctx, existingBlogCategory.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Blog Category deleted!"})
}