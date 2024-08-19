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
	"github.com/rs/zerolog/log"
)

func newPurchaseOrder(orders []db.PurchaseOrder) ([]purchaseOrdersResponse, error) {
	var result []purchaseOrdersResponse
	for _, order := range orders {
		var orderData []map[string]interface{}
		if order.LpoData != nil {
			if unerr := json.Unmarshal(order.LpoData, &orderData); unerr != nil {
				return []purchaseOrdersResponse{}, unerr
			}
		}

		data := purchaseOrdersResponse{
			ID:              order.ID,
			SupplierName:    order.SupplierName,
			SupplierAddress: order.SupplierAddress,
			OrderData:       orderData,
			CreatedAt:       order.CreatedAt,
		}

		result = append(result, data)
	}

	return result, nil
}

type purchaseOrdersResponse struct {
	ID              string                   `json:"id"`
	SupplierName    string                   `json:"supplier_name"`
	SupplierAddress string                   `json:"supplier_addres"`
	OrderData       []map[string]interface{} `json:"order_data"`
	CreatedAt       time.Time                `json:"created_at"`
}

type createPurchaseOrderRequest struct {
	SupplierName string                   `json:"supplier_name" binding:"required"`
	PoBox        string                   `json:"po_box" binding:"required"`
	Address      string                   `json:"address" binding:"required"`
	Data         []map[string]interface{} `json:"data" binding:"required"`
}

func (server *Server) createPurchaseOrder(ctx *gin.Context) {
	var req createPurchaseOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// data := map[string]interface{}{
	// 	"supplier_name": "Eldohosp Pharmaceuticals LTD",
	// 	"po_box":        "66859-00800",
	// 	"address":       "epl plaza nairobi",
	// 	"order_date":    time.Now().Format("2006-01-02"),
	// 	"order_number":  time.Now().Format("20060102150405"),
	// 	"Data": []map[string]interface{}{
	// 		{
	// 			"product_name": "HIV 1-2 Test Cassette (Self test)",
	// 			"quantity":     100,
	// 			"unit_price":   1000,
	// 		},
	// 	},
	// }

	// save to db then create

	purchaseOrderNumber := time.Now().Format("20060102150405")

	data := map[string]interface{}{
		"supplier_name": req.SupplierName,
		"po_box":        req.PoBox,
		"address":       req.Address,
		"order_date":    time.Now().Format("2006-01-02"),
		"order_number":  purchaseOrderNumber,
		"Data":          req.Data,
	}

	result := make(chan struct {
		pdfBytes []byte
		err      error
	}, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func(data map[string]interface{}) {
		defer wg.Done()
		pdfBytes, err := utils.GeneratePurchaseOrder(data)
		result <- struct {
			pdfBytes []byte
			err      error
		}{pdfBytes: pdfBytes, err: err}
	}(data)

	wg.Wait()
	close(result)

	lpoResult := <-result
	if lpoResult.err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(lpoResult.err))
		return
	}

	jsonLpoData, err := json.Marshal(req.Data)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(lpoResult.err))
		return
	}

	orderCreated, err := server.store.CreatePurchaseOrder(ctx, db.CreatePurchaseOrderParams{
		ID:              purchaseOrderNumber,
		SupplierName:    req.SupplierName,
		SupplierAddress: fmt.Sprintf("P.O Box %s - %s", req.PoBox, req.Address),
		LpoData:         jsonLpoData,
		LpoPdf:          lpoResult.pdfBytes,
	})

	rsp := map[string]interface{}{
		"purchase_order_pdf": orderCreated.LpoPdf,
		// "purchase_order_pdf": lpoResult.pdfBytes,
		// "data":          data,
	}

	ctx.JSON(http.StatusOK, rsp)
}

type listPurchaseOrdersRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
}

type listPurchaseOrdersResponse struct {
	Data     []purchaseOrdersResponse `json:"data"`
	Metadata PaginationMetadata       `json:"metadata"`
}

func (server *Server) listPurchaseOrders(ctx *gin.Context) {
	var req listPurchaseOrdersRequest
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

	listOrders, err := server.store.ListPurchaseOrder(ctx, db.ListPurchaseOrderParams{
		Limit:  PageSize,
		Offset: (req.PageID - 1) * PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPurchaseOrders, err := server.store.CountPurchaseOrders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	totalPages := totalPurchaseOrders / int64(PageSize)
	if totalPurchaseOrders/int64(PageSize) != 0 {
		totalPages++
	}

	data, err := newPurchaseOrder(listOrders)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := listPurchaseOrdersResponse{
		Data: data,
		Metadata: PaginationMetadata{
			CurrentPage: req.PageID,
			TotalPages:  int32(totalPages),
			TotalData:   int32(totalPurchaseOrders),
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

type deletePurchaseOrderUri struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) deletePurchaseOrders(ctx *gin.Context) {
	var uri deletePurchaseOrderUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeletePurchaseOrder(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := server.redis.Del(ctx, ListPurchaseOrders+fmt.Sprintf(":1")).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "deleted successful"})
}

type downloadPurchaseOrderUri struct {
	ID string `uri:"id" binding:"required"`
}

func (server *Server) downloadPurchaseOrders(ctx *gin.Context) {
	var uri downloadPurchaseOrderUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	purchaseOrder, err := server.store.GetPurchaseOrder(ctx, uri.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := map[string]interface{}{
		"purchase_order_pdf": purchaseOrder.LpoPdf,
	}

	ctx.JSON(http.StatusOK, rsp)
}
