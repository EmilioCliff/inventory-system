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
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil, err
	}

	imagePath := filepath.Join(currentDir, "Kokomed-Logo-small.png")

	marginX := 10.0
	marginY := 20.0
	gapY := 2.0
	pdf := fpdf.New("P", "mm", "A4", "") // 210mm x 297mm
	pdf.SetMargins(marginX, marginY, marginX)
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetX(marginX)
		pdf.SetFont("Arial", "I", 8)
		pdf.Cell(marginX+10, 10, "Your satisfaction is our priority. If you have any concerns, please let us know.")
		pdf.SetX(-marginX)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	pdf.AddPage()
	pageW, _ := pdf.GetPageSize()
	safeAreaW := pageW - 2*marginX

	pdf.SetXY(0, 0)
	x, y := pdf.GetXY()

	pdf.SetFillColor(97, 201, 221)
	pdf.Rect(x, y, 210, 40, "F")
	pdf.Ln(7)

	pdf.SetFillColor(255, 255, 255)

	pdf.SetXY(-(marginX + 110), marginY+50)
	pdf.ImageOptions(imagePath, marginX, 5, 35, 30, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdf.SetFont("Arial", "B", 6)
	pdf.SetXY(marginX+3, 19)
	pdf.Cell(70, 35, "Bridging Your Healthcare")

	pdf.SetXY(marginX+35, 5)
	pdf.SetFont("Arial", "B", 30)
	pdf.SetTextColor(10, 47, 70)
	pdf.CellFormat(70, 30, "KOKOMED SUPPLIES LTD", "", 0, "T", false, 0, "")
	pdf.Ln(10)
	pdf.SetX(marginX + 50)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(10)
	pdf.Cell(70, 8, "Laboratory supplies  /  Medical supplies  / Medical devices")
	pdf.Ln(5)
	pdf.SetX(marginX + 45)
	pdf.Cell(70, 8, "/ Hospital supplies  /  Hospital equipment  /  Laboratory instrument")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 32)
	pdf.Cell(70, 35, "INVOICE")
	pdf.Ln(31)

	pdf.SetFont("Arial", "B", 16)
	_, lineHeight := pdf.GetFontSize()
	currentY := pdf.GetY() - gapY - 5
	pdf.SetXY(marginX, currentY)

	pdf.SetFont("Arial", "", 12)
	_, lineHeight = pdf.GetFontSize()
	pdf.Cell(40, lineHeight, "RG Center, Ground")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "Floor Room A10,")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "Eastern Bypass")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "Road Utawala,")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "Nairobi.")
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, "kokomed421@gmail.com")
	pdf.Ln(lineHeight + gapY)
	// pdf.SetFontStyle("I")
	pdf.Cell(40, lineHeight, "Phone: 0713851482")
	pdf.Ln(lineHeight * 3)

	headerY := pdf.GetY()

	pdf.SetFont("Arial", "B", 16)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetXY(marginX, headerY)
	pdf.Cell(40, lineHeight, "INVOICE TO")
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
	pdf.Cell(40, lineHeight, "INVOICE")
	pdf.Ln(lineHeight + gapY + 2)
	pdf.SetX(-middleX)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(42, lineHeight, "INVOICE NO: ")
	pdf.Cell(40, lineHeight, user["invoiceNo"])
	pdf.Ln(lineHeight + gapY)
	pdf.SetX(-middleX)
	pdf.Cell(42, lineHeight, "INVOICE DATE:")
	pdf.Cell(40, lineHeight, user["invoiceDate"])
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

			editedUnitPrice := FormatNumberWithCommas(int64(unitPrice))

			editedTotalBill := FormatNumberWithCommas(int64(totalBill))

			pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", rowJ+1), "1", 0, "CM", true, 0, "")             // No
			pdf.CellFormat(colWidth[1], lineHt, desc, "1", 0, "LM", true, 0, "")                                  // Description
			pdf.CellFormat(colWidth[2], lineHt, fmt.Sprintf("%d", unit), "1", 0, "CM", true, 0, "")               // Quantity
			pdf.CellFormat(colWidth[3], lineHt, fmt.Sprintf("%s.00", editedUnitPrice), "1", 0, "CM", true, 0, "") // Unit Price
			pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%s.00", editedTotalBill), "1", 0, "CM", true, 0, "") // Total Bill
			pdf.Ln(-1)
		}
	}

	xPos, yPos := pdf.GetXY()

	leftIndent := 0.0
	for i := 0; i < 3; i++ {
		leftIndent += colWidth[i]
	}

	// Calculate the subtotal
	pdf.SetFontStyle("B")

	editedSubtotal := FormatNumberWithCommas(int64(subtotal))

	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Subtotal", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%s.00", editedSubtotal), "1", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	// grandTotal := subtotal
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Grand total", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%s.00", editedSubtotal), "1", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	pdf.SetXY(xPos, yPos)

	pdf.Ln(2)

	pdf.SetFont("Times", "B", 12.0)
	pdf.Cell(40, 8.0, "PAYMENT INFORMATION")
	pdf.SetFontStyle("")
	pdf.Ln(4.5)
	pdf.Cell(60, 8, "Bank Name: I&M BANK LTD")
	pdf.Ln(4.5)
	pdf.Cell(60, 8, "Account Name: KOKOMED SUPPLIES LIMITED")
	pdf.Ln(4.5)
	pdf.Cell(60, 8, "Account Number: 00103333091850")
	pdf.Ln(6)

	// pdf.SetX(leftIndent / 2)
	pdf.SetFontStyle("B")
	pdf.Cell(60, 8, "OR")
	pdf.SetFontStyle("")

	pdf.Ln(6)
	pdf.Cell(60, 8.0, fmt.Sprintf("Safaricom Lipa na M-Pesa"))
	pdf.Ln(4.5)
	pdf.Cell(60, 8.0, fmt.Sprintf("Till Number: %s", "9090757"))

	pdf.Ln(30)

	pdf.SetFont("Times", "B", 12.0)
	pdf.Cell(40, 8.0, "TERMS AND CONDITIONS")
	pdf.SetFontStyle("")
	pdf.Ln(5)
	pdf.Cell(60, 8.0, fmt.Sprintf("Please make payment upon reciept of this invoice."))

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	pdf.Close()

	return buffer.Bytes(), nil
}
