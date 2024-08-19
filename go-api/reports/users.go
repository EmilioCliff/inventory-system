package reports

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
)

var (
	NoInvoiceData = errors.New("no invoice data")
	NoReceiptData = errors.New("no receipt data")
)

type ReportSummaryData struct {
	UserID   int64     `json:"user_id,omitempty"`
	ToDate   time.Time `json:"to_date"`
	FromDate time.Time `json:"from_date"`
}

type ProductSummary struct {
	ProductName  string `json:"product_name"`
	ProductPrice int32  `json:"product_price"`
	Quantity     int64  `json:"quantity"`
}

type GetUserInvoiceSummaryResponse struct {
	InvoiceNumber  string    `json:"invoice_number"`
	AdditionalData string    `json:"additional_data"`
	TotalInvoice   int64     `json:"total_invoice"`
	InvoiceDate    time.Time `json:"invoice_date"`
	Username       string    `json:"username,omitempty"`
}

func (store *ReportStore) GetUserInvoiceSummary(ctx context.Context, payload ReportSummaryData) ([]GetUserInvoiceSummaryResponse, error) {
	invoiceSummary, err := store.dbStore.GetUserInvoiceSummaryBtwnPeriod(ctx, db.GetUserInvoiceSummaryBtwnPeriodParams{
		UserID:   int32(payload.UserID),
		FromDate: payload.FromDate,
		ToDate:   formartDate(payload.ToDate),
	})
	if err != nil {
		return []GetUserInvoiceSummaryResponse{}, fmt.Errorf("failed to get invoice summary: %w", err)
	}

	var rsp []GetUserInvoiceSummaryResponse
	for _, invoice := range invoiceSummary {
		var invoiceData []map[string]interface{}
		if err := json.Unmarshal(invoice.InvoiceData, &invoiceData); err != nil {
			return []GetUserInvoiceSummaryResponse{}, fmt.Errorf("failed to unmarshal invoice data: %w", err)
		}

		additionaData := ""
		for idx, data := range invoiceData {
			if idx == 0 {
				continue
			}
			productName := data["productName"].(string)
			productQuantity := data["productQuantity"].(float64)
			productPrice := data["totalBill"].(float64) / productQuantity
			additionaData += fmt.Sprintf("%s: %vx%v; ", productName, productQuantity, productPrice)
		}

		rsp = append(rsp, GetUserInvoiceSummaryResponse{
			InvoiceNumber:  invoice.InvoiceNumber,
			TotalInvoice:   invoice.Total.(int64),
			InvoiceDate:    invoice.InvoiceDate,
			AdditionalData: additionaData,
		})
	}

	return rsp, nil
}

type GetUserReceiptSummaryResponse struct {
	ReceiptNumber      string `json:"receipt_number"`
	MpesaReceiptNumber string `json:"mpesa_receipt_number"`
	PhoneNumber        string `json:"phone_number"`
	// ReceiptData    []ProductSummary `json:"receipt_data"`
	AdditionalData string    `json:"additional_data"`
	TotalReceipt   int64     `json:"total_receipt"`
	PaymentMethod  string    `json:"payment_method"`
	ReceiptDate    time.Time `json:"receipt_date"`
	Username       string    `json:"username,omitempty"`
}

func (store *ReportStore) GetUserReceiptSummary(ctx context.Context, payload ReportSummaryData) ([]GetUserReceiptSummaryResponse, error) {
	receiptSummary, err := store.dbStore.GetUserReceiptSummaryBtwnPeriod(ctx, db.GetUserReceiptSummaryBtwnPeriodParams{
		UserID:   int32(payload.UserID),
		FromDate: payload.FromDate,
		ToDate:   formartDate(payload.ToDate),
	})

	if err != nil {
		return []GetUserReceiptSummaryResponse{}, fmt.Errorf("failed to get receipt summary: %w", err)
	}

	var rsp []GetUserReceiptSummaryResponse
	for _, receipt := range receiptSummary {
		var receiptData []map[string]interface{}
		if err := json.Unmarshal(receipt.ReceiptData, &receiptData); err != nil {
			return []GetUserReceiptSummaryResponse{}, fmt.Errorf("failed to unmarshal receipt data: %w", err)
		}

		var additionaData string
		for _, data := range receiptData {
			productName := data["productName"].(string)
			productPrice := data["productPrice"].(float64)
			quantity := data["productQuantity"].(float64)
			additionaData += fmt.Sprintf("%s: %vx%v; ", productName, quantity, productPrice)
		}

		rsp = append(rsp, GetUserReceiptSummaryResponse{
			TotalReceipt:  int64(receipt.Amount),
			ReceiptNumber: receipt.ReceiptNumber,
			// ReceiptData:   receiptDataSold,
			AdditionalData:     additionaData,
			PaymentMethod:      receipt.PaymentMethod,
			ReceiptDate:        receipt.CreatedAt,
			MpesaReceiptNumber: receipt.MpesaReceiptNumber,
			PhoneNumber:        receipt.PhoneNumber,
		})
	}

	return rsp, nil
}

type GetAdminPurchaseHistoryResponse struct {
	ProductName  string    `json:"product_name"`
	ProductPrice int32     `json:"product_price"`
	Quantity     int32     `json:"quantity"`
	PurchaseDate time.Time `json:"purchase_date"`
}

func (store *ReportStore) GetAdminPurchaseHistory(ctx context.Context, payload ReportSummaryData) ([]GetAdminPurchaseHistoryResponse, error) {
	entries, err := store.dbStore.ListEntries(ctx, db.ListEntriesParams{
		FromDate: payload.FromDate,
		ToDate:   formartDate(payload.ToDate),
	})
	if err != nil {
		return []GetAdminPurchaseHistoryResponse{}, fmt.Errorf("failed to list entries: %w", err)
	}

	var rsp []GetAdminPurchaseHistoryResponse

	for _, entry := range entries {
		data := GetAdminPurchaseHistoryResponse{
			ProductName:  entry.ProductName,
			ProductPrice: entry.ProductPrice,
			Quantity:     entry.QuantityAdded,
			PurchaseDate: entry.CreatedAt,
		}

		rsp = append(rsp, data)
	}
	return rsp, nil
}

type GetAdminSalesHistoryResponse struct {
	Username          string `json:"username"`
	StockDistributed  int64  `json:"stock_distributed" default:"0"`
	StockPaid         int64  `json:"stock_paid" default:"0"`
	CurrentStockValue int64  `json:"current_stock_value" default:"0"`
}

func (store *ReportStore) GetUserHistorySummary(ctx context.Context, payload ReportSummaryData) (GetAdminSalesHistoryResponse, error) {
	totalInvoice, err := store.dbStore.GetUserStockDistributed(ctx, db.GetUserStockDistributedParams{
		UserID:   int32(payload.UserID),
		FromDate: payload.FromDate,
		ToDate:   payload.ToDate,
	})
	if err != nil {
		log.Printf("error 1: %v, value: %v", err, totalInvoice)
		return GetAdminSalesHistoryResponse{}, NoInvoiceData
	}

	totalReceipt, err := store.dbStore.GetUserReceiptPaidTotal(ctx, db.GetUserReceiptPaidTotalParams{
		UserID:   int32(payload.UserID),
		FromDate: payload.FromDate,
		ToDate:   formartDate(payload.ToDate),
	})
	if err != nil {
		log.Printf("error 2: %v", err)
		return GetAdminSalesHistoryResponse{}, fmt.Errorf("failed to get user receipt paid: %w", err)
	}

	userCurrentStockValue, err := store.dbStore.GetUserStockValue(ctx, int32(payload.UserID))
	if err != nil {
		log.Printf("error 3: %v", err)
		return GetAdminSalesHistoryResponse{}, fmt.Errorf("failed to get user current stock value: %w", err)
	}

	username, err := store.dbStore.GetUserUsername(ctx, payload.UserID)
	if err != nil {
		return GetAdminSalesHistoryResponse{}, fmt.Errorf("failed to get user username: %w", err)
	}
	log.Printf("totalInvoice: %v; totalReceipt: %v", totalInvoice, totalReceipt)

	rsp := GetAdminSalesHistoryResponse{
		Username:          username,
		StockDistributed:  totalInvoice.(int64),
		StockPaid:         totalReceipt.(int64),
		CurrentStockValue: userCurrentStockValue.Value,
	}
	return rsp, nil
}

func formartDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, int(time.Millisecond*999999), date.Location())
}
