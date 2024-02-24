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

// var invoiceNo string
// var invoiceDate string
// var toName string
// var toAddress string
// var toContact string
// var toEmail string

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
	// fmt.Println(products)
	pdfBytes, err := generateReceipt(products, userDetails)
	if err != nil {
		return nil, err
	}
	return pdfBytes, nil
}

func generateReceipt(data [][]string, user map[string]string) ([]byte, error) {
	marginX := 10.0
	marginY := 10.0
	gapY := 1.0
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginX, marginY, marginX)

	pdf.SetHeaderFuncMode(func() {
		pdf.SetFont("Arial", "B", 15)
		pdf.Cell(80, 0, "")
		pdf.CellFormat(30, 10, "Payment Receipt", "", 0, "C", false, 0, "")
		pdf.Ln(2)
	}, true)

	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.Cell(marginX+10, 10, "Your satisfaction is our priority. If you have any concerns, please let us know.")
		pdf.SetX(-15)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/%d", pdf.PageNo(), pdf.PageCount()), "", 0, "C", false, 0, "")
	})
	pdf.AddPage()
	pageW, _ := pdf.GetPageSize()
	safeW := pageW - 2*marginX
	centerW := safeW/2 - 10

	pdf.SetXY(safeW/2-20, marginY+25)
	pdf.SetFont("Arial", "B", 18)
	_, lineHeight := pdf.GetFontSize()
	pdf.Cell(40, 10, "PAYMENT RECEIPT")
	pdf.Ln(lineHeight)
	pdf.Ln(lineHeight)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil, err
	}

	imagePath := filepath.Join(currentDir, "logi.png")

	pdf.ImageOptions(imagePath, safeW/2, pdf.GetY()+lineHeight*3, 30, 20, true, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdf.SetFont("Arial", "", 12)
	_, lineHeight = pdf.GetFontSize()
	currentY := pdf.GetY()
	// pdf.SetXY(safeW/2 - 30, currentY+lineHeight*3)
	// pdf.Cell(safeW/2, lineHeight, defaultFromAddress)

	// currentY = pdf.GetY()
	pdf.SetFontStyle("B")
	pdf.SetXY(centerW, currentY+lineHeight*3)
	pdf.Cell(50, lineHeight, defaultFromName)
	pdf.Ln(lineHeight + gapY)

	pdf.SetX(centerW)
	pdf.SetFontStyle("")
	pdf.Cell(50, lineHeight, defaultFromAddress)
	pdf.Ln(lineHeight + gapY)

	pdf.SetX(centerW)
	pdf.SetFontStyle("")
	pdf.Cell(50, lineHeight, defaultFromContact)
	pdf.Ln(lineHeight)
	// pdf.Ln(lineHeight)

	pdf.Line(marginX, pdf.GetY()+lineHeight, marginX+safeW, pdf.GetY()+lineHeight)
	pdf.Ln(lineHeight)
	pdf.Ln(lineHeight)

	pdf.SetX(centerW)
	pdf.Cell(40, lineHeight, user["toName"])
	pdf.Ln(lineHeight + gapY)

	pdf.SetX(centerW)
	pdf.Cell(40, lineHeight, user["toContact"])
	pdf.Ln(lineHeight + gapY)

	pdf.SetX(centerW)
	pdf.Cell(40, lineHeight, user["toEmail"])
	pdf.Ln(lineHeight + gapY)
	pdf.Ln(lineHeight + gapY)
	pdf.Ln(lineHeight + gapY)

	pdf.SetX(centerW)
	pdf.SetFontStyle("B")
	pdf.Cell(30, lineHeight, "Receipt No: ")
	pdf.SetFontStyle("")
	pdf.Cell(30, lineHeight, user["receiptNo"])
	pdf.Ln(lineHeight + gapY)

	pdf.SetX(centerW)
	pdf.SetFontStyle("B")
	pdf.Cell(30, lineHeight, "Receipt Date: ")
	pdf.SetFontStyle("")
	pdf.Cell(30, lineHeight, user["receiptDate"])
	pdf.Ln(lineHeight + gapY)
	pdf.Ln(lineHeight + gapY)

	const colNum = 4
	headers := [colNum]string{"Product", "Quantity", "Unit Price (Ksh)", "Total Price (Ksh)"}
	colW := [colNum]float64{75.0, 25.0, 40.0, 40.0}

	pdf.SetX(marginX)
	pdf.SetFontStyle("B")
	pdf.SetFillColor(200, 200, 200)
	for i := 0; i < colNum; i++ {
		pdf.CellFormat(colW[i], 10, headers[i], "", 0, "CM", true, 0, "")
	}

	pdf.Ln(-1)
	pdf.SetFillColor(255, 255, 255)

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

			// pdf.CellFormat(colW[0], 10, fmt.Sprintf("%d", rowJ+1), "", 0, "CM", true, 0, "")      // No
			pdf.CellFormat(colW[0], 10, desc, "B", 0, "LM", true, 0, "")                           // Description
			pdf.CellFormat(colW[1], 10, fmt.Sprintf("%d", unit), "B", 0, "CM", true, 0, "")        // Quantity
			pdf.CellFormat(colW[2], 10, fmt.Sprintf("%.2f", unitPrice), "B", 0, "CM", true, 0, "") // Unit Price
			pdf.CellFormat(colW[3], 10, fmt.Sprintf("%.2f", totalBill), "B", 0, "CM", true, 0, "") // Total Bill
			pdf.Ln(-1)
		}
	}

	// Calculate the subtotal
	pdf.SetFontStyle("B")
	leftIndent := 0.0
	for i := 0; i < 2; i++ {
		leftIndent += colW[i]
	}
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colW[2], 10, "Subtotal", "B", 0, "CM", true, 0, "")
	pdf.CellFormat(colW[3], 10, fmt.Sprintf("%.2f", subtotal), "B", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	grandTotal := subtotal
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colW[2], 10, "Grand total", "B", 0, "CM", true, 0, "")
	pdf.CellFormat(colW[3], 10, fmt.Sprintf("%.2f", grandTotal), "B", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	// pdf.SetFontStyle("")
	// pdf.Ln(lineBreak)
	// pdf.Cell(safeAreaW, lineHeight, "Your satisfaction is our priority. If you have any concerns, please let us know.")

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	pdf.Close()

	return buffer.Bytes(), nil
}
