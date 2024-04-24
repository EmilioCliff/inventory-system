package api

import (
	"fmt"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/token"
	"github.com/EmilioCliff/inventory-system/worker"
	"github.com/gin-gonic/gin"
)

const (
	PageSize = 2 // return page size to 10
)

type Server struct {
	config          utils.Config
	store           *db.Store
	router          *gin.Engine
	tokenMaker      token.Maker
	emailSender     utils.GmailSender
	taskDistributor worker.TaskDistributor
}

func NewServer(config utils.Config, store *db.Store, emailSender utils.GmailSender, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPaseto(config.TOKEN_SYMMETRY_KEY)
	if err != nil {
		return nil, fmt.Errorf("Couldnt open tokenmaker %w", err)
	}
	server := &Server{
		tokenMaker:      tokenMaker,
		store:           store,
		config:          config,
		emailSender:     emailSender,
		taskDistributor: taskDistributor,
	}

	server.setRoutes()
	return server, nil
}

func (server *Server) setRoutes() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(loggerMiddleware())

	auth := router.Group("/").Use(authMiddleware(server.tokenMaker))
	auth.GET("/users/products/:id", server.getUserProducts)
	auth.GET("/products/", server.listProducts)
	auth.GET("/allproducts/", server.listAllProducts)
	auth.GET("/products/:id", server.getProduct)
	auth.POST("/products/admin/add", server.createProduct)
	auth.DELETE("/products/admin/delete/:id", server.deleteProduct)
	auth.PUT("/products/admin/edit/:id", server.editProduct)

	router.GET("/users/login", server.loginUser)
	router.POST("/reset", server.resetPassword)
	router.POST("/resetit", server.resetIt)
	router.Any("/transaction/:id", server.mpesaCallback)
	auth.GET("/users/:id", server.getUser)
	auth.PUT("/users/:id/edit", server.editUser)
	auth.POST("/users/admin/add", server.createUser)
	auth.DELETE("/users/admin/:id", server.deleteUser)
	auth.GET("/users/admin", server.listUsers)
	auth.PUT("/users/admin/manage/:id", server.manageUser)
	auth.POST("/users/admin/manage/add/:id", server.addClientStock)
	auth.POST("/users/products/admin/add/:id", server.addAdminStock)
	auth.POST("/users/products/sell/:id", server.reduceClientStock)
	auth.GET("/users/invoices/:id", server.getUserInvoices)
	auth.GET("/users/receipts/:id", server.getUserReceipts)

	auth.GET("/search/users", server.searchUser)
	auth.GET("/search/products", server.searchProduct)
	auth.GET("/search/transactions", server.searchTransaction)
	auth.GET("/search/invoices", server.searchInvoice)
	auth.GET("/search/receipts", server.searchReceipt)
	auth.GET("/search/user/invoices", server.searchUserInvoice)
	auth.GET("/search/user/receipts", server.searchUserReceipt)

	auth.GET("/invoices/admin", server.listInvoices)
	auth.GET("/invoices/:id", server.getInvoice)

	auth.GET("/receipts/admin", server.listReceipts)
	auth.GET("/receipts/:id", server.getReceipt)

	auth.GET("/transactions/all", server.allTransactions)
	auth.GET("/transactions/successfull", server.succussfulTransactions)
	auth.GET("/transactions/failed", server.failedTransactions)
	auth.GET("/user/transactions/all/:id", server.getUsersTransactions)
	auth.GET("/user/transactions/successful/:id", server.getUserSuccessfulTransaction)
	auth.GET("/user/transactions/failed/:id", server.getUserFailedTransaction)
	auth.GET("/user/transactions/:id", server.getUserTransaction)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

type PaginationMetadata struct {
	CurrentPage int32 `json:"current_page"`
	TotalData   int32 `json:"total_data"`
	TotalPages  int32 `json:"total_pages"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) GeneratePythonToken(username string) (string, error) {
	return server.tokenMaker.CreateToken(username, server.config.PYTHON_APP_TOKEN_DURATION)
}
