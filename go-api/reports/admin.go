package reports

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
)

func (store *ReportStore) GetAdminInvoiceSummary(ctx context.Context, payload ReportSummaryData) ([]GetUserInvoiceSummaryResponse, error) {
	invoiceSummary, err := store.dbStore.GetAllInvoiceSummaryBtwnPeriod(ctx, db.GetAllInvoiceSummaryBtwnPeriodParams{
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
			Username:       invoice.Username,
		})
	}

	return rsp, nil
}

func (store *ReportStore) GetAdminReceiptSummary(ctx context.Context, payload ReportSummaryData) ([]GetUserReceiptSummaryResponse, error) {
	receiptSummary, err := store.dbStore.GetAllReceiptSummaryBtwnPeriod(ctx, db.GetAllReceiptSummaryBtwnPeriodParams{
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
			Username:           receipt.Username,
		})
	}

	return rsp, nil
}
