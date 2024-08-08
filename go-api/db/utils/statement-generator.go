package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-pdf/fpdf"
)

func GenerateStatement(data []map[string]interface{}, userDetails map[string]string) ([]byte, error) {
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

	pdf.SetFillColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 16)
	_, lineHeight := pdf.GetFontSize()

	pdf.SetFont("Arial", "", 13)
	_, lineHeight = pdf.GetFontSize()
	pdf.Ln(lineHeight + gapY + 2)
	pdf.SetFontStyle("B")
	pdf.Cell(38, lineHeight, "Name: ")
	pdf.SetFontStyle("")
	pdf.Cell(40, lineHeight, fmt.Sprintf("%s", userDetails["username"]))
	pdf.Ln(lineHeight + gapY)
	pdf.SetFontStyle("B")
	pdf.Cell(38, lineHeight, "Phone Number: ")
	pdf.SetFontStyle("")
	pdf.Cell(40, lineHeight, fmt.Sprintf("%s", userDetails["phone"]))
	pdf.Ln(lineHeight + gapY)
	pdf.SetFontStyle("B")
	pdf.Cell(38, lineHeight, "Email Address:")
	pdf.SetFontStyle("")
	pdf.Cell(40, lineHeight, fmt.Sprintf("%s", userDetails["email"]))
	pdf.Ln(lineHeight + gapY)
	pdf.SetFontStyle("B")
	pdf.Cell(38, lineHeight, "Date: ")
	pdf.SetFontStyle("")
	pdf.Cell(40, lineHeight, fmt.Sprintf("%s", userDetails["date"]))
	pdf.Ln(lineHeight + gapY)

	pdf.Ln(8)
	pdf.SetX(safeAreaW/2 - 20)
	x, y = pdf.GetXY()
	pdf.SetTextColor(50, 127, 168)
	pdf.Cell(40, lineHeight, fmt.Sprintf("Transaction Details"))
	pdf.SetTextColor(0, 0, 0)

	textWidth := pdf.GetStringWidth("Transaction Details")

	// Draw the underline
	underlineY := y + 7 // position of the underline (adjust as needed)
	pdf.Line(x, underlineY, x+textWidth+2, underlineY)

	pdf.Ln(12)

	// Draw the table
	pdf.SetX(marginX)
	lineHt := 10.0
	const colNumber = 5
	header := [colNumber]string{"Receipt Number", "Mpesa Ref", "Amount (Ksh)", "Completion Time"}
	colWidth := [colNumber]float64{50.0, 50.0, 40.0, 50.0}

	// Headers
	pdf.SetFont("Arial", "", 12)
	pdf.SetFillColor(200, 200, 200)
	for colJ := 0; colJ < colNumber; colJ++ {
		pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "", 0, "CM", true, 0, "")
	}

	pdf.Ln(-1)
	pdf.SetFillColor(255, 255, 255)

	// Table data
	pdf.SetFontStyle("")

	for rowJ := 0; rowJ < len(data); rowJ++ {
		total := FormatNumberWithCommas(int64(data[rowJ]["amount"].(int32)))
		pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%v", data[rowJ]["receipt_number"]), "B", 0, "CM", true, 0, "")
		pdf.CellFormat(colWidth[1], lineHt, fmt.Sprintf("%v", data[rowJ]["mpesa_number"]), "B", 0, "CM", true, 0, "")
		pdf.CellFormat(colWidth[2], lineHt, fmt.Sprintf("%v.00", total), "B", 0, "CM", true, 0, "")
		pdf.CellFormat(colWidth[3], lineHt, fmt.Sprintf("%v", data[rowJ]["created_at"]), "B", 0, "CM", true, 0, "")
		pdf.Ln(-1)
	}

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	pdf.Close()

	return buffer.Bytes(), nil
}
