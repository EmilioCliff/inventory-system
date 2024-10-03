package api

import (
	"encoding/json"
	"fmt"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/reports"
	"github.com/EmilioCliff/inventory-system/token"
	"github.com/EmilioCliff/inventory-system/worker"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	PageSize     = 10
	CacheDuraton = time.Second * 30
)

type PaginationMetadata struct {
	CurrentPage int32 `json:"current_page"`
	TotalData   int32 `json:"total_data"`
	TotalPages  int32 `json:"total_pages"`
}

type Server struct {
	config          utils.Config
	store           *db.Store
	router          *gin.Engine
	tokenMaker      token.Maker
	emailSender     utils.GmailSender
	taskDistributor worker.TaskDistributor
	redis           *redis.Client
	reportMaker     reports.ReportMaker // added
}

func NewServer(
	config utils.Config,
	store *db.Store,
	emailSender utils.GmailSender,
	taskDistributor worker.TaskDistributor,
	redis *redis.Client,
) (*Server, error) {
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
		redis:           redis,
		reportMaker:     reports.NewReportMaker(store),
	}

	server.setRoutes()
	return server, nil
}

func (server *Server) setRoutes() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// router.Use(loggerMiddleware())

	auth := router.Group("/").Use(authMiddleware(server.tokenMaker))
	cacheAuth := router.Group("/").Use(authMiddleware(server.tokenMaker), redisCacheMiddleware(server.redis))

	// cors := router.Group("/").Use(CORSmiddleware())
	cacheAuth.GET("/users/products/:id", server.getUserProducts)
	cacheAuth.GET("/products/", server.listProducts)
	cacheAuth.GET("/allproducts/", server.listAllProducts)
	auth.GET("/products/:id", server.getProduct)
	auth.POST("/products/admin/add", server.createProduct)
	auth.DELETE("/products/admin/delete/:id", server.deleteProduct)
	auth.PUT("/products/admin/edit/:id", server.editProduct)
	auth.POST("/products/calculate", server.calculatePrice)

	router.POST("/users/login", server.loginUser)
	router.POST("/reset", server.resetPassword)
	router.POST("/resetit", server.resetIt)
	router.GET("/c2b/list", server.listC2BTransactions)       // added
	auth.POST("/c2b/register", server.registerUrl)            // added
	router.Any("/c2b/complete", server.completeTransaction)   // added
	router.Any("/c2b/validation", server.validateTransaction) // added
	router.Any("/transaction/:id", server.mpesaCallback)
	auth.GET("/search/all", server.searchAll)
	cacheAuth.GET("/users/:id", server.getUser)
	auth.PUT("/users/:id/edit", server.editUser)
	auth.POST("/users/admin/add", server.createUser)
	auth.DELETE("/users/admin/:id", server.deleteUser)
	cacheAuth.GET("/users/admin", server.listUsers)
	auth.PUT("/users/admin/manage/:id", server.manageUser)
	auth.POST("/users/admin/manage/add/:id", server.addClientStock)
	auth.POST("/users/products/admin/add/:id", server.addAdminStock)
	auth.POST("/users/products/sell/:id", server.reduceClientStock)
	cacheAuth.GET("/users/invoices/:id", server.getUserInvoices)
	cacheAuth.GET("/users/receipts/:id", server.getUserReceipts)
	auth.POST("/users/request_stock/:id", server.requestStock)
	auth.POST("/users/admin/reduce_client_stock/:id", server.reduceClientProductByAdmin)

	auth.GET("/search/users", server.searchUser)
	auth.GET("/search/products", server.searchProduct)
	auth.GET("/search/transactions", server.searchTransaction)
	auth.GET("/search/invoices", server.searchInvoice)
	auth.GET("/search/receipts", server.searchReceipt)
	auth.GET("/search/user/invoices", server.searchUserInvoice)

	cacheAuth.GET("/invoices/admin", server.listInvoices)
	cacheAuth.GET("/invoices/:id", server.getInvoice)
	auth.GET("invoice/download/:id", server.downloadInvoice)

	cacheAuth.GET("/receipts/admin", server.listReceipts)
	cacheAuth.GET("/receipts/:id", server.getReceipt)
	auth.GET("receipt/download/:id", server.downloadReceipt)

	cacheAuth.GET("/transactions/all", server.allTransactions)
	cacheAuth.GET("/transactions/successfull", server.succussfulTransactions)
	cacheAuth.GET("/transactions/failed", server.failedTransactions)
	cacheAuth.GET("/user/transactions/all/:id", server.getUsersTransactions)
	cacheAuth.GET("/user/transactions/successful/:id", server.getUserSuccessfulTransaction)
	cacheAuth.GET("/user/transactions/failed/:id", server.getUserFailedTransaction)
	cacheAuth.GET("/user/transactions/:id", server.getUserTransaction)
	auth.GET("/statements/:id", server.downloadStatement)

	auth.POST("/admin/purchase-order", server.createPurchaseOrder)
	auth.DELETE("/admin/purchase-orders/:id", server.deletePurchaseOrders)
	cacheAuth.GET("/admin/purchase-orders", server.listPurchaseOrders)
	auth.GET("/admin/purchase-orders/:id", server.downloadPurchaseOrders)

	cacheAuth.GET("/history/received/:id", server.getUserReceivedHistory)
	cacheAuth.GET("/history/all_received", server.getAllUsersReceivedHistory)
	cacheAuth.GET("/history/sold/:id", server.getUserSoldHistory)
	cacheAuth.GET("/history/debt/:id", server.getUserDebt)
	cacheAuth.GET("/history/all_debt", server.getAllUserDebt)
	cacheAuth.GET("/history/admin", server.adminHistory)
	cacheAuth.GET("/history/test", server.testGroup)

	auth.POST("/admin/users_reports", server.downloadUserReports)
	auth.POST("/admin/admin_reports", server.downloadAdminReports)

	// router.GET("/test", server.test)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setCache(ctx *gin.Context, key string, value any) error {
	rspByte, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return server.redis.Set(ctx, key, rspByte, CacheDuraton).Err()
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) GeneratePythonToken(username string) (string, error) {
	return server.tokenMaker.CreateToken(username, server.config.PYTHON_APP_TOKEN_DURATION)
}
