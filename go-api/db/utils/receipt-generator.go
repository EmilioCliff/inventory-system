package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-pdf/fpdf"
)

func GenerateReceipt(data map[string]string) ([]byte, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return nil, err
	}

	imagePath := filepath.Join(currentDir, "Kokomed-Logo-small.png")

	marginX := 3.0
	marginY := 3.0
	gapY := 2.0
	pdf := fpdf.New("P", "mm", "A5", "") // 148mm x 210mm
	pdf.SetMargins(marginX, marginY, marginX)
	pdf.SetFooterFunc(func() {
		pdf.SetXY(marginX, -15)
		pdf.SetFont("Arial", "I", 8)
		pdf.Cell(marginX+10, 10, "Your satisfaction is our priority. If you have any concerns, please let us know.")
		pdf.SetX(-(marginX + 18))
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/%d", pdf.PageNo(), pdf.PageCount()), "", 0, "C", false, 0, "")
	})
	pdf.AddPage()
	pageW, _ := pdf.GetPageSize()
	safeAreaW := pageW - 2*marginX

	pdf.SetXY(0, 0)
	x, y := pdf.GetXY()

	pdf.SetFillColor(97, 201, 221)
	pdf.Rect(x, y, 210, 36, "F")
	pdf.Ln(7)

	pdf.SetFillColor(255, 255, 255)

	pdf.SetXY(-(marginX + 110), marginY+50)
	pdf.ImageOptions(imagePath, marginX, 5, 27, 25, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdf.SetFont("Arial", "B", 6)
	pdf.SetXY(marginX, 15)
	pdf.Cell(70, 35, "Bridging Your Healthcare")

	pdf.SetXY(marginX+30, 5)
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(10, 47, 70)
	pdf.CellFormat(70, 30, "KOKOMED SUPPLIES LTD", "", 0, "T", false, 0, "")
	pdf.Ln(7)
	pdf.SetX(marginX + 35)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFontSize(8)
	pdf.Cell(70, 8, "Laboratory supplies  /  Medical supplies  / Medical devices")
	pdf.Ln(5)
	pdf.SetX(marginX + 30)
	pdf.Cell(70, 8, "/ Hospital supplies  /  Hospital equipment  /  Laboratory instrument")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 28)
	pdf.Cell(70, 35, "RECEIPT")
	pdf.Ln(32)

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
	pdf.Cell(40, lineHeight, data["receipt_username"])
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, data["user_address"])
	pdf.Ln(lineHeight + gapY)
	pdf.Cell(40, lineHeight, data["user_email"])
	pdf.Ln(lineHeight + gapY)
	// pdf.SetFontStyle("I")
	pdf.Cell(40, lineHeight, fmt.Sprintf("Phone: %s", data["user_contact"]))

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
	pdf.Cell(40, lineHeight, data["receipt_number"])
	pdf.Ln(lineHeight + gapY)
	pdf.SetX(-middleX)
	pdf.Cell(42, lineHeight, "RECEIPT DATE:")
	pdf.Cell(40, lineHeight, data["created_at"])
	pdf.Ln(lineHeight * 4)

	// Draw the table
	pdf.SetX(marginX)
	lineHt := 10.0
	// const colNumber = 4
	// header := [colNumber]string{"M-Pesa Receipt Number", "Phone Number", "Amount (Ksh)", "Status"}
	// colWidth := [colNumber]float64{60.0, 40.0, 40.0, 50.0}

	// Headers
	pdf.SetFontStyle("B")
	pdf.SetFillColor(200, 200, 200)
	// for colJ := 0; colJ < colNumber; colJ++ {
	// 	pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "CM", true, 0, "")
	// }
	pdf.CellFormat(50, lineHt, "M-Pesa Receipt Number", "1", 0, "LM", true, 0, "")
	pdf.CellFormat(40, lineHt, "Phone Number", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(30, lineHt, "Amount (Ksh)", "1", 0, "CM", true, 0, "")
	pdf.CellFormat(22, lineHt, "Status", "1", 0, "CM", true, 0, "")

	pdf.Ln(-1)
	pdf.SetFillColor(255, 255, 255)

	// Table data
	pdf.SetFontStyle("")

	intAmount, _ := strconv.Atoi(data["amount"])
	editedAmount := FormatNumberWithCommas(int64(intAmount))

	pdf.CellFormat(50, lineHt, data["mpesa_receipt_number"], "1", 0, "LM", true, 0, "")
	pdf.CellFormat(40, lineHt, data["user_contact"], "1", 0, "CM", true, 0, "")
	pdf.CellFormat(30, lineHt, fmt.Sprintf("%s.00", editedAmount), "1", 0, "CM", true, 0, "")
	pdf.CellFormat(22, lineHt, data["status"], "1", 0, "CM", true, 0, "")
	pdf.Ln(-1)

	var buffer bytes.Buffer
	if err := pdf.Output(&buffer); err != nil {
		return nil, err
	}

	pdf.Close()

	return buffer.Bytes(), nil
}
