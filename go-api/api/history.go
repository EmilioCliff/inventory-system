package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (server *Server) testGroup(ctx *gin.Context) {
	invoices, err := server.store.StoreGetUserInvoicesByDate(ctx, 2)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var jsonResponse []map[string]interface{}
	for _, invoice := range invoices {
		var invoiceData []interface{}
		if err := json.Unmarshal(invoice.InvoiceData, &invoiceData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		response := map[string]interface{}{
			"issued_date":  invoice.IssuedDate,
			"num_invoices": invoice.NumInvoices,
			"invoice_data": invoiceData,
		}
		jsonResponse = append(jsonResponse, response)
	}

	err = server.setCache(ctx, TestGroup, jsonResponse)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, jsonResponse)
}

type userHistoryResponse struct {
	IssuedDate  string `json:"issued_date"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
}

type allUserHistoryResponse struct {
	Data userHistoryResponse `json:"data"`
	User string              `json:"user"`
}

type userDeptResponse struct {
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
}

type userDeptResponseName struct {
	Data []userDeptResponse
	User string `json:"user"`
}

type userHistoryRequest struct {
	UserID int32 `uri:"id" binding:"required"`
}

func (server *Server) getAllUsersReceivedHistory(ctx *gin.Context) {
	invoices, err := server.store.StoreGetInvoicesByDate(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	aggregatedData := make(map[string]map[string]map[string]int)

	for _, invoice := range invoices {
		var invoiceData [][]interface{}
		if err := json.Unmarshal(invoice.InvoiceData, &invoiceData); err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		for _, data := range invoiceData {
			for idx, item := range data {
				product := make(map[string]interface{})
				if idx == 0 {
					continue
				}

				itemMap := item.(map[string]interface{})
				for key, value := range itemMap {
					product[key] = value
				}

				productName := product["productName"].(string)
				quantity := int(product["productQuantity"].(float64))
				price := int(product["totalBill"].(float64))

				// issuedDate := invoice.IssuedDate.Format("2006-01-02")
				issuedDate := invoice.IssuedDate.Truncate(10 * time.Minute).Format("2006-01-02 15:04:05")

				// issuedDate := string(invoice.IssuedDate)
				if _, ok := aggregatedData[issuedDate]; !ok {
					aggregatedData[issuedDate] = make(map[string]map[string]int)
				}

				if _, ok := aggregatedData[issuedDate][productName]; !ok {
					aggregatedData[issuedDate][productName] = make(map[string]int)
				}

				aggregatedData[issuedDate][productName]["quantity"] += quantity
				aggregatedData[issuedDate][productName]["totalPrice"] += price
			}
		}
	}

	// var jsonResponse []userHistoryResponse
	// for issuedDate, productData := range aggregatedData {
	// 	for productName, data := range productData {
	// 		response := userHistoryResponse{
	// 			IssuedDate:  issuedDate,
	// 			ProductName: productName,
	// 			Quantity:    data["quantity"],
	// 			Price:       data["totalPrice"],
	// 		}
	// 		jsonResponse = append(jsonResponse, response)
	// 	}
	// }
	err = server.setCache(ctx, AllUserReceiverHistory, aggregatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, aggregatedData)
}

func (server *Server) getUserReceivedHistory(ctx *gin.Context) {
	var req userHistoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	invoices, err := server.store.StoreGetUserInvoicesByDate(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	aggregatedData, err := structureInvoiceProducts(invoices)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// log.Println("aggregated data: ", aggregatedData)
	// var jsonResponse []userHistoryResponse
	// for issuedDate, productData := range aggregatedData {
	// 	for productName, data := range productData {
	// 		response := userHistoryResponse{
	// 			IssuedDate:  issuedDate,
	// 			ProductName: productName,
	// 			Quantity:    data["quantity"],
	// 			Price:       data["totalPrice"],
	// 		}
	// 		jsonResponse = append(jsonResponse, response)
	// 	}
	// }
	err = server.setCache(ctx, UserReceivedHistory+fmt.Sprintf("%v", req.UserID), aggregatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, aggregatedData)
}

func (server *Server) getUserSoldHistory(ctx *gin.Context) {
	var req userHistoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	receipts, err := server.store.StoreGetUserReceiptsByDate(ctx, req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	aggregatedData, err := structureReceiptProducts(receipts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// var jsonResponse []userHistoryResponse
	// for issuedDate, productData := range aggregatedData {
	// 	for productName, data := range productData {
	// 		response := userHistoryResponse{
	// 			IssuedDate:  issuedDate,
	// 			ProductName: productName,
	// 			Quantity:    data["quantity"],
	// 			Price:       data["totalPrice"],
	// 		}
	// 		jsonResponse = append(jsonResponse, response)
	// 	}
	// }
	err = server.setCache(ctx, UserSoldHistory+fmt.Sprintf("%v", req.UserID), aggregatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, aggregatedData)
}

func (server *Server) getUserDebt(ctx *gin.Context) {
	var req userHistoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, int64(req.UserID))
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp, err := structureDebt(user, server, ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.setCache(ctx, UserDebt+fmt.Sprintf("%v", req.UserID), rsp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getAllUserDebt(ctx *gin.Context) {
	users, err := server.store.ListUserNoPagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rsp []userDeptResponseName
	for idx, user := range users {
		if idx == 0 {
			continue
		}
		response, err := structureDebt(user, server, ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		rsp = append(rsp, response)
	}
	log.Println(rsp)

	err = server.setCache(ctx, AllUserDebt, rsp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) adminHistory(ctx *gin.Context) {
	entries, err := server.store.StoreGetEntryByDate(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.setCache(ctx, AdminHistory, entries)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

func structureInvoiceProducts(invoices []db.StoreGetUserInvoicesByDateRow) (map[string]map[string]map[string]int, error) {
	aggregatedData := make(map[string]map[string]map[string]int)

	for _, invoice := range invoices {
		var invoiceData [][]interface{}
		if err := json.Unmarshal(invoice.InvoiceData, &invoiceData); err != nil {
			return nil, err
		}

		for _, data := range invoiceData {
			for idx, item := range data {
				product := make(map[string]interface{})
				if idx == 0 {
					continue
				}

				itemMap := item.(map[string]interface{})
				for key, value := range itemMap {
					product[key] = value
				}

				productName := product["productName"].(string)
				quantity := int(product["productQuantity"].(float64))
				price := int(product["totalBill"].(float64))

				// issuedDate := invoice.IssuedDate.Format("2006-01-02")
				issuedDate := invoice.IssuedDate.Truncate(10 * time.Minute).Format("2006-01-02 15:04:05")

				// issuedDate := string(invoice.IssuedDate)
				if _, ok := aggregatedData[issuedDate]; !ok {
					aggregatedData[issuedDate] = make(map[string]map[string]int)
				}

				if _, ok := aggregatedData[issuedDate][productName]; !ok {
					aggregatedData[issuedDate][productName] = make(map[string]int)
				}

				aggregatedData[issuedDate][productName]["quantity"] += quantity
				aggregatedData[issuedDate][productName]["totalPrice"] += price
			}
		}
	}

	return aggregatedData, nil
}

func structureReceiptProducts(receipts []db.StoreGetUserReceiptsByDateRow) (map[string]map[string]map[string]int, error) {
	aggregatedData := make(map[string]map[string]map[string]int)

	for _, receipt := range receipts {
		var receiptData [][]interface{}
		if err := json.Unmarshal(receipt.ReceiptData, &receiptData); err != nil {
			return nil, err
		}

		for _, data := range receiptData {
			for idx, item := range data {
				product := make(map[string]interface{})
				if idx == 0 {
					continue
				}

				itemMap := item.(map[string]interface{})
				for key, value := range itemMap {
					product[key] = value
				}

				productName := product["productName"].(string)
				quantity := int(product["productQuantity"].(float64))
				price := int(product["totalBill"].(float64))

				// issuedDate := receipt.IssuedDate.Format("2006-01-02")

				issuedDate := receipt.IssuedDate.Truncate(10 * time.Minute).Format("2006-01-02 15:04:05")
				// issuedDate := string(receipt.IssuedDate)
				if _, ok := aggregatedData[issuedDate]; !ok {
					aggregatedData[issuedDate] = make(map[string]map[string]int)
				}

				if _, ok := aggregatedData[issuedDate][productName]; !ok {
					aggregatedData[issuedDate][productName] = make(map[string]int)
				}

				aggregatedData[issuedDate][productName]["quantity"] += quantity
				aggregatedData[issuedDate][productName]["totalPrice"] += price
			}
		}
	}

	return aggregatedData, nil
}

func structureDebt(user db.User, server *Server, ctx *gin.Context) (userDeptResponseName, error) {
	var userProductData []map[string]interface{}
	if user.Stock != nil {
		if err := json.Unmarshal(user.Stock, &userProductData); err != nil {
			return userDeptResponseName{}, err
		}
	}

	var rsp []userDeptResponse
	for _, productData := range userProductData {
		if productData["productQuantity"].(float64) < 1 {
			continue
		}
		product, err := server.store.GetProduct(ctx, int64(productData["productID"].(float64)))
		if err != nil {
			return userDeptResponseName{}, err
		}

		response := userDeptResponse{
			ProductName: productData["productName"].(string),
			Price:       productData["productQuantity"].(float64) * float64(product.UnitPrice),
			Quantity:    productData["productQuantity"].(float64),
		}

		rsp = append(rsp, response)
	}

	rspName := userDeptResponseName{
		Data: rsp,
		User: user.Username,
	}

	return rspName, nil
}

func structureAllInvoiceProducts(invoices []db.StoreGetInvoicesByDateRow) (map[string]map[string]map[string]int, error) {
	aggregatedData := make(map[string]map[string]map[string]int)

	for _, invoice := range invoices {
		var invoiceData [][]interface{}
		if err := json.Unmarshal(invoice.InvoiceData, &invoiceData); err != nil {
			return nil, err
		}

		for _, data := range invoiceData {
			for idx, item := range data {
				product := make(map[string]interface{})
				if idx == 0 {
					continue
				}

				itemMap := item.(map[string]interface{})
				for key, value := range itemMap {
					product[key] = value
				}

				productName := product["productName"].(string)
				quantity := int(product["productQuantity"].(float64))
				price := int(product["totalBill"].(float64))

				issuedDate := invoice.IssuedDate.Format("2006-01-02")

				// issuedDate := string(invoice.IssuedDate)
				if _, ok := aggregatedData[issuedDate]; !ok {
					aggregatedData[issuedDate] = make(map[string]map[string]int)
				}

				if _, ok := aggregatedData[issuedDate][productName]; !ok {
					aggregatedData[issuedDate][productName] = make(map[string]int)
				}

				aggregatedData[issuedDate][productName]["quantity"] += quantity
				aggregatedData[issuedDate][productName]["totalPrice"] += price
			}
		}
	}

	return aggregatedData, nil
}
