package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

type createProductRequest struct {
	ProductName string `json:"product_name" binding:"required"`
	UnitPrice   int32  `json:"unit_price"  binding:"required"`
	Packsize    string `json:"packsize"  binding:"required"`
}

func (server *Server) createProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateProductParams{
		ProductName: req.ProductName,
		UnitPrice:   req.UnitPrice,
		Packsize:    req.Packsize,
	}

	product, err := server.store.CreateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
	return
}

type deleteProductRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteProduct(ctx *gin.Context) {
	var req deleteProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteProduct(ctx, int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"item": "deleted successfully"})
	return
}

type editProductUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
type editProductRequest struct {
	ProductName string `json:"product_name" binding:"required"`
	UnitPrice   int32  `json:"unit_price"  binding:"required"`
	Packsize    string `json:"packsize"  binding:"required"`
}

func (server *Server) editProduct(ctx *gin.Context) {
	var uri editProductUri

	rawData, _ := ctx.GetRawData()

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var rawMap map[string]interface{}
	if err := json.Unmarshal(rawData, &rawMap); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	unitPrice, _ := rawMap["unit_price"].(float64)
	productName, _ := rawMap["product_name"].(string)
	packsize, _ := rawMap["packsize"].(string)

	product := db.Product{
		ProductID:   uri.ID,
		UnitPrice:   int32(unitPrice),
		ProductName: productName,
		Packsize:    packsize,
	}
	edited_product, err := server.store.EditStockTx(ctx, db.EditProductParams{
		ProductToEdit: product,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, edited_product.ProductEdited)
	return
}

func (server *Server) listProducts(ctx *gin.Context) {
	list_product, err := server.store.ListProduct(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, list_product)
	return
}

type getProductRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	product, err := server.store.GetProduct(ctx, int64(req.ID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, product)
	return
}
