package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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
		ReceiptPdf:          receipt.ReceiptPdf,
	}, nil
}

type receiptResponse struct {
	ReceiptID           int64                    `json:"receipt_id"`
	ReceiptNumber       string                   `json:"receipt_number"`
	UserreceiptID       int64                    `json:"user_receipt_id"`
	UserreceiptUsername string                   `json:"user_receipt_username"`
	ReceiptData         []map[string]interface{} `json:"receipt_data"`
	ReceiptCreateTime   string                   `json:"receipt_create_time"`
	ReceiptPdf          []byte                   `json:"receipt_pdf"`
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

	receipts, err := server.store.ListReceipts(ctx, db.ListReceiptsParams{
		Limit:  PageSize,
		Offset: (req.PageID - 1) * PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []receiptResponse
	for _, receipt := range receipts {
		updatedReceipt, _ := newReceiptResponse(receipt)
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

func (server *Server) searchUserReceipt(ctx *gin.Context) {
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

	rst, err := server.store.SearchUserReceipts(ctx, pgQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rst)
	return
}
