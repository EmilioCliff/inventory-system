package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) listReceipts(ctx *gin.Context) {
	receipts, err := server.store.ListReceipts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, receipts)
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

	ctx.JSON(http.StatusOK, receipts)
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

	ctx.JSON(http.StatusOK, receipt)
	return
}
