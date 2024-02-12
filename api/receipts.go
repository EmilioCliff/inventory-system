package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
)

func newReceiptResponse(receipt db.Receipt) (receiptResponse, error) {
	var receiptData []map[string]interface{}
	if receipt.ReceiptData != nil {
		if unerr := json.Unmarshal(receipt.ReceiptData, &receiptData); unerr != nil {
			return receiptResponse{}, unerr
		}
	}

	return receiptResponse{
		ReceiptNumber:       receipt.ReceiptNumber,
		UserreceiptID:       int64(receipt.UserReceiptID),
		UserreceiptUsername: receipt.UserReceiptUsername,
		ReceiptData:         receiptData,
		ReceiptCreateTime:   receipt.CreatedAt,
	}, nil
}

type receiptResponse struct {
	ReceiptNumber       string                   `json:"receipt_number"`
	UserreceiptID       int64                    `json:"user_receipt_id"`
	UserreceiptUsername string                   `json:"user_receipt_username"`
	ReceiptData         []map[string]interface{} `json:"receipt_data"`
	ReceiptCreateTime   time.Time                `json:"receipt_create_time"`
}

func (server *Server) listReceipts(ctx *gin.Context) {
	receipts, err := server.store.ListReceipts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []receiptResponse
	for _, receipt := range receipts {
		updatedReceipt, _ := newReceiptResponse(receipt)
		rsp = append(rsp, updatedReceipt)
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type getUserReceiptsRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) getUserReceipts(ctx *gin.Context) {
	var req getUserReceiptsRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	receipts, err := server.store.GetUserReceiptsByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	var rsp []receiptResponse
	for _, receipt := range receipts {
		updatedReceipt, _ := newReceiptResponse(receipt)
		rsp = append(rsp, updatedReceipt)
	}

	ctx.JSON(http.StatusOK, rsp)
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

	ctx.JSON(http.StatusOK, rsp)
	return
}
