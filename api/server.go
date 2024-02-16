package api

import (
	"fmt"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/token"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     utils.Config
	store      *db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPaseto(config.TOKEN_SYMMETRY_KEY)
	if err != nil {
		return nil, fmt.Errorf("Couldnt open tokenmaker %w", err)
	}
	server := &Server{
		tokenMaker: tokenMaker,
		store:      store,
		config:     config,
	}

	server.setRoutes()
	return server, nil
}

func (server *Server) setRoutes() {
	router := gin.Default()

	auth := router.Group("/").Use(authMiddleware(server.tokenMaker))
	auth.GET("users/products/:id", server.getUserProducts)
	auth.GET("/products/", server.listProducts)
	auth.GET("/products/:id", server.getProduct)
	auth.POST("/products/admin/add", server.createProduct)
	auth.DELETE("/products/admin/delete/:id", server.deleteProduct)
	auth.PUT("/products/admin/edit/:id", server.editProduct)

	router.GET("/users/login", server.loginUser)
	auth.GET("/users/:id", server.getUser)
	auth.PUT("/users/:id/edit", server.editUser)
	auth.POST("/users/admin/add", server.createUser)
	auth.DELETE("/users/admin/:id", server.deleteUser)
	auth.POST("/reset", server.resetPassword)
	auth.POST("/resetit", server.resetIt)
	auth.GET("/users/admin", server.listUsers)
	auth.PUT("/users/admin/manage/:id", server.manageUser)
	auth.POST("/users/admin/manage/add/:id", server.addClientStock)
	auth.POST("/users/products/admin/add/:id", server.addAdminStock)
	auth.POST("/users/products/sell/:id", server.reduceClientStock)
	auth.GET("/users/invoices/:id", server.getUserInvoices)
	auth.GET("/users/receipts/:id", server.getUserReceipts)
	auth.GET("/search/users", server.searchUsers)
	auth.GET("/search/products", server.searchProduct)

	auth.GET("/invoices/admin", server.listInvoices)
	auth.GET("/invoices/:id", server.getInvoice)

	auth.GET("/receipts/admin", server.listReceipts)
	auth.GET("/receipts/:id", server.getReceipt)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) GeneratePythonToken(username string) (string, error) {
	return server.tokenMaker.CreateToken(username, server.config.PYTHON_APP_TOKEN_DURATION)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
