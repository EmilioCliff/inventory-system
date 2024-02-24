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

const defaultFromName = "My Company Inc"
const defaultFromAddress = "My Company Address"
const defaultFromContact = "Company Contact"
const companyNo = "01234567890"
const fromName = "Emilio Cliff"
const fromAddress = "00100 - Nairobi"
const fromEmail = "company@gmail.com"
const fromContact = "1234567890"

// Input flags
// var invoiceNo string
// var invoiceDate string
// var toName string
// var toAddress string
// var toContact string
// var toEmail string

func SetInvoiceVariables(invoiceData map[string]string, data []map[string]interface{}) ([]byte, error) {
	userDetails := map[string]string{
		"invoiceNo":   fmt.Sprintf("INV - %v", invoiceData["invoice_number"]),
		"invoiceDate": invoiceData["created_at"],
		"toName":      invoiceData["invoice_username"],
		"toAddress":   data[0]["user_address"].(string),
		"toContact":   data[0]["user_contact"].(string),
		"toEmail":     data[0]["user_email"].(string),
	}

	var products [][]string
	for _, entry := range data {
		if _, isUserData := entry["user_contact"].(string); !isUserData {

			product := []string{
				fmt.Sprintf("%v", entry["totalBill"]),
				fmt.Sprintf("%v", entry["productName"]),
				fmt.Sprintf("%v", entry["productQuantity"]),
			}
			products = append(products, product)
		}
	}
	// fmt.Println(products)
	pdfBytes, err := generateInvoice(products, userDetails)
	if err != nil {
		return nil, err
	}
	return pdfBytes, err
}

func generateInvoice(data [][]string, user map[string]string) ([]byte, error) {
	marginX := 10.0
	marginY := 20.0
	gapY := 2.0
	pdf := fpdf.New("P", "mm", "A4", "") // 210mm x 297mm
	pdf.SetMargins(marginX, marginY, marginX)
	pdf.AddPage()
	pageW, _ := pdf.GetPageSize()
	safeAreaW := pageW - 2*marginX

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil, err
	}

	imagePath := filepath.Join(currentDir, "logi.png")

	pdf.ImageOptions(imagePath, 0, 0, 50, 40, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	pdf.SetFont("Arial", "B", 16)
	_, lineHeight := pdf.GetFontSize()
	currentY := pdf.GetY() + lineHeight + gapY + 16
	pdf.SetXY(marginX, currentY)

	pdf.Cell(40, 10, defaultFromName)

	if companyNo != "" {
		pdf.SetFont("Arial", "BI", 12)
		_, lineHeight = pdf.GetFontSize()
		pdf.SetXY(marginX, pdf.GetY()+lineHeight+gapY)
		pdf.Cell(40, 10, fmt.Sprintf("Company No : %v", defaultFromContact))
	}

	leftY := pdf.GetY() + lineHeight + gapY
	// Build invoice word on right
	pdf.SetFont("Arial", "B", 32)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetXY(130, currentY)
	pdf.Cell(100, 40, "INVOICE")

	newY := leftY
	if (pdf.GetY() + gapY) > newY {
		newY = pdf.GetY() + gapY
	}

	newY += 10.0 // Add margin

	pdf.SetXY(marginX, newY)
	pdf.SetFont("Arial", "", 12)
	_, lineHeight = pdf.GetFontSize()
	lineBreak := lineHeight + float64(1)

	pdf.SetFontStyle("B")
	pdf.Cell(safeAreaW/2, lineHeight, fromName)
	pdf.Ln(lineBreak)

	pdf.SetFontStyle("")
	pdf.Cell(safeAreaW/2, lineHeight, fromAddress)
	pdf.Ln(lineBreak)

	pdf.Cell(safeAreaW/2, lineHeight, "Kenya")
	pdf.Ln(lineBreak)

	pdf.Cell(safeAreaW/2, lineHeight, fromEmail)
	pdf.Ln(lineBreak)

	pdf.SetFontStyle("I")
	pdf.Cell(safeAreaW/2, lineHeight, fmt.Sprintf("Tel: %s", fromContact))
	pdf.Ln(lineBreak)
	pdf.Ln(lineBreak)
	pdf.Ln(lineBreak)

	pdf.SetFontStyle("B")
	pdf.Cell(safeAreaW/2, lineHeight, "Invoice To:")
	pdf.Line(marginX, pdf.GetY()+lineHeight, marginX+safeAreaW/2, pdf.GetY()+lineHeight)
	pdf.Ln(lineBreak)
	pdf.Cell(safeAreaW/2, lineHeight, user["toName"])

	pdf.SetFontStyle("")
	pdf.Ln(lineBreak)

	pdf.Cell(safeAreaW/2, lineHeight, user["toAddress"])
	pdf.Ln(lineBreak)

	pdf.Cell(safeAreaW/2, lineHeight, "Kenya")
	pdf.Ln(lineBreak)

	pdf.Cell(safeAreaW/2, lineHeight, user["toEmail"])
	pdf.Ln(lineBreak)

	pdf.SetFontStyle("I")
	pdf.Cell(safeAreaW/2, lineHeight, fmt.Sprintf("Tel: %s", user["toContact"]))

	endOfInvoiceDetailY := pdf.GetY() + lineHeight
	pdf.SetFontStyle("")

	// Right hand side info, invoice no & invoice date
	invoiceDetailW := float64(30)
	pdf.SetXY(safeAreaW/2+30, newY)
	pdf.Cell(invoiceDetailW, lineHeight, "Invoice No.:")
	pdf.Cell(invoiceDetailW, lineHeight, user["invoiceNo"])
	pdf.Ln(lineBreak)
	pdf.SetX(safeAreaW/2 + 30)
	pdf.Cell(invoiceDetailW, lineHeight, "Invoice Date:")
	pdf.Cell(invoiceDetailW, lineHeight, user["invoiceDate"])
	pdf.Ln(lineBreak)

	// Draw the table
	pdf.SetXY(marginX, endOfInvoiceDetailY+10.0)
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

	pdf.SetFontStyle("")
	pdf.Ln(lineBreak)
	pdf.Cell(safeAreaW, lineHeight, "Your satisfaction is our priority. If you have any concerns, please let us know.")

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	pdf.Close()

	return buffer.Bytes(), nil
}
