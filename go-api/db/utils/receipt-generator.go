package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-pdf/fpdf"
)

func SetReceiptVariables(receiptData map[string]string, data []map[string]interface{}) ([]byte, error) {
	userDetails := map[string]string{
		"receiptNo":   fmt.Sprintf("INV - %v", receiptData["receipt_number"]),
		"receiptDate": receiptData["created_at"],
		"toName":      receiptData["user_username"],
		"toAddress":   receiptData["user_address"],
		"toContact":   receiptData["user_contact"],
		"toEmail":     receiptData["user_email"],
	}

	if receiptData["by_admin"] == "true" {
		userDetails["byAdmin"] = "true"
	}

	var products [][]string
	for _, entry := range data {
		product := []string{
			fmt.Sprintf("%v", entry["totalBill"]),
			fmt.Sprintf("%v", entry["productName"]),
			fmt.Sprintf("%v", entry["productQuantity"]),
			fmt.Sprintf("%v", entry["productPrice"]),
		}
		products = append(products, product)
	}
	// fmt.Println(products)
	pdfBytes, err := generateReceipt(userDetails, products)
	if err != nil {
		return nil, err
	}
	return pdfBytes, err
}

func generateReceipt(user map[string]string, data [][]string) ([]byte, error) {
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
		pdf.SetFont("Arial", "I", 8)
		pdf.Cell(marginX+10, 10, "Your satisfaction is our priority. If you have any concerns, please let us know.")
		pdf.SetX(-marginX)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/%d", pdf.PageNo(), pdf.PageCount()), "", 0, "C", false, 0, "")
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
	pdf.Cell(70, 35, "RECEIPT")
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
	pdf.Ln(5)
	// pdf.SetFontStyle("I")
	pdf.Cell(40, lineHeight, "Phone: 0713851482")
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
	pdf.Cell(40, lineHeight, fmt.Sprintf("Phone: %s", user["toContact"]))

	pdf.SetY(headerY)
	pdf.SetFont("Arial", "B", 16)
	_, lineHeight = pdf.GetFontSize()
	middleX := safeAreaW/2 + 10
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
		quantity, _ := strconv.Atoi(val[2])
		unitPrice, _ := strconv.Atoi(val[3])
		totalBill, _ := strconv.ParseFloat(val[0], 64)
		subtotal += totalBill

		editedQuantity := FormatNumberWithCommas(int64(quantity))
		editedUnitPrice := FormatNumberWithCommas(int64(unitPrice))
		editedTotalBill := FormatNumberWithCommas(int64(totalBill))

		pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", rowJ+1), "1", 0, "CM", true, 0, "")             // No
		pdf.CellFormat(colWidth[1], lineHt, val[1], "1", 0, "LM", true, 0, "")                                // Description
		pdf.CellFormat(colWidth[2], lineHt, fmt.Sprintf("%s", editedQuantity), "1", 0, "CM", true, 0, "")     // Quantity
		pdf.CellFormat(colWidth[3], lineHt, fmt.Sprintf("%s.00", editedUnitPrice), "1", 0, "CM", true, 0, "") // Unit Price
		pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%s.00", editedTotalBill), "1", 0, "CM", true, 0, "") // Total Bill
		pdf.Ln(-1)
	}

	x, y = pdf.GetXY()

	// add something like kaNote
	if user["byAdmin"] == "true" {
		pdf.Ln(5)
		pdf.SetFontSize(8)
		pdf.SetFontStyle("I")
		pdf.Cell(0, 0, "Payment By Admin")
		pdf.Ln(5)
	}

	pdf.SetXY(x, y)
	pdf.SetFontSize(12)

	leftIndent := 0.0
	for i := 0; i < 3; i++ {
		leftIndent += colWidth[i]
	}

	pdf.SetFontStyle("B")

	editedSubtotal := FormatNumberWithCommas(int64(subtotal))

	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Subtotal", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%s.00", editedSubtotal), "1", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Grand total", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%s.00", editedSubtotal), "1", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	pdf.Close()

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
