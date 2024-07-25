package utils

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-pdf/fpdf"
)

func SetReceiptVariables(receipt map[string]string, data []map[string]interface{}) ([]byte, error) {
	userDetails := map[string]string{
		"receiptNo":   fmt.Sprintf("RCT - %v", receipt["receipt_number"]),
		"receiptDate": receipt["created_at"],
		"toName":      receipt["receipt_username"],
		"toAddress":   data[0]["user_address"].(string),
		"toContact":   data[0]["user_contact"].(string),
		"toEmail":     data[0]["user_email"].(string),
	}

	var products [][]string
	for _, entry := range data {
		if _, ok := entry["user_contact"].(string); !ok {

			product := []string{
				fmt.Sprintf("%v", entry["totalBill"]),
				fmt.Sprintf("%v", entry["productName"]),
				fmt.Sprintf("%v", entry["productQuantity"]),
			}
			products = append(products, product)
		}
	}
	pdfBytes, err := generateReceipt(products, userDetails)
	if err != nil {
		return nil, err
	}
	return pdfBytes, nil
}

func generateReceipt(data [][]string, user map[string]string) ([]byte, error) {
	marginX := 10.0
	marginY := 20.0
	gapY := 2.0
	pdf := fpdf.New("P", "mm", "A4", "") // 210mm x 297mm
	pdf.SetMargins(marginX, marginY, marginX)
	pdf.SetFooterFunc(func() {
		pdf.SetFont("Arial", "I", 12)
		pdf.SetXY((190/2)-20, -36)
		pdf.CellFormat(40, 5, "PAYMENT METHOD", "", 0, "CM", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX((190 / 2) - 20)
		pdf.CellFormat(40, 5, "BUYGOODS", "", 0, "CM", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX((190 / 2) - 20)
		pdf.CellFormat(40, 5, "TILL NO: 9090757", "", 0, "CM", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX((190 / 2) - 20)
		pdf.CellFormat(40, 5, "Inventory System", "", 0, "CM", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX((190 / 2) - 20)

		pdf.SetXY(marginX, -15)
		pdf.SetFont("Arial", "I", 8)
		pdf.Cell(marginX+10, 10, "Your satisfaction is our priority. If you have any concerns, please let us know.")
		pdf.SetX(-marginX)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/%d", pdf.PageNo(), pdf.PageCount()), "", 0, "C", false, 0, "")
	})
	pdf.AddPage()
	pageW, _ := pdf.GetPageSize()
	safeAreaW := pageW - 2*marginX

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil, err
	}

	imagePath := filepath.Join(currentDir, "Kokomed-Logo.png")

	pdf.SetXY(-(marginX + 110), marginY+40)
	pdf.ImageOptions(imagePath, -40, -(marginY - 10), 120, 110, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdf.SetXY(marginX, marginY-15)
	pdf.SetFont("Arial", "B", 32)
	pdf.Cell(70, 35, "RECEIPT")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "B", 16)
	_, lineHeight := pdf.GetFontSize()
	currentY := pdf.GetY() - gapY - 5
	pdf.SetXY(marginX, currentY)

	pdf.SetFont("Arial", "", 12)
	_, lineHeight = pdf.GetFontSize()
	pdf.Cell(40, lineHeight, "RG Center, Ground Floor")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "Room A10")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "Eastern Bypass Road")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "Utawala, Nairobi")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "kokomed421@gmail.com")
	pdf.Ln(lineHeight + gapY)
	pdf.SetFontStyle("I")
	pdf.Cell(40, lineHeight, "Tel: 0713851482")
	pdf.Ln(lineHeight * 3)

	headerY := pdf.GetY()

	pdf.SetFont("Arial", "B", 16)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetXY(marginX, headerY)
	pdf.Cell(40, lineHeight, "RECEIPT TO")
	pdf.SetFont("Arial", "", 12)
	_, lineHeight = pdf.GetFontSize()
	pdf.Ln(lineHeight + gapY + 2)
	pdf.Cell(40, lineHeight, user["toName"])
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, user["toAddress"])
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, user["toEmail"])
	pdf.Ln(lineHeight + gapY)
	pdf.SetFontStyle("I")
	pdf.Cell(40, lineHeight, fmt.Sprintf("Tel: %s", user["toContact"]))

	pdf.SetY(headerY)
	pdf.SetFont("Arial", "B", 16)
	_, lineHeight = pdf.GetFontSize()
	middleX := safeAreaW / 2
	pdf.SetX(-middleX)
	pdf.Cell(40, lineHeight, "RECEIPT")
	pdf.Ln(lineHeight + gapY + 2)
	pdf.SetX(-middleX)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(42, lineHeight, "RECEIPT NO: ")
	pdf.Cell(40, lineHeight, user["receiptNo"])
	pdf.Ln(lineHeight + gapY)
	pdf.SetX(-middleX)
	pdf.Cell(42, lineHeight, "RECEIPT DATE:")
	pdf.Cell(40, lineHeight, user["receiptDate"])
	pdf.Ln(lineHeight * 4)

	// Draw the table
	pdf.SetX(marginX)
	lineHt := 10.0
	const colNumber = 5
	header := [colNumber]string{"No", "Description", "Quantity", "Unit Price (Ksh)", "Price (Ksh)"}
	colWidth := [colNumber]float64{10.0, 75.0, 25.0, 40.0, 40.0}

	// Headers
	pdf.SetFontStyle("B")
	pdf.SetFillColor(200, 200, 200)
	for colJ := 0; colJ < colNumber; colJ++ {
		pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "CM", true, 0, "")
	}

	pdf.Ln(-1)
	pdf.SetFillColor(255, 255, 255)

	// Table data
	pdf.SetFontStyle("")
	subtotal := 0.0

	for rowJ := 0; rowJ < len(data); rowJ++ {
		val := data[rowJ]
		if len(val) == 3 {
			// Column 1: Unit
			// Column 2: Description
			// Column 3: Price per unit
			unit, _ := strconv.Atoi(val[2])
			desc := val[1]
			quantity, _ := strconv.ParseFloat(val[2], 64)
			totalBill, _ := strconv.ParseFloat(val[0], 64)
			unitPrice := math.Round((totalBill/quantity)*100) / 100
			subtotal += totalBill

			pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", rowJ+1), "1", 0, "CM", true, 0, "")      // No
			pdf.CellFormat(colWidth[1], lineHt, desc, "1", 0, "LM", true, 0, "")                           // Description
			pdf.CellFormat(colWidth[2], lineHt, fmt.Sprintf("%d", unit), "1", 0, "CM", true, 0, "")        // Quantity
			pdf.CellFormat(colWidth[3], lineHt, fmt.Sprintf("%.2f", unitPrice), "1", 0, "CM", true, 0, "") // Unit Price
			pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%.2f", totalBill), "1", 0, "CM", true, 0, "") // Total Bill
			pdf.Ln(-1)
		}
	}

	xPos, yPos := pdf.GetXY()
	pdf.Ln(5)
	_, height := pdf.GetFontSize()

	pdf.SetFont("Arial", "", 8.0)
	pdf.MultiCell(75.0, height, fmt.Sprintf("Payment made from %s. This receipt has been generated electronically and does not require a signature.", user["toContact"]), "0", "", false)
	pdf.SetFont("Arial", "", 12)
	pdf.SetXY(xPos, yPos)

	// Calculate the subtotal
	pdf.SetFontStyle("B")
	leftIndent := 0.0
	for i := 0; i < 3; i++ {
		leftIndent += colWidth[i]
	}
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Subtotal", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%.2f", subtotal), "1", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	grandTotal := subtotal
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Grand total", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%.2f", grandTotal), "1", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	// const colNum = 5
	// headers := [colNum]string{"No", "Product", "Quantity", "Unit Price (Ksh)", "Total Price (Ksh)"}
	// colW := [colNum]float64{75.0, 25.0, 40.0, 40.0}

	// pdf.SetX(marginX)
	// pdf.SetFontStyle("B")
	// pdf.SetFillColor(200, 200, 200)
	// for i := 0; i < colNum; i++ {
	// 	pdf.CellFormat(colW[i], 10, headers[i], "", 0, "CM", true, 0, "")
	// }

	// pdf.Ln(-1)
	// pdf.SetFillColor(255, 255, 255)

	// pdf.SetFontStyle("")
	// subtotal := 0.0

	// for rowJ := 0; rowJ < len(data); rowJ++ {
	// 	val := data[rowJ]
	// 	if len(val) == 3 {
	// 		unit, _ := strconv.Atoi(val[2])
	// 		desc := val[1]
	// 		quantity, _ := strconv.ParseFloat(val[2], 64)
	// 		totalBill, _ := strconv.ParseFloat(val[0], 64)
	// 		unitPrice := math.Round((totalBill/quantity)*100) / 100
	// 		subtotal += totalBill

	// 		pdf.CellFormat(colW[0], 10, fmt.Sprintf("%d", rowJ+1), "", 0, "CM", true, 0, "")       // No
	// 		pdf.CellFormat(colW[0], 10, desc, "B", 0, "LM", true, 0, "")                           // Description
	// 		pdf.CellFormat(colW[1], 10, fmt.Sprintf("%d", unit), "B", 0, "CM", true, 0, "")        // Quantity
	// 		pdf.CellFormat(colW[2], 10, fmt.Sprintf("%.2f", unitPrice), "B", 0, "CM", true, 0, "") // Unit Price
	// 		pdf.CellFormat(colW[3], 10, fmt.Sprintf("%.2f", totalBill), "B", 0, "CM", true, 0, "") // Total Bill
	// 		pdf.Ln(-1)
	// 	}
	// }

	// // Calculate the subtotal
	// pdf.SetFontStyle("B")
	// leftIndent := 0.0
	// for i := 0; i < 3; i++ {
	// 	leftIndent += colW[i]
	// }
	// pdf.SetX(marginX + leftIndent)
	// pdf.CellFormat(colW[2], 10, "Subtotal", "B", 0, "CM", true, 0, "")
	// pdf.CellFormat(colW[3], 10, fmt.Sprintf("%.2f", subtotal), "B", 0, "CM", true, 0, "")
	// pdf.Ln(-1)

	// grandTotal := subtotal
	// pdf.SetX(marginX + leftIndent)
	// pdf.CellFormat(colW[2], 10, "Grand total", "B", 0, "CM", true, 0, "")
	// pdf.CellFormat(colW[3], 10, fmt.Sprintf("%.2f", grandTotal), "B", 0, "CM", true, 0, "")
	// pdf.Ln(-1)

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	pdf.Close()

	return buffer.Bytes(), nil
}
