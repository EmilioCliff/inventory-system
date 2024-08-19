package reports

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/xuri/excelize/v2"
)

func (r *ReportStore) invoiceSummary(payload ReportsPayload, f *excelize.File, sheetName string, ctx context.Context, columnHeaders []string, user_id int64, styles []int) error {
	col1 := fmt.Sprintf("%s", columnHeaders[0])
	col2 := fmt.Sprintf("%s", columnHeaders[1])
	col3 := fmt.Sprintf("%s", columnHeaders[2])
	col4 := fmt.Sprintf("%s", columnHeaders[3])

	f.SetColWidth(sheetName, col1, col1, 20)
	f.SetColWidth(sheetName, col2, col2, 40)
	f.SetColWidth(sheetName, col3, col4, 20)

	f.SetCellValue(sheetName, col1+"1", "Invoice Number")
	f.SetCellValue(sheetName, col2+"1", "Invoice Data")
	f.SetCellValue(sheetName, col3+"1", "Amount")
	f.SetCellValue(sheetName, col4+"1", "Invoice Date")

	err := f.SetColStyle(sheetName, col4, styles[1]) // date style
	if err != nil {
		return err
	}

	err = f.SetColStyle(sheetName, col3, styles[2]) // money style
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheetName, col1+"1", col4+"1", styles[0]) // header style
	if err != nil {
		return err
	}

	summaries, err := r.GetUserInvoiceSummary(ctx, ReportSummaryData{
		UserID:   user_id,
		FromDate: payload.FromDate,
		ToDate:   payload.ToDate,
	})
	if err != nil {
		return err
	}

	rowNumber := 1
	for _, summary := range summaries {
		rowNumber += 1

		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col1, rowNumber), summary.InvoiceNumber)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col2, rowNumber), summary.AdditionalData)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col3, rowNumber), summary.TotalInvoice)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col4, rowNumber), summary.InvoiceDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReportStore) receiptSummary(payload ReportsPayload, f *excelize.File, sheetName string, ctx context.Context, columnHeaders []string, user_id int64, styles []int) error {
	col1 := fmt.Sprintf("%s", columnHeaders[0])
	col2 := fmt.Sprintf("%s", columnHeaders[1])
	col3 := fmt.Sprintf("%s", columnHeaders[2])
	col4 := fmt.Sprintf("%s", columnHeaders[3])
	col5 := fmt.Sprintf("%s", columnHeaders[4])
	col6 := fmt.Sprintf("%s", columnHeaders[5])
	col7 := fmt.Sprintf("%s", columnHeaders[6])

	f.SetColWidth(sheetName, col1, col3, 20)
	f.SetColWidth(sheetName, col4, col4, 40)
	f.SetColWidth(sheetName, col5, col7, 20)

	f.SetCellValue(sheetName, col1+"1", "Mpesa Receipt Number")
	f.SetCellValue(sheetName, col2+"1", "Receipt Number")
	f.SetCellValue(sheetName, col3+"1", "Phone Number")
	f.SetCellValue(sheetName, col4+"1", "Receipt Data")
	f.SetCellValue(sheetName, col5+"1", "Amount")
	f.SetCellValue(sheetName, col6+"1", "Payment Method")
	f.SetCellValue(sheetName, col7+"1", "Receipt Date")

	err := f.SetColStyle(sheetName, col7, styles[1]) // date style
	if err != nil {
		return err
	}

	err = f.SetColStyle(sheetName, col5, styles[2]) // money style
	if err != nil {
		return err
	}

	paymentStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})
	if err != nil {
		return err
	}

	err = f.SetColStyle(sheetName, col6, paymentStyle) // payment alignment
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheetName, col1+"1", col7+"1", styles[0]) // header style
	if err != nil {
		return err
	}

	summaries, err := r.GetUserReceiptSummary(ctx, ReportSummaryData{
		UserID:   user_id,
		FromDate: payload.FromDate,
		ToDate:   payload.ToDate,
	})
	if err != nil {
		return err
	}

	rowNumber := 1
	for _, summary := range summaries {
		rowNumber += 1

		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col1, rowNumber), summary.MpesaReceiptNumber)
		if err != nil {
			return err
		}

		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col2, rowNumber), summary.ReceiptNumber)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col3, rowNumber), summary.PhoneNumber)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col4, rowNumber), summary.AdditionalData)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col5, rowNumber), summary.TotalReceipt)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col6, rowNumber), summary.PaymentMethod)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col7, rowNumber), summary.ReceiptDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReportStore) userAvailableStock(f *excelize.File, sheetName string, columnHeaders []string, userData []byte, styles []int) error {
	col1 := fmt.Sprintf("%s", columnHeaders[0])
	col2 := fmt.Sprintf("%s", columnHeaders[1])

	f.SetColWidth(sheetName, col1, col2, 20)

	f.SetCellValue(sheetName, col1+"1", "Product Name")
	f.SetCellValue(sheetName, col2+"1", "Product Quantity")

	err := f.SetColStyle(sheetName, col2, styles[1]) // quantity style
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheetName, col1+"1", col2+"1", styles[0]) // header style
	if err != nil {
		return err
	}

	if userData == nil {
		return nil
	}

	var userStockData []map[string]interface{}
	if err := json.Unmarshal(userData, &userStockData); err != nil {
		return err
	}

	rowNumber := 1
	for _, data := range userStockData {
		rowNumber += 1

		productName := data["productName"].(string)
		quantity := data["productQuantity"].(float64)

		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col1, rowNumber), productName)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col2, rowNumber), quantity)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *ReportStore) adminInvoiceSummary(payload ReportsPayload, f *excelize.File, sheetName string, ctx context.Context, columnHeaders []string, styles []int) error {
	col1 := fmt.Sprintf("%s", columnHeaders[0])
	col2 := fmt.Sprintf("%s", columnHeaders[1])
	col3 := fmt.Sprintf("%s", columnHeaders[2])
	col4 := fmt.Sprintf("%s", columnHeaders[3])
	col5 := fmt.Sprintf("%s", columnHeaders[4])

	f.SetColWidth(sheetName, col1, col2, 20)
	f.SetColWidth(sheetName, col3, col3, 40)
	f.SetColWidth(sheetName, col4, col5, 20)

	f.SetCellValue(sheetName, col1+"1", "Username")
	f.SetCellValue(sheetName, col2+"1", "Invoice Number")
	f.SetCellValue(sheetName, col3+"1", "Invoice Data")
	f.SetCellValue(sheetName, col4+"1", "Amount")
	f.SetCellValue(sheetName, col5+"1", "Invoice Date")

	err := f.SetColStyle(sheetName, col5, styles[1]) // date style
	if err != nil {
		return err
	}

	err = f.SetColStyle(sheetName, col4, styles[2]) // money style
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheetName, col1+"1", col5+"1", styles[0]) // header style
	if err != nil {
		return err
	}

	summaries, err := r.GetAdminInvoiceSummary(ctx, ReportSummaryData{
		FromDate: payload.FromDate,
		ToDate:   payload.ToDate,
	})
	if err != nil {
		return err
	}

	rowNumber := 1
	for _, summary := range summaries {
		rowNumber += 1

		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col1, rowNumber), summary.Username)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col2, rowNumber), summary.InvoiceNumber)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col3, rowNumber), summary.AdditionalData)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col4, rowNumber), summary.TotalInvoice)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col5, rowNumber), summary.InvoiceDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReportStore) adminReceiptSummary(payload ReportsPayload, f *excelize.File, sheetName string, ctx context.Context, columnHeaders []string, styles []int) error {
	col1 := fmt.Sprintf("%s", columnHeaders[0])
	col2 := fmt.Sprintf("%s", columnHeaders[1])
	col3 := fmt.Sprintf("%s", columnHeaders[2])
	col4 := fmt.Sprintf("%s", columnHeaders[3])
	col5 := fmt.Sprintf("%s", columnHeaders[4])
	col6 := fmt.Sprintf("%s", columnHeaders[5])
	col7 := fmt.Sprintf("%s", columnHeaders[6])
	col8 := fmt.Sprintf("%s", columnHeaders[7])

	f.SetColWidth(sheetName, col1, col4, 20)
	f.SetColWidth(sheetName, col5, col5, 40)
	f.SetColWidth(sheetName, col6, col8, 20)

	f.SetCellValue(sheetName, col1+"1", "Username")
	f.SetCellValue(sheetName, col2+"1", "Mpesa Receipt Number")
	f.SetCellValue(sheetName, col3+"1", "Receipt Number")
	f.SetCellValue(sheetName, col4+"1", "Phone Number")
	f.SetCellValue(sheetName, col5+"1", "Receipt Data")
	f.SetCellValue(sheetName, col6+"1", "Amount")
	f.SetCellValue(sheetName, col7+"1", "Payment Method")
	f.SetCellValue(sheetName, col8+"1", "Receipt Date")

	err := f.SetColStyle(sheetName, col8, styles[1]) // date style
	if err != nil {
		return err
	}

	err = f.SetColStyle(sheetName, col6, styles[2]) // money style
	if err != nil {
		return err
	}

	paymentStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})
	if err != nil {
		return err
	}

	err = f.SetColStyle(sheetName, col7, paymentStyle) // payment alignment
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheetName, col1+"1", col8+"1", styles[0]) // header style
	if err != nil {
		return err
	}

	summaries, err := r.GetAdminReceiptSummary(ctx, ReportSummaryData{
		FromDate: payload.FromDate,
		ToDate:   payload.ToDate,
	})
	if err != nil {
		return err
	}

	rowNumber := 1
	for _, summary := range summaries {
		rowNumber += 1

		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col1, rowNumber), summary.Username)
		if err != nil {
			return err
		}

		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col2, rowNumber), summary.MpesaReceiptNumber)
		if err != nil {
			return err
		}

		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col3, rowNumber), summary.ReceiptNumber)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col4, rowNumber), summary.PhoneNumber)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col5, rowNumber), summary.AdditionalData)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col6, rowNumber), summary.TotalReceipt)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col7, rowNumber), summary.PaymentMethod)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col8, rowNumber), summary.ReceiptDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReportStore) localPurchaseOrders(payload ReportsPayload, f *excelize.File, sheetName string, ctx context.Context, columnHeaders []string, styles []int) error {
	col1 := fmt.Sprintf("%s", columnHeaders[0])
	col2 := fmt.Sprintf("%s", columnHeaders[1])
	col3 := fmt.Sprintf("%s", columnHeaders[2])
	col4 := fmt.Sprintf("%s", columnHeaders[3])
	col5 := fmt.Sprintf("%s", columnHeaders[4])

	f.SetColWidth(sheetName, col1, col2, 20)
	f.SetColWidth(sheetName, col3, col3, 40)
	f.SetColWidth(sheetName, col4, col5, 20)

	f.SetCellValue(sheetName, col1+"1", "LPO ID")
	f.SetCellValue(sheetName, col2+"1", "Supplier Name")
	f.SetCellValue(sheetName, col3+"1", "LPO Data")
	f.SetCellValue(sheetName, col4+"1", "Amount")
	f.SetCellValue(sheetName, col5+"1", "LPO Date")

	err := f.SetColStyle(sheetName, col5, styles[1]) // date style
	if err != nil {
		return err
	}

	err = f.SetColStyle(sheetName, col4, styles[2]) // money style
	if err != nil {
		return err
	}

	err = f.SetCellStyle(sheetName, col1+"1", col5+"1", styles[0]) // header style
	if err != nil {
		return err
	}

	orders, err := r.dbStore.GetAdminPurchaseOrders(ctx, db.GetAdminPurchaseOrdersParams{
		FromDate: payload.FromDate,
		ToDate:   payload.ToDate,
	})
	if err != nil {
		return err
	}

	rowNumber := 1
	for _, order := range orders {
		rowNumber += 1

		dataString, err := structureLPOData(order.LpoData)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col1, rowNumber), order.ID)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col2, rowNumber), order.SupplierName)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col3, rowNumber), dataString)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col4, rowNumber), order.Total)
		if err != nil {
			return err
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("%s%v", col5, rowNumber), order.CreatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func structureLPOData(data []byte) (string, error) {
	var jsonData []map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return "", err
	}

	dataString := ""
	for _, d := range jsonData {
		productName := d["product_name"].(string)
		productPrice := d["unit_price"].(float64)
		quantity := d["quantity"].(float64)
		dataString += fmt.Sprintf("%s: %vx%v; ", productName, quantity, productPrice)
	}

	return dataString, nil
}
