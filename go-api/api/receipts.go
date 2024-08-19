package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

func newReceiptResponse(receipt db.Receipt) (receiptResponse, error) {
	var receiptData []map[string]interface{}
	if receipt.ReceiptData != nil {
		if unerr := json.Unmarshal(receipt.ReceiptData, &receiptData); unerr != nil {
			return receiptResponse{}, unerr
		}
	}

	mytime := receipt.CreatedAt.Format("02-January-2006")
	return receiptResponse{
		ReceiptID:           receipt.ReceiptID,
		ReceiptNumber:       receipt.ReceiptNumber,
		UserreceiptID:       int64(receipt.UserReceiptID),
		UserreceiptUsername: receipt.UserReceiptUsername,
		ReceiptData:         receiptData,
		ReceiptCreateTime:   mytime,
		PaymentMethod:       receipt.PaymentMethod,
	}, nil
}

type receiptResponse struct {
	ReceiptID           int64                    `json:"receipt_id"`
	ReceiptNumber       string                   `json:"receipt_number"`
	UserreceiptID       int64                    `json:"user_receipt_id"`
	UserreceiptUsername string                   `json:"user_receipt_username"`
	ReceiptData         []map[string]interface{} `json:"receipt_data"`
	PaymentMethod       string                   `json:"payment_method"`
	ReceiptCreateTime   string                   `json:"receipt_create_time"`
}

type listReceiptRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

type listReceiptResponse struct {
	Data     []receiptResponse  `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

func (server *Server) listReceipts(ctx *gin.Context) {
	var req listReceiptRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("%v:%v", ctx.Request.URL.Path, req.PageID)
	cacheData, err := server.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Info().Msgf("cached hit for: %v", cacheKey)
		ctx.Data(http.StatusOK, "application/json", cacheData)
		return
	}

	receipts, err := server.store.ListReceipts(ctx, db.ListReceiptsParams{
		Limit:  PageSize,
		Offset: (req.PageID - 1) * PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// return to receipt details
	var rsp []receiptResponse
	for _, receipt := range receipts {
		updatedReceipt, _ := newReceiptResponse(receipt)
		transaction, err := server.store.GetTransaction(ctx, receipt.ReceiptNumber)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		additionalData1 := []map[string]interface{}{
			{
				"mpesa_ref":    transaction.MpesaReceiptNumber,
				"amount":       transaction.Amount,
				"phone_number": transaction.PhoneNumber,
			},
		}

		// updatedReceipt.ReceiptData = append(updatedReceipt.ReceiptData, additionalData1)
		updatedReceipt.ReceiptData = append(additionalData1, updatedReceipt.ReceiptData...)
		// updatedReceipt := receiptResponse{
		// 	ReceiptID:           receipt.ReceiptID,
		// 	ReceiptNumber:       receipt.ReceiptNumber,
		// 	UserreceiptID:       int64(receipt.UserReceiptID),
		// 	UserreceiptUsername: receipt.UserReceiptUsername,
		// 	ReceiptData: []map[string]interface{}{
		// 		{
		// 			"user_contact": "dummy_data",
		// 		},
		// 		{
		// 			"mpesa_ref":    transaction.MpesaReceiptNumber,
		// 			"amount":       transaction.Amount,
		// 			"phone_number": transaction.PhoneNumber,
		// 		},
		// 	},
		// 	ReceiptCreateTime: mytime,
		// }
		rsp = append(rsp, updatedReceipt)
	}

	totalReceipt, err := server.store.CountReceipts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalReceipt / int64(PageSize)
	if totalReceipt%int64(PageSize) != 0 {
		totalPages++
	}

	newRsp := listReceiptResponse{
		Data: rsp,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalReceipt),
		},
	}

	err = server.setCache(ctx, cacheKey, newRsp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newRsp)
	return
}

type getUserReceiptsRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

type getUserReceiptsFormRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

type getUserReceiptsResponse struct {
	Data     []receiptResponse  `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

func (server *Server) getUserReceipts(ctx *gin.Context) {
	var req getUserReceiptsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var page getUserReceiptsFormRequest
	if err := ctx.ShouldBindQuery(&page); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	cacheKey := fmt.Sprintf("%v/%v:%v", ctx.Request.URL.Path, req.ID, page.PageID)
	cacheData, err := server.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Info().Msgf("cached hit for: %v", cacheKey)
		ctx.Data(http.StatusOK, "application/json", cacheData)
		return
	}

	receipts, err := server.store.GetUserReceiptsByID(ctx, db.GetUserReceiptsByIDParams{
		UserReceiptID: req.ID,
		Limit:         PageSize,
		Offset:        (page.PageID - 1) * PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	var rsp []receiptResponse
	for _, receipt := range receipts {
		updatedReceipt, _ := newReceiptResponse(receipt)
		rsp = append(rsp, updatedReceipt)
	}

	totalReceipt, err := server.store.CountUserReceiptsByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	totalPages := totalReceipt / PageSize
	if totalReceipt%PageSize != 0 {
		totalPages++
	}

	newRsp := getUserReceiptsResponse{
		Data: rsp,
		Metadata: PaginationMetadata{
			CurrentPage: page.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalReceipt),
		},
	}

	err = server.setCache(ctx, cacheKey, newRsp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newRsp)
	return
}

type getReceiptRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) getReceipt(ctx *gin.Context) {
	var req getReceiptRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	receipt, err := server.store.GetReceipt(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, _ := newReceiptResponse(receipt)

	err = server.setCache(ctx, GetReceipt+fmt.Sprintf("%v", req.ID), rsp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type downloadReceiptRequest struct {
	ID string `uri:"id" binding:"required"`
}

type dowloadReceiptResponse struct {
	ReceiptPdf []byte `json:"receipt_pdf"`
}

func (server *Server) downloadReceipt(ctx *gin.Context) {
	var req getReceiptRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	receipt, err := server.store.GetReceipt(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := dowloadReceiptResponse{
		ReceiptPdf: receipt.ReceiptPdf,
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type searchReceipt struct {
	SearchWord string `form:"search_word" binding:"required"`
}

func (server *Server) searchReceipt(ctx *gin.Context) {
	var req searchReceipt

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var pgQuery pgtype.Text
	if err := pgQuery.Scan(req.SearchWord); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rst, err := server.store.SearchILikeReceipts(ctx, pgQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rst)
	return
}

type downloadStatementRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) downloadStatement(ctx *gin.Context) {
	var uri downloadStatementRequest
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transactions, err := server.store.AllUserTransactionsNoLimit(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(transactions) == 0 {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "user has no successful transactions"})
		return
	}

	var statementData []map[string]interface{}
	for _, transaction := range transactions {
		data := map[string]interface{}{
			"receipt_number": transaction.TransactionID,
			"mpesa_number":   transaction.MpesaReceiptNumber,
			"amount":         transaction.Amount,
			"created_at":     transaction.CreatedAt.Format("2006-01-02 15:04"),
		}

		statementData = append(statementData, data)
	}

	user, err := server.store.GetUser(ctx, int64(uri.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// send data and user to generate statement pdf

	userDetails := map[string]string{
		"username": user.Username,
		"phone":    user.PhoneNumber,
		"email":    user.Email,
		"date":     time.Now().Format("2006-01-02"),
	}

	result := make(chan struct {
		pdfBytes []byte
		err      error
	}, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func([]map[string]interface{}, map[string]string) {
		defer wg.Done()
		pdfBytes, err := utils.GenerateStatement(statementData, userDetails)
		result <- struct {
			pdfBytes []byte
			err      error
		}{pdfBytes: pdfBytes, err: err}
	}(statementData, userDetails)

	wg.Wait()
	close(result)

	data := <-result
	if data.err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(data.err))
		return
	}

	rsp := map[string]interface{}{
		"statement_pdf": data.pdfBytes,
		"data":          statementData,
		"user":          userDetails,
	}

	ctx.JSON(http.StatusOK, rsp)
}

// func (server *Server) searchUserReceipt(ctx *gin.Context) {
// 	var req searchReceipt

// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	var pgQuery pgtype.Text
// 	if err := pgQuery.Scan(req.SearchWord); err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	rst, err := server.store.SearchUserReceipts(ctx, pgQuery)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, rst)
// 	return
// }
