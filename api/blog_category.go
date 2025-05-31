package api

import (
	"net/http"
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type createBlogCategoriesRequest struct {
	Name    string `json:"name" binding:"required"`
	Status 	bool   `json:"status" binding:"required"`
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
