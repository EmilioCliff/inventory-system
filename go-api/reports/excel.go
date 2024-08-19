package reports

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/xuri/excelize/v2"
)

type ReportsPayload struct {
	FromDate time.Time `json:"from_date"`
	ToDate   time.Time `json:"to_date"`
}

func (r *ReportStore) GenerateUserExcel(ctx context.Context, payload ReportsPayload) ([]byte, error) {
	f := excelize.NewFile()

	headerStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "center",
			Indent:          1,
			JustifyLastLine: true,
			ReadingOrder:    0,
			RelativeIndent:  1,
			ShrinkToFit:     true,
			Vertical:        "top",
			WrapText:        true,
		},
		Font: &excelize.Font{
			Bold:   true,
			Italic: true,
			Family: "Times New Roman",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#0094B7"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	dateStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 14,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})
	if err != nil {
		return nil, err
	}

	moneyStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 4,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})

	quantityStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 3,
	})
	if err != nil {
		return nil, err
	}

	sheet := "Sheet1"

	f.SetColWidth(sheet, "A", "D", 20)

	f.SetCellValue(sheet, "A1", "Username")
	f.SetCellValue(sheet, "B1", "Stock Distributed")
	f.SetCellValue(sheet, "C1", "Stock Paid")
	f.SetCellValue(sheet, "D1", "Current Stock Value")

	err = f.SetColStyle(sheet, "B:D", moneyStyle)
	if err != nil {
		return nil, err
	}

	err = f.SetCellStyle("Sheet1", "A1", "D1", headerStyle)
	if err != nil {
		return nil, err
	}

	users, err := r.dbStore.ListUserNoPagination(ctx)
	if err != nil {
		return nil, err
	}

	ch := make(chan error, 1)
	wg := sync.WaitGroup{}

	rowNumber := 1
	for _, user := range users {
		if user.Role == "admin" {
			continue
		}
		data, err := r.GetUserHistorySummary(ctx, ReportSummaryData{
			UserID:   user.UserID,
			FromDate: payload.FromDate,
			ToDate:   payload.ToDate,
		})
		log.Printf("%v", data)
		if err != nil {
			if errors.Is(err, NoInvoiceData) || errors.Is(err, NoReceiptData) {
				continue
			}

			return nil, err
		}

		rowNumber += 1

		err = f.SetCellValue(sheet, fmt.Sprintf("A%v", rowNumber), data.Username)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("B%v", rowNumber), data.StockDistributed)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("C%v", rowNumber), data.StockPaid)
		if err != nil {
			return nil, err
		}
		err = f.SetCellValue(sheet, fmt.Sprintf("D%v", rowNumber), data.CurrentStockValue)
		if err != nil {
			return nil, err
		}

		sheetName := fmt.Sprintf("%s Sheet", user.Username)
		_, err = f.NewSheet(sheetName)
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(user db.User) {
			defer wg.Done()
			err = r.invoiceSummary(payload, f, sheetName, ctx, []string{"A", "B", "C", "D"}, user.UserID, []int{headerStyle, dateStyle, moneyStyle})
			if err != nil {
				ch <- err
				return
			}

			err = r.receiptSummary(payload, f, sheetName, ctx, []string{"F", "G", "H", "I", "J", "K", "L"}, user.UserID, []int{headerStyle, dateStyle, moneyStyle})
			if err != nil {
				ch <- err
				return
			}

			err = r.userAvailableStock(f, sheetName, []string{"N", "O"}, user.Stock, []int{headerStyle, quantityStyle})
			if err != nil {
				ch <- err
				return
			}
		}(user)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for err := range ch {
		if err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	// if err := f.SaveAs("book1.xlsx"); err != nil {
	// 	fmt.Println(err)
	// }

	if err := f.Close(); err != nil {
		fmt.Println(err)
	}

	// return nil, nil

	return buf.Bytes(), nil
}

func (r *ReportStore) GenerateManagerReports(ctx context.Context, payload ReportsPayload) ([]byte, error) {
	f := excelize.NewFile()

	headerStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "center",
			Indent:          1,
			JustifyLastLine: true,
			ReadingOrder:    0,
			RelativeIndent:  1,
			ShrinkToFit:     true,
			Vertical:        "top",
			WrapText:        true,
		},
		Font: &excelize.Font{
			Bold:   true,
			Italic: true,
			Family: "Times New Roman",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#0094B7"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, err
	}

	dateStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 14,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})
	if err != nil {
		return nil, err
	}

	moneyStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 4,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})

	quantityStyle, err := f.NewStyle(&excelize.Style{
		NumFmt: 3,
	})
	if err != nil {
		return nil, err
	}

	sheet := "Sheet1"

	f.SetColWidth(sheet, "A", "D", 20)

	f.SetCellValue(sheet, "A1", "Product Name")
	f.SetCellValue(sheet, "B1", "Quantity")
	f.SetCellValue(sheet, "C1", "Price")
	f.SetCellValue(sheet, "D1", "Date")

	err = f.SetColStyle(sheet, "C", moneyStyle)
	if err != nil {
		return nil, err
	}

	err = f.SetColStyle(sheet, "D", dateStyle)
	if err != nil {
		return nil, err
	}

	err = f.SetColStyle(sheet, "B", quantityStyle)
	if err != nil {
		return nil, err
	}

	err = f.SetCellStyle("Sheet1", "A1", "D1", headerStyle)
	if err != nil {
		return nil, err
	}

	ch := make(chan error, 1)
	wg := sync.WaitGroup{}

	// writes this sheet in a go routine
	wg.Add(1)
	go func() {
		defer wg.Done()
		admin, err := r.dbStore.GetUser(ctx, 1)
		if err != nil {
			ch <- err
			return
		}

		entries, err := r.GetAdminPurchaseHistory(ctx, ReportSummaryData{
			FromDate: payload.FromDate,
			ToDate:   payload.ToDate,
		})
		if err != nil {
			ch <- err
			return
		}

		log.Printf("Entries response: %v", entries)

		// purchase history
		rowNumber := 1
		for _, entry := range entries {
			rowNumber += 1
			err = f.SetCellValue(sheet, fmt.Sprintf("A%v", rowNumber), entry.ProductName)
			if err != nil {
				ch <- err
				return
			}
			err = f.SetCellValue(sheet, fmt.Sprintf("B%v", rowNumber), entry.Quantity)
			if err != nil {
				ch <- err
				return
			}
			err = f.SetCellValue(sheet, fmt.Sprintf("C%v", rowNumber), entry.ProductPrice)
			if err != nil {
				ch <- err
				return
			}
			err = f.SetCellValue(sheet, fmt.Sprintf("D%v", rowNumber), entry.PurchaseDate)
			if err != nil {
				ch <- err
				return
			}
		}

		log.Println("Done with purchase history")

		// in stock history
		err = r.userAvailableStock(f, sheet, []string{"F", "G"}, admin.Stock, []int{headerStyle, quantityStyle})
		if err != nil {
			ch <- err
			return
		}

		log.Println("Done with in stock history")
	}()

	sheets := []string{"Invoices", "Receipts", "lpo"}

	for _, sheet := range sheets {
		switch sheet {
		case "Invoices":
			sheetName := "Invoices Sheet"
			_, err = f.NewSheet(sheetName)
			if err != nil {
				return nil, err
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = r.adminInvoiceSummary(payload, f, sheetName, ctx, []string{"A", "B", "C", "D", "E"}, []int{headerStyle, dateStyle, moneyStyle})
				if err != nil {
					ch <- err
					return
				}
				log.Println("Done with invoices")
			}()

		case "Receipts":
			sheetName := "Receipts Sheet"
			_, err = f.NewSheet(sheetName)
			if err != nil {
				return nil, err
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = r.adminReceiptSummary(payload, f, sheetName, ctx, []string{"A", "B", "C", "D", "E", "F", "G", "H"}, []int{headerStyle, dateStyle, moneyStyle})
				if err != nil {
					ch <- err
					return
				}
				log.Println("Done with receipts")
			}()

		case "lpo":
			sheetName := "Local Purchase Orders Sheet"
			_, err = f.NewSheet(sheetName)
			if err != nil {
				return nil, err
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				err = r.localPurchaseOrders(payload, f, sheetName, ctx, []string{"A", "B", "C", "D", "E"}, []int{headerStyle, dateStyle, moneyStyle})
				if err != nil {
					ch <- err
					return
				}
				log.Println("Done with lpo")
			}()
		}
	}

	go func() {
		wg.Wait()
		close(ch)
		log.Println("Done waiting")
	}()

	for err := range ch {
		if err != nil {
			return nil, err
		}
		log.Println("Done with error")
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	// if err := f.SaveAs("book2.xlsx"); err != nil {
	// 	fmt.Println(err)
	// }

	if err := f.Close(); err != nil {
		fmt.Println(err)
	}
	// return nil, nil

	return buf.Bytes(), nil
}
