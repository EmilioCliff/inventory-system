package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

func newInvoiceResponse(invoice db.Invoice) (invoiceResponse, error) {
	var invoiceData []map[string]interface{}
	if invoice.InvoiceData != nil {
		if unerr := json.Unmarshal(invoice.InvoiceData, &invoiceData); unerr != nil {
			return invoiceResponse{}, unerr
		}
	}
	mytime := invoice.CreatedAt.Format("02-January-2006")
	return invoiceResponse{
		InvoiceID:           invoice.InvoiceID,
		InvoiceNumber:       invoice.InvoiceNumber,
		UserInvoiceID:       int64(invoice.UserInvoiceID),
		UserInvoiceUsername: invoice.UserInvoiceUsername,
		InvoiceData:         invoiceData,
		InvoiceCreateTime:   mytime,
	}, nil
}

type invoiceResponse struct {
	InvoiceID           int64                    `json:"invoice_id"`
	InvoiceNumber       string                   `json:"invoice_number"`
	UserInvoiceID       int64                    `json:"user_invoice_id"`
	UserInvoiceUsername string                   `json:"user_invoice_username"`
	InvoiceData         []map[string]interface{} `json:"invoice_data"`
	InvoiceCreateTime   string                   `json:"invoice_create_time"`
}

type listInvoiceRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

type listInvoiceResponse struct {
	Data     []invoiceResponse  `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

func (server *Server) listInvoices(ctx *gin.Context) {
	var req listInvoiceRequest
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

	invoices, err := server.store.ListInvoices(ctx, db.ListInvoicesParams{
		Limit:  int32(PageSize),
		Offset: int32((req.PageID - 1) * PageSize),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var returnInvoices []invoiceResponse
	for _, invoice := range invoices {
		updatedInvoice, _ := newInvoiceResponse(invoice)
		returnInvoices = append(returnInvoices, updatedInvoice)
	}

	totalInvoice, err := server.store.CountInvoices(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalInvoice / int64(PageSize)
	if totalInvoice%int64(PageSize) != 0 {
		totalPages++
	}

	rsp := listInvoiceResponse{
		Data: returnInvoices,
		Metadata: PaginationMetadata{
			CurrentPage: int32(req.PageID),
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalInvoice),
		},
	}

	err = server.setCache(ctx, cacheKey, rsp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type getUserInvoicesRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

type getUserInvoicesFormRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

type getUserInvoicesResponse struct {
	Data     []invoiceResponse  `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

func (server *Server) getUserInvoices(ctx *gin.Context) {
	var req getUserInvoicesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var page getUserInvoicesFormRequest
	if err := ctx.ShouldBindQuery(&page); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// path/id:pageID
	cacheKey := fmt.Sprintf("%v/%v:%v", ctx.Request.URL.Path, req.ID, page.PageID)
	cacheData, err := server.redis.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Info().Msgf("cached hit for: %v", cacheKey)
		ctx.Data(http.StatusOK, "application/json", cacheData)
		return
	}

	invoices, err := server.store.GetUserInvoicesByID(ctx, db.GetUserInvoicesByIDParams{
		Limit:         PageSize,
		Offset:        (page.PageID - 1) * PageSize,
		UserInvoiceID: req.ID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	var rsp []invoiceResponse
	for _, invoice := range invoices {
		updatedInvoice, _ := newInvoiceResponse(invoice)
		rsp = append(rsp, updatedInvoice)
	}

	totalInvoice, err := server.store.CountUserInvoicesByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	totalPages := totalInvoice / PageSize
	if totalInvoice%PageSize != 0 {
		totalPages++
	}

	newRsp := getUserInvoicesResponse{
		Data: rsp,
		Metadata: PaginationMetadata{
			CurrentPage: page.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalInvoice),
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

type getInvoiceRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) getInvoice(ctx *gin.Context) {
	var req getInvoiceRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	invoice, err := server.store.GetInvoice(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, _ := newInvoiceResponse(invoice)

	err = server.setCache(ctx, GetInvoice+fmt.Sprintf("%v", req.ID), rsp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type downloadInvoiceRequest struct {
	ID string `uri:"id" binding:"required"`
}

type downloadInvoiceResponse struct {
	InvoicePdf []byte `json:"invoice_pdf"`
}

func (server *Server) downloadInvoice(ctx *gin.Context) {
	var req getInvoiceRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	invoice, err := server.store.GetInvoice(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := downloadInvoiceResponse{
		InvoicePdf: invoice.InvoicePdf,
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type searchInvoice struct {
	SearchWord string `form:"search_word" binding:"required"`
}

func (server *Server) searchInvoice(ctx *gin.Context) {
	var req searchInvoice

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var pgQuery pgtype.Text
	if err := pgQuery.Scan(req.SearchWord); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rst, err := server.store.SearchILikeInvoices(ctx, pgQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rst)
	return
}

func (server *Server) searchUserInvoice(ctx *gin.Context) {
	var req searchInvoice

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var pgQuery pgtype.Text
	if err := pgQuery.Scan(req.SearchWord); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rst, err := server.store.SearchUserInvoices(ctx, pgQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rst)
	return
}
