package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

	admin, err := server.store.GetUserForUpdate(ctx, 1)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = server.store.AddAdminStockTx(ctx, db.AddAdminStockParams{
		Admin:       admin,
		ProducToAdd: product,
		Amount:      0,
	})

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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

type listProductsRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

type listProductsResponse struct {
	Data     []db.Product       `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

func (server *Server) listProducts(ctx *gin.Context) {
	var req listProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	list_product, err := server.store.ListProduct(ctx, db.ListProductParams{
		Limit:  int32(PageSize),
		Offset: int32((req.PageID - 1) * PageSize),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalProduct, err := server.store.CountProducts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalProduct / int64(PageSize)
	if totalProduct%int64(PageSize) != 0 {
		totalPages++
	}

	rsp := listProductsResponse{
		Data: list_product,
		Metadata: PaginationMetadata{
			TotalPages:  int32(totalPages),
			CurrentPage: int32(req.PageID),
			TotalData:   int32(totalProduct),
		},
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

func (server *Server) listAllProducts(ctx *gin.Context) {

	list_product, err := server.store.ListAllProduct(ctx)
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

type getUserProductsRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

func (server *Server) getUserProducts(ctx *gin.Context) {
	var uri getUserProductsRequest

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := []map[string]interface{}{}

	if user.Stock != nil {
		if unerr := json.Unmarshal(user.Stock, &rsp); unerr != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, rsp)
	return
}

type searchProduct struct {
	SearchWord string `form:"search_word" binding:"required"`
}

func (server *Server) searchProduct(ctx *gin.Context) {
	var req searchProduct

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var pgQuery pgtype.Text
	if err := pgQuery.Scan(req.SearchWord); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rst, err := server.store.SearchILikeProducts(ctx, pgQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rst)
	return
}

type listUserSearchResponse struct {
	Data     []userResponse     `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

type listProductsSearchResponse struct {
	Data     []db.Product       `json:"data"`
	Metadata PaginationMetadata `json:"metadata"`
}

type listReceiptSearchResponse struct {
	Data      []receiptResponse  `json:"data"`
	QueryWord string             `json:"query_word"`
	Metadata  PaginationMetadata `json:"metadata"`
}

type listInvoiceSearchResponse struct {
	Data      []invoiceResponse  `json:"data"`
	QueryWord string             `json:"query_word"`
	Metadata  PaginationMetadata `json:"metadata"`
}

type searchAllRequest struct {
	SearchQuery   string `form:"search_query" binding:"required"`
	PageID        int32  `form:"page_id" binding:"required"`
	SearchContext string `form:"search_context" binding:"required"`
}

func (server *Server) searchAll(ctx *gin.Context) {
	var req searchAllRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	context := req.SearchContext
	switch context {
	case "users":
		var pgQuery pgtype.Text
		if err := pgQuery.Scan(req.SearchQuery); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		results, err := server.store.SearchILikeUsers(ctx, pgQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var rsp []userResponse
		for _, userFound := range results {
			user, err := server.store.GetUserByUsename(ctx, userFound)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			resp, _ := newUserResponse(user)
			rsp = append(rsp, resp)
		}

		rspSend := listUserSearchResponse{
			Data: rsp,
			Metadata: PaginationMetadata{
				CurrentPage: 1,
				TotalPages:  1,
				TotalData:   int32(len(rsp)),
			},
		}

		ctx.JSON(http.StatusOK, rspSend)

	case "products":
		var pgQuery pgtype.Text
		if err := pgQuery.Scan(req.SearchQuery); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		results, err := server.store.SearchILikeProducts(ctx, pgQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var productsFound []db.Product
		for _, result := range results {
			product, err := server.store.GetProductByProductName(ctx, result)
			if err != nil {
				if err == sql.ErrNoRows {
					ctx.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			productsFound = append(productsFound, product)
		}

		rsp := listProductsSearchResponse{
			Data: productsFound,
			Metadata: PaginationMetadata{
				TotalPages:  1,
				CurrentPage: 1,
				TotalData:   int32(len(productsFound)),
			},
		}

		ctx.JSON(http.StatusOK, rsp)

	case "receipts":
		var pgQuery pgtype.Text
		if err := pgQuery.Scan(req.SearchQuery); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		results, err := server.store.SearchUserReceipts(ctx, pgQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		var totalReceipts int32
		processedUsernames := make(map[string]struct{})
		var rsp []receiptResponse
		for _, result := range results {
			if _, ok := processedUsernames[result]; ok {
				continue
			}
			receipts, err := server.store.GetUserReceiptsByUsername(ctx, db.GetUserReceiptsByUsernameParams{
				UserReceiptUsername: result,
				Limit:               PageSize,
				Offset:              (req.PageID - 1) * PageSize,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalReceipt, err := server.store.CountUserReceiptsByUsername(ctx, result)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			totalReceipts = totalReceipts + int32(totalReceipt)

			for _, receipt := range receipts {
				updatedReceipt, _ := newReceiptResponse(receipt)
				rsp = append(rsp, updatedReceipt)
			}

			processedUsernames[result] = struct{}{}
		}

		totalPages := totalReceipts / PageSize
		if totalReceipts%PageSize != 0 {
			totalPages++
		}

		newRsp := listReceiptSearchResponse{
			Data: rsp,
			Metadata: PaginationMetadata{
				CurrentPage: req.PageID,
				TotalData:   int32(totalReceipts),
				TotalPages:  int32(totalPages),
			},
			QueryWord: req.SearchQuery,
		}

		ctx.JSON(http.StatusOK, newRsp)

	case "invoices":
		var pgQuery pgtype.Text
		if err := pgQuery.Scan(req.SearchQuery); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		results, err := server.store.SearchUserInvoices(ctx, pgQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		processedUsernames := make(map[string]struct{})
		var totalInvoices int32
		var returnInvoices []invoiceResponse
		for _, result := range results {
			if _, ok := processedUsernames[result]; ok {
				continue
			}

			invoices, err := server.store.GetUserInvoicesByUsername(ctx, db.GetUserInvoicesByUsernameParams{
				UserInvoiceUsername: result,
				Limit:               PageSize,
				Offset:              (req.PageID - 1) * PageSize,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalInvoice, err := server.store.CountUserInvoicesByUsername(ctx, result)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

			totalInvoices = totalInvoices + int32(totalInvoice)

			for _, invoice := range invoices {
				updatedInvoice, _ := newInvoiceResponse(invoice)
				returnInvoices = append(returnInvoices, updatedInvoice)
			}

			processedUsernames[result] = struct{}{}
		}

		totalPages := totalInvoices / PageSize
		if totalInvoices%PageSize != 0 {
			totalPages++
		}

		rsp := listInvoiceSearchResponse{
			Data: returnInvoices,
			Metadata: PaginationMetadata{
				CurrentPage: int32(req.PageID),
				TotalPages:  int32(totalPages),
				TotalData:   int32(totalInvoices),
			},
			QueryWord: req.SearchQuery,
		}

		ctx.JSON(http.StatusOK, rsp)
	}

	return
}
