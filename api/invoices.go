package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
)

func newInvoiceResponse(invoice db.Invoice) (invoiceResponse, error) {
	var invoiceData []map[string]interface{}
	if invoice.InvoiceData != nil {
		if unerr := json.Unmarshal(invoice.InvoiceData, &invoiceData); unerr != nil {
			return invoiceResponse{}, unerr
		}
	}

	return invoiceResponse{
		InvoiceNumber:       invoice.InvoiceNumber,
		UserInvoiceID:       int64(invoice.UserInvoiceID),
		UserInvoiceUsername: invoice.UserInvoiceUsername,
		InvoiceData:         invoiceData,
		InvoiceCreateTime:   invoice.CreatedAt,
	}, nil
}

type invoiceResponse struct {
	InvoiceNumber       string                   `json:"invoice_number"`
	UserInvoiceID       int64                    `json:"user_invoice_id"`
	UserInvoiceUsername string                   `json:"user_invoice_username"`
	InvoiceData         []map[string]interface{} `json:"invoice_data"`
	InvoiceCreateTime   time.Time                `json:"invoice_create_time"`
}

func (server *Server) listInvoices(ctx *gin.Context) {
	invoices, err := server.store.ListInvoices(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []invoiceResponse
	for _, invoice := range invoices {
		updatedInvoice, _ := newInvoiceResponse(invoice)
		rsp = append(rsp, updatedInvoice)
	}
	ctx.JSON(http.StatusOK, rsp)
	return
}

type getUserInvoicesRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (server *Server) getUserInvoices(ctx *gin.Context) {
	var req getUserInvoicesRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	invoices, err := server.store.GetUserInvoicesByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	var rsp []invoiceResponse
	for _, invoice := range invoices {
		updatedInvoice, _ := newInvoiceResponse(invoice)
		rsp = append(rsp, updatedInvoice)
	}

	ctx.JSON(http.StatusOK, rsp)
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

	ctx.JSON(http.StatusOK, rsp)
	return
}
