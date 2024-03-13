package db

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/stretchr/testify/require"
)

func createReceiptTest() (CreateReceiptParams, Receipt, error) {
	_, user, err := CreateRandomUserTest()
	if err != nil {
		fmt.Println("Error creating user", err)
	}

	randomOffset := time.Duration(rand.Intn(10)) * time.Second
	uniqueTime := time.Now().Add(randomOffset)
	transactionParam := CreateTransactionParams{
		TransactionID: uniqueTime.Format("20060102150405"),
		Amount:        1000,
	}

	args := CreateReceiptParams{
		ReceiptNumber:       transactionParam.TransactionID,
		UserReceiptID:       int32(user.UserID),
		UserReceiptUsername: user.Username,
		ReceiptPdf:          []byte(utils.RandomString(6)),
	}

	receiptData := []map[string]interface{}{
		{
			"userdata": "userdata stored here",
		},
		{
			"productID":       2,
			"productName":     "HIV testkit 2",
			"productQuantity": 9,
			"totalBill":       600,
		},
		{
			"productID":       3,
			"productName":     "HIV testkit 3",
			"productQuantity": 5,
			"totalBill":       250,
		},
	}
	jsonReceiptData, err := json.Marshal(receiptData)
	if err != nil {
		fmt.Println("error marshaling recceipt data: ", err)
	}
	args.ReceiptData = jsonReceiptData
	transactionParam.DataSold = jsonReceiptData
	_, _ = testStore.CreateTransaction(context.Background(), transactionParam)
	newReceipt, err := testStore.CreateReceipt(context.Background(), args)
	if err != nil {
		fmt.Println("error creating recceipt: ", err)
	}
	return args, newReceipt, err
}

func TestCreateReceipt(t *testing.T) {
	args, receipt, err := createReceiptTest()
	require.NoError(t, err)
	require.NotEmpty(t, args)
	require.NotEmpty(t, receipt)

	require.Equal(t, args.ReceiptNumber, receipt.ReceiptNumber)
	require.Equal(t, args.UserReceiptID, receipt.UserReceiptID)
}

func TestGetReceipt(t *testing.T) {
	args, receipt, err := createReceiptTest()
	require.NoError(t, err)
	require.NotEmpty(t, args)
	require.NotEmpty(t, receipt)

	receiptGet, err := testStore.GetReceipt(context.Background(), receipt.ReceiptNumber)
	require.NoError(t, err)
	require.NotEmpty(t, receiptGet)

	require.Equal(t, receipt.ReceiptNumber, receiptGet.ReceiptNumber)
	require.Equal(t, receipt.UserReceiptID, receiptGet.UserReceiptID)
}

func TestGetUSerReceipts(t *testing.T) {
	// _, user, err := CreateRandomUserTest()
	// if err != nil {
	// 	fmt.Println("Error creating user", err)
	// }
	// require.NotEmpty(t, user)

	var receipts []Receipt
	receipts, err := testStore.GetUserReceiptsByID(context.Background(), int32(18))
	require.NoError(t, err)
	// require.NotEmpty(t, receipts)

	for _, receipt := range receipts {
		require.NotEmpty(t, receipt.ReceiptID)
		require.NotEmpty(t, receipt.ReceiptNumber)
		require.NotEmpty(t, receipt.UserReceiptID)
		require.NotEmpty(t, receipt.CreatedAt)
	}
}

func TestListReceipts(t *testing.T) {
	receipts, err := testStore.ListReceipts(context.Background(), ListReceiptsParams{
		Limit:  1,
		Offset: 0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, receipts)

	for _, receipt := range receipts {
		require.NotEmpty(t, receipt.ReceiptID)
		require.NotEmpty(t, receipt.ReceiptNumber)
		require.NotEmpty(t, receipt.UserReceiptID)
		require.NotEmpty(t, receipt.CreatedAt)
		// fmt.Printf("%v\t%v\t%v\t%v", receipt.ReceiptID, receipt.ReceiptNumber, receipt.UserReceiptID, receipt.CreatedAt)
	}
}
