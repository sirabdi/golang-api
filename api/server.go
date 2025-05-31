package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Account
	router.POST("/accounts", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/account/:id", server.updateAccount)
	router.DELETE("/account/:id", server.deleteAccount)

	// Blog
	router.POST("/blog", server.createBlog)
	router.GET("/blog/:id", server.getBlog)
	router.GET("/blogs", server.listBlogs)
	router.PUT("/blog/:id", server.updateBlog)
	router.DELETE("/blog/:id", server.deleteBlog)

	// Blog Categories
	router.POST("/blog-category", server.createBlogCategories)
	router.GET("/blog-categories", server.listBlogCategories)
	router.PUT("/blog-category/:id", server.updateBlogCategories)
	router.DELETE("/blog-category/:id", server.deleteBlogCategories)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}