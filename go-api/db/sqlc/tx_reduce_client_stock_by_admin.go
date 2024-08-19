package db

import (
	"context"
	"fmt"
	"time"
)

type ReduceClientStockByAdmin struct {
	Amount int64 `json:"amount"`
	// ProducToReduce     []Product `json:"productstoadd"`
	// Quantities         []int64   `json:"quantities"`
	PhoneNumber        string `json:"phone_number"`
	MpesaReceiptNumber string `json:"mpesa_receipt_number"`
	Description        string `json:"description"`
	UserID             int32  `json:"user_id"`
	TransactionData    []byte `json:"transaction_data"`
}

func (store *Store) ReduceClientStockByAdminTx(ctx context.Context, arg ReduceClientStockByAdmin) (Transaction, error) {
	var result Transaction

	err := store.execTx(ctx, func(q *Queries) error {
		timestamp := time.Now().Format("20060102150405")
		transaction, err := q.CreateTransaction(ctx, CreateTransactionParams{
			TransactionID:     timestamp,
			TransactionUserID: arg.UserID,
			DataSold:          arg.TransactionData,
			PhoneNumber:       arg.PhoneNumber,
			ResultDescription: arg.Description,
			Amount:            int32(arg.Amount),
		})
		if err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		updatedTransaction, err := q.UpdateTransaction(ctx, UpdateTransactionParams{
			TransactionID:      transaction.TransactionID,
			MpesaReceiptNumber: arg.MpesaReceiptNumber,
			PhoneNumber:        transaction.PhoneNumber,
			ResultDescription:  transaction.ResultDescription,
		})
		if err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		result = updatedTransaction

		return nil
	})

	return result, err
}
