package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-pdf/fpdf"
)

func GeneratePurchaseOrder(data map[string]interface{}) ([]byte, error) {
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
	pdf.SetMargins(marginX, marginX, marginX)
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

	pdf.SetFont("Arial", "B", 28)
	pdf.Cell(70, 35, "PURCHASE ORDER")
	pdf.Ln(32)

	pdf.SetFont("Arial", "B", 16)
	_, lineHeight := pdf.GetFontSize()
	currentY := pdf.GetY() - gapY - 5
	pdf.SetXY(marginX, currentY)

	pdf.SetFont("Arial", "B", 13)
	_, lineHeight = pdf.GetFontSize()
	pdf.Cell(40, lineHeight, "RG Center, Ground")
	pdf.Ln(lineHeight)
	pdf.Cell(40, lineHeight, "Floor Room A10,")
	pdf.Ln(lineHeight)
	pdf.Cell(40, lineHeight, "Eastern Bypass")
	pdf.Ln(lineHeight)
	pdf.Cell(40, lineHeight, "Road Utawala,")
	pdf.Ln(lineHeight)
	pdf.Cell(40, lineHeight, "Nairobi.")
	pdf.Ln(lineHeight)
	pdf.Cell(40, lineHeight, "kokomed421@gmail.com")
	pdf.Ln(lineHeight)
	pdf.Ln(5)
	pdf.Cell(40, lineHeight, "Phone: 0713851482")
	pdf.Ln(lineHeight * 3)

	headerY := pdf.GetY()

	pdf.SetFont("Arial", "B", 16)
	_, lineHeight = pdf.GetFontSize()
	pdf.SetXY(marginX, headerY)
	pdf.Cell(40, lineHeight, "PURCHASE ORDER TO")
	pdf.SetFont("Arial", "", 12)
	_, lineHeight = pdf.GetFontSize()
	pdf.Ln(lineHeight + gapY + 2)
	pdf.Cell(40, lineHeight, strings.ToUpper(data["supplier_name"].(string)))
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, strings.ToUpper(fmt.Sprintf("P.O Box %s", data["po_box"])))
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, strings.ToUpper(data["address"].(string)))

	pdf.SetY(headerY)
	pdf.SetFont("Arial", "B", 16)
	_, lineHeight = pdf.GetFontSize()
	middleX := safeAreaW / 2
	pdf.SetX(-middleX)
	pdf.Cell(40, lineHeight, "PURCHASE ORDER")
	pdf.Ln(lineHeight + gapY + 2)
	pdf.SetX(-middleX)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(42, lineHeight, "ORDER NO: ")
	pdf.Cell(40, lineHeight, data["order_number"].(string))
	pdf.Ln(lineHeight + gapY)
	pdf.SetX(-middleX)
	pdf.Cell(42, lineHeight, "ORDER DATE:")
	pdf.Cell(40, lineHeight, data["order_date"].(string))
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
	var subtotal int64

	openData := data["Data"].([]map[string]interface{})

	for rowJ := 0; rowJ < len(openData); rowJ++ {
		product := openData[rowJ]
		totalBill := int64(product["quantity"].(float64) * product["unit_price"].(float64))
		editedTotalBill := FormatNumberWithCommas(totalBill)
		subtotal += totalBill

		pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", rowJ+1), "1", 0, "CM", true, 0, "")                            // No
		pdf.CellFormat(colWidth[1], lineHt, product["product_name"].(string), "1", 0, "LM", true, 0, "")                     // Description
		pdf.CellFormat(colWidth[2], lineHt, fmt.Sprintf("%.0f", product["quantity"].(float64)), "1", 0, "CM", true, 0, "")   // Quantity
		pdf.CellFormat(colWidth[3], lineHt, fmt.Sprintf("%.2f", product["unit_price"].(float64)), "1", 0, "CM", true, 0, "") // Unit Price
		pdf.SetFontStyle("B")
		pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%s.00", editedTotalBill), "1", 0, "CM", true, 0, "") // Total Bill
		pdf.SetFontStyle("")
		pdf.Ln(-1)
	}

	leftIndent := 0.0
	for i := 0; i < 3; i++ {
		leftIndent += colWidth[i]
	}

	// Calculate the total
	pdf.SetFontStyle("B")

	editedTotalBill := FormatNumberWithCommas(subtotal)

	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Total", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%s.00", editedTotalBill), "1", 0, "CM", true, 0, "")
	pdf.Ln(30)

	pdf.Cell(40, 8.0, "COMMENTS:")

	pdf.Ln(30)

	pdf.Cell(40, 8.0, "PREPARED BY:")

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	pdf.Close()

	return buffer.Bytes(), nil
}

func FormatNumberWithCommas(num int64) string {
	str := fmt.Sprintf("%d", num)
	n := len(str)

	if n <= 3 {
		return str
	}

	var result strings.Builder
	mod := n % 3

	// If there are leading digits before the first comma
	if mod > 0 {
		result.WriteString(str[:mod])
		result.WriteString(",")
	}

	// Add the rest of the digits in groups of three
	for i := mod; i < n; i += 3 {
		if i > mod {
			result.WriteString(str[i : i+3])
		} else {
			result.WriteString(str[i : i+3])
		}
		if i+3 < n {
			result.WriteString(",")
		}
	}

	return result.String()
}

// func FormatFloatNumberWithCommas(num float64) string {
// 	// Format the number to have two decimal places
// 	formattedNumber := fmt.Sprintf("%.2f", num)

// 	// Split the number into integer and decimal parts
// 	parts := strings.Split(formattedNumber, ".")
// 	integerPart := parts[0]
// 	decimalPart := parts[1]

// 	// Format the integer part with commas
// 	n := len(integerPart)
// 	if n <= 3 {
// 		return integerPart + "." + decimalPart
// 	}

// 	var result strings.Builder
// 	mod := n % 3

// 	// If there are leading digits before the first comma
// 	if mod > 0 {
// 		result.WriteString(integerPart[:mod])
// 		result.WriteString(",")
// 	}

// 	// Add the rest of the digits in groups of three
// 	for i := mod; i < n; i += 3 {
// 		if i > mod {
// 			result.WriteString(integerPart[i : i+3])
// 		} else {
// 			result.WriteString(integerPart[i : i+3])
// 		}
// 		if i+3 < n {
// 			result.WriteString(",")
// 		}
// 	}

// 	return result.String() + "." + decimalPart
// }
