package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
)

type listTransactions struct {
	Data     []transactionResponse `json:"data"`
	Metadata PaginationMetadata    `json:"metadata"`
}

type productDataResponse struct {
	Product  string `json:"product"`
	Quantity int64  `json:"quantity"`
}

type transactionResponse struct {
	TransactionID      string                `json:"transaction_id"`
	TransactionOwner   int32                 `json:"transaction_owner"`
	PhoneNumber        string                `json:"phone_number"`
	Amount             int64                 `json:"amount"`
	MpesaReceiptNumber string                `json:"mpesa_receipt_number"`
	DataSold           []productDataResponse `json:"data_sold"`
	Status             bool                  `json:"status"`
	CreatedAt          time.Time             `json:"created_at"`
}

func newTransactionResponse(transaction db.Transaction, dataReturn []productDataResponse) (transactionResponse, error) {
	// var transactionData map[string][]int
	// if err := json.Unmarshal(transaction.DataSold, &transactionData); err != nil {
	// 	return transactionResponse{}, err
	// }

	// var transactionProducts []productDataResponse
	// for idx, productID := range transactionData["products_id"] {
	// 	product, err := server.store.GetProduct(ctx, int64(productID))
	// 	if err != nil {
	// 		if err == sql.ErrNoRows {
	// 			return transactionResponse{}, err
	// 		}
	// 		return transactionResponse{}, err
	// 	}

	// 	totalAmount := product.UnitPrice * int32(transactionData["quantities"][idx])
	// 	transactionProducts = append(transactionProducts, productDataResponse{
	// 		Product:  product.ProductName,
	// 		Quantity: int64(totalAmount),
	// 	})
	// }

	return transactionResponse{
		TransactionID:      transaction.TransactionID,
		TransactionOwner:   transaction.TransactionUserID,
		PhoneNumber:        transaction.PhoneNumber,
		Amount:             int64(transaction.Amount),
		MpesaReceiptNumber: transaction.MpesaReceiptNumber,
		DataSold:           dataReturn,
		Status:             transaction.Status,
		CreatedAt:          transaction.CreatedAt,
	}, nil
}

type listAllTransactionsRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

func (server *Server) allTransactions(ctx *gin.Context) {
	var req listAllTransactionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transactions, err := server.store.ListTransactions(ctx, db.ListTransactionsParams{
		Limit:  PageSize,
		Offset: (req.PageID - 1) * PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalTransations, err := server.store.CountTransactions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalTransations / PageSize
	if totalTransations%PageSize != 0 {
		totalPages++
	}

	var formatedTransaction []transactionResponse
	for _, transaction := range transactions {
		var transactionData map[string][]int
		if err := json.Unmarshal(transaction.DataSold, &transactionData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var transactionProducts []productDataResponse
		for idx, productID := range transactionData["products_id"] {
			product, err := server.store.GetProduct(ctx, int64(productID))
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalAmount := product.UnitPrice * int32(transactionData["quantities"][idx])
			transactionProducts = append(transactionProducts, productDataResponse{
				Product:  product.ProductName,
				Quantity: int64(totalAmount),
			})
		}

		newTransaction, err := newTransactionResponse(transaction, transactionProducts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		formatedTransaction = append(formatedTransaction, newTransaction)
	}

	rsp := listTransactions{
		Data: formatedTransaction,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalTransations),
		},
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

func (server *Server) succussfulTransactions(ctx *gin.Context) {
	var req listAllTransactionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	succussfulTransactions, err := server.store.SuccessTransactions(ctx, db.SuccessTransactionsParams{
		Limit:  PageSize,
		Offset: (req.PageID - 1) * PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalTransactions, err := server.store.CountSuccessfulTransactions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalTransactions / PageSize
	if totalPages%PageSize != 0 {
		totalPages++
	}

	var formatedTransaction []transactionResponse
	for _, transaction := range succussfulTransactions {
		var transactionData map[string][]int
		if err := json.Unmarshal(transaction.DataSold, &transactionData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var transactionProducts []productDataResponse
		for idx, productID := range transactionData["products_id"] {
			product, err := server.store.GetProduct(ctx, int64(productID))
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalAmount := product.UnitPrice * int32(transactionData["quantities"][idx])
			transactionProducts = append(transactionProducts, productDataResponse{
				Product:  product.ProductName,
				Quantity: int64(totalAmount),
			})
		}

		newTransaction, err := newTransactionResponse(transaction, transactionProducts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		formatedTransaction = append(formatedTransaction, newTransaction)
	}

	rsp := listTransactions{
		Data: formatedTransaction,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalTransactions),
		},
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

func (server *Server) failedTransactions(ctx *gin.Context) {
	var req listAllTransactionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	failedTransactions, err := server.store.FailedTransactions(ctx, db.FailedTransactionsParams{
		Limit:  PageSize,
		Offset: (req.PageID - 1) * PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalTransactions, err := server.store.CountFailedTransactions(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalTransactions / PageSize
	if totalPages%PageSize != 0 {
		totalPages++
	}

	var formatedTransaction []transactionResponse
	for _, transaction := range failedTransactions {
		var transactionData map[string][]int
		if err := json.Unmarshal(transaction.DataSold, &transactionData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var transactionProducts []productDataResponse
		for idx, productID := range transactionData["products_id"] {
			product, err := server.store.GetProduct(ctx, int64(productID))
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalAmount := product.UnitPrice * int32(transactionData["quantities"][idx])
			transactionProducts = append(transactionProducts, productDataResponse{
				Product:  product.ProductName,
				Quantity: int64(totalAmount),
			})
		}

		newTransaction, err := newTransactionResponse(transaction, transactionProducts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		formatedTransaction = append(formatedTransaction, newTransaction)
	}

	rsp := listTransactions{
		Data: formatedTransaction,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalTransactions),
		},
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type listUserTransactionsRequest struct {
	UserID int64 `uri:"id" binding:"required"`
}

func (server *Server) getUsersTransactions(ctx *gin.Context) {
	var uri listUserTransactionsRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req listAllTransactionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userTransactions, err := server.store.AllUserTransactions(ctx, db.AllUserTransactionsParams{
		Limit:             PageSize,
		Offset:            (req.PageID - 1) * PageSize,
		TransactionUserID: int32(uri.UserID),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalTransactions, err := server.store.CountAllUserTransactions(ctx, int32(uri.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalTransactions / PageSize
	if totalPages%PageSize != 0 {
		totalPages++
	}

	var formatedTransaction []transactionResponse
	for _, transaction := range userTransactions {
		var transactionData map[string][]int
		if err := json.Unmarshal(transaction.DataSold, &transactionData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var transactionProducts []productDataResponse
		for idx, productID := range transactionData["products_id"] {
			product, err := server.store.GetProduct(ctx, int64(productID))
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalAmount := product.UnitPrice * int32(transactionData["quantities"][idx])
			transactionProducts = append(transactionProducts, productDataResponse{
				Product:  product.ProductName,
				Quantity: int64(totalAmount),
			})
		}

		newTransaction, err := newTransactionResponse(transaction, transactionProducts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		formatedTransaction = append(formatedTransaction, newTransaction)
	}

	rsp := listTransactions{
		Data: formatedTransaction,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalTransactions),
		},
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

func (server *Server) getUserSuccessfulTransaction(ctx *gin.Context) {
	var uri listUserTransactionsRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req listAllTransactionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	successfulUserTransactions, err := server.store.SuccessUserTransactions(ctx, db.SuccessUserTransactionsParams{
		Limit:             PageSize,
		Offset:            (req.PageID - 1) * PageSize,
		TransactionUserID: int32(uri.UserID),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalTransactions, err := server.store.CountSuccessfulUserTransactions(ctx, int32(uri.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalTransactions / PageSize
	if totalPages%PageSize != 0 {
		totalPages++
	}

	var formatedTransaction []transactionResponse
	for _, transaction := range successfulUserTransactions {
		var transactionData map[string][]int
		if err := json.Unmarshal(transaction.DataSold, &transactionData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var transactionProducts []productDataResponse
		for idx, productID := range transactionData["products_id"] {
			product, err := server.store.GetProduct(ctx, int64(productID))
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalAmount := product.UnitPrice * int32(transactionData["quantities"][idx])
			transactionProducts = append(transactionProducts, productDataResponse{
				Product:  product.ProductName,
				Quantity: int64(totalAmount),
			})
		}

		newTransaction, err := newTransactionResponse(transaction, transactionProducts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		formatedTransaction = append(formatedTransaction, newTransaction)
	}

	rsp := listTransactions{
		Data: formatedTransaction,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalTransactions),
		},
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

func (server *Server) getUserFailedTransaction(ctx *gin.Context) {
	var uri listUserTransactionsRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req listAllTransactionsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	failedUserTransactions, err := server.store.FailedUserTransactions(ctx, db.FailedUserTransactionsParams{
		Limit:             PageSize,
		Offset:            (req.PageID - 1) * PageSize,
		TransactionUserID: int32(uri.UserID),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalTransactions, err := server.store.CountFailedUserTransactions(ctx, int32(uri.UserID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalTransactions / PageSize
	if totalPages%PageSize != 0 {
		totalPages++
	}

	var formatedTransaction []transactionResponse
	for _, transaction := range failedUserTransactions {
		var transactionData map[string][]int
		if err := json.Unmarshal(transaction.DataSold, &transactionData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var transactionProducts []productDataResponse
		for idx, productID := range transactionData["products_id"] {
			product, err := server.store.GetProduct(ctx, int64(productID))
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalAmount := product.UnitPrice * int32(transactionData["quantities"][idx])
			transactionProducts = append(transactionProducts, productDataResponse{
				Product:  product.ProductName,
				Quantity: int64(totalAmount),
			})
		}

		newTransaction, err := newTransactionResponse(transaction, transactionProducts)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		formatedTransaction = append(formatedTransaction, newTransaction)
	}

	rsp := listTransactions{
		Data: formatedTransaction,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalTransactions),
		},
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type getUserTransactionRequest struct {
	TrasanctionID string `uri:"id" binding:"required"`
}

func (server *Server) getUserTransaction(ctx *gin.Context) {
	var req getUserTransactionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transactions, err := server.store.GetTransaction(ctx, req.TrasanctionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var transactionData map[string][]int
	if err := json.Unmarshal(transactions.DataSold, &transactionData); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var transactionProducts []productDataResponse
	for idx, productID := range transactionData["products_id"] {
		product, err := server.store.GetProduct(ctx, int64(productID))
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		totalAmount := product.UnitPrice * int32(transactionData["quantities"][idx])
		transactionProducts = append(transactionProducts, productDataResponse{
			Product:  product.ProductName,
			Quantity: int64(totalAmount),
		})
	}

	rsp, err := newTransactionResponse(transactions, transactionProducts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}
