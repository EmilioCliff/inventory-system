package db

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/stretchr/testify/require"
)

func createInvoiceTest() (CreateInvoiceParams, Invoice, error) {
	_, user, err := CreateRandomUserTest()
	if err != nil {
		fmt.Println("Error creating user", err)
	}

	arg := CreateInvoiceParams{
		InvoiceNumber:       utils.RandomInvoiceReceiptNumber(),
		UserInvoiceID:       int32(user.UserID),
		UserInvoiceUsername: user.Username,
		InvoicePdf:          []byte(utils.RandomString(6)),
	}

	invoiceData := []map[string]interface{}{
		{
			"userdata": "userdata stored here",
		},
		{
			"productID":       1,
			"productName":     "HIV testkit 1",
			"productQuantity": 4,
			"totalBill":       400,
		},
		{
			"productID":       2,
			"productName":     "HIV testkit 2",
			"productQuantity": 5,
			"totalBill":       250,
		},
	}

	jsonInvoiceData, err := json.Marshal(invoiceData)
	if err != nil {
		fmt.Println("Error marshaling invoice data:", err)
	}

	arg.InvoiceData = jsonInvoiceData
	result, err := testStore.CreateInvoice(context.Background(), arg)

	// var invoiceDataFromDB []map[string]interface{}

	// if err := json.Unmarshal(result.InvoiceData, &invoiceDataFromDB); err != nil {
	// 	fmt.Println("Error unmarshaling invoice data:", err)
	// }
	return arg, result, err
}

func TestCreateInvoice(t *testing.T) {
	arg, invoice, err := createInvoiceTest()
	require.NoError(t, err)
	require.NotEmpty(t, invoice)

	require.Equal(t, arg.InvoiceNumber, arg.InvoiceNumber)
	require.Equal(t, arg.UserInvoiceID, arg.UserInvoiceID)

	require.NotZero(t, invoice.InvoiceID)
	require.NotZero(t, invoice.CreatedAt)
}

func TestGetInvoice(t *testing.T) {
	_, invoice, err := createInvoiceTest()
	require.NoError(t, err)
	require.NotEmpty(t, invoice)

	returnedInvoice, err := testStore.GetInvoice(context.Background(), invoice.InvoiceNumber)
	require.NoError(t, err)
	require.NotEmpty(t, returnedInvoice)

	require.Equal(t, invoice.InvoiceNumber, returnedInvoice.InvoiceNumber)
	require.Equal(t, invoice.UserInvoiceID, returnedInvoice.UserInvoiceID)

	require.NotZero(t, returnedInvoice.InvoiceID)
	require.NotZero(t, returnedInvoice.CreatedAt)
}

func TestGetUserInvoice(t *testing.T) {
	arg, invoice, err := createInvoiceTest()
	require.NoError(t, err)
	require.NotEmpty(t, invoice)
	require.NotEmpty(t, arg)

	invoices, err := testStore.GetUserInvoicesByID(context.Background(), GetUserInvoicesByIDParams{
		Limit:         10,
		Offset:        0,
		UserInvoiceID: 8,
	})
	require.NoError(t, err)
	for _, invoice := range invoices {
		require.NotEmpty(t, invoice)
		require.NotEmpty(t, invoice.InvoiceID)
		require.NotEmpty(t, invoice.InvoiceNumber)
		require.NotEmpty(t, invoice.UserInvoiceID)
		require.NotEmpty(t, invoice.CreatedAt)
		// fmt.Printf("%v\t%v\t%v\t%v", invoice.InvoiceID, invoice.InvoiceNumber, invoice.UserInvoiceID, invoice.CreatedAt)
	}

}

func TestListInvoices(t *testing.T) {
	var invoices []Invoice
	var err error
	invoices, err = testStore.ListInvoices(context.Background(), ListInvoicesParams{
		Limit:  1,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, invoices)
	for _, invoice := range invoices {
		require.NotEmpty(t, invoice.InvoiceID)
		require.NotEmpty(t, invoice.InvoiceNumber)
		require.NotEmpty(t, invoice.UserInvoiceID)
		require.NotEmpty(t, invoice.CreatedAt)
		// fmt.Printf("%v\t%v\t%v\t%v", invoice.InvoiceID, invoice.InvoiceNumber, invoice.UserInvoiceID, invoice.CreatedAt)
	}
}
