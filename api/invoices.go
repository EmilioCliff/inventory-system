package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) listInvoices(ctx *gin.Context) {
	invoices, err := server.store.ListInvoices(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, invoices)
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

	ctx.JSON(http.StatusOK, invoices)
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

	ctx.JSON(http.StatusOK, invoice)
	return
}
