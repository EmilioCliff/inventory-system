package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const ProcessMpesaCallbackTask = "task:process_mpesa_callback"

type ProcessMpesaCallbackPayload struct {
	UserID        string                 `json:"user_id"`
	TransactionID string                 `json:"transaction_id"`
	Body          map[string]interface{} `json:"body"`
}

func (distributor *RedisTaskDistributor) DistributeProcessMpesaCallback(ctx context.Context, payload ProcessMpesaCallbackPayload, opts ...asynq.Option) error {
	jsonProcessMpesaPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal process mpesa payload: %w", err)
	}

	// optsNew := []asynq.Option{
	// 	asynq.MaxRetry(2),
	// }
	task := asynq.NewTask(ProcessMpesaCallbackTask, jsonProcessMpesaPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue process mpesa callback task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("body", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessMpesaCallback(ctx context.Context, task *asynq.Task) error {
	var mpesaCallbackPayload ProcessMpesaCallbackPayload
	if err := json.Unmarshal(task.Payload(), &mpesaCallbackPayload); err != nil {
		return fmt.Errorf("failed to unmarshal mpesa callback payload: %w", err)
	}

	intUserID, _ := strconv.Atoi(mpesaCallbackPayload.UserID)
	user, err := processor.store.GetUserForUpdate(ctx, int64(intUserID))
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", err)
		}
		return fmt.Errorf("internal server error: %w", err)
	}

	transaction, err := processor.store.GetTransaction(ctx, mpesaCallbackPayload.TransactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("transaction not found: %w", err)
		}
		return fmt.Errorf("internal server error: %w", err)
	}

	log.Info().Msgf("In processMpesaCallbackData: %s\nUser and Transaction: %v:%v", mpesaCallbackPayload.Body["Body"].(map[string]interface{}), user, transaction)

	bodyValue, _ := mpesaCallbackPayload.Body["Body"].(map[string]interface{})
	stkCallbackValue, _ := bodyValue["stkCallback"].(map[string]interface{})

	if len(stkCallbackValue) != 5 {
		// ("No CallbackMetadata in the response: %w", asynq.SkipRetry)
		log.Error().
			Str("type", task.Type()).
			Bytes("body", task.Payload()).
			Str("transaction_id", transaction.TransactionID).
			Str("mpesa_receipt_number", transaction.MpesaReceiptNumber).
			Str("info", "transaction not successful").
			// Str("invoice_number", invoiceGenerated.InvoiceNumber).
			Msg("tasked processed successfull")
		return nil
	}
	// var resultCode int
	// if val, ok := stkCallbackValue["ResultCode"].(int); ok {
	// 	resultCode = int(val)
	// 	if resultCode != 0 {
	// 		resultDesc, _ := stkCallbackValue["ResultDesc"].(string)
	// 		newError := errors.New(fmt.Sprintf("resultCode %v not same as 0. Description: %v", resultCode, resultDesc))
	// 		log.Error().Err(newError)
	// 		// redirectToPythonApp(user, transaction, fmt.Errorf(resultDesc))
	// 		return fmt.Errorf("different resultcode not 0: %w", newError)
	// 	}
	// 	log.Info().Msg("ResultCode is zero can continue")
	// }

	metaData, _ := stkCallbackValue["CallbackMetadata"].(map[string]interface{})
	items, _ := metaData["Item"].([]interface{})

	var phoneNumber, mpesaReceiptNumber string
	if len(items) > 0 {
		if val, ok := items[3].(map[string]interface{}); ok {
			phoneNumber, _ = val["Value"].(string)
			log.Info().Msgf("number: %s", val["Value"].(string))
		}

		if val, ok := items[1].(map[string]interface{}); ok {
			mpesaReceiptNumber, _ = val["Value"].(string)
			log.Info().Msgf("mpesa_receipt%s", val["Value"].(string))
		}
	}
	// if len(items) > 3 {
	// 	if val, ok := items[3].(map[string]interface{}); ok {
	// 		phoneNumber, _ = val["Value"].(string)
	// 	}
	// }

	// if len(items) > 1 {
	// 	if val, ok := items[1].(map[string]interface{}); ok {
	// 		mpesaReceiptNumber, _ = val["Value"].(string)
	// 	}
	// }

	_, err = processor.store.UpdateTransaction(ctx, db.UpdateTransactionParams{
		TransactionID:      transaction.TransactionID,
		MpesaReceiptNumber: mpesaReceiptNumber,
		PhoneNumber:        phoneNumber,
	})
	if err != nil {
		// redirectToPythonApp(user, transaction, err)
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	var data map[string][]int8
	if unerr := json.Unmarshal(transaction.DataSold, &data); unerr != nil {
		// redirectToPythonApp(user, transaction, err)
		return fmt.Errorf("failed to unmarshal transaction data sold: %w", unerr)
	}

	var newProducts []db.Product
	for _, id := range data["products_id"] {
		addProduct, err := processor.store.GetProduct(ctx, int64(id))
		if err != nil {
			if err == sql.ErrNoRows {
				// redirectToPythonApp(user, transaction, err)
				return fmt.Errorf("product not found: %w", err)
			}
			// redirectToPythonApp(user, transaction, err)
			return fmt.Errorf("error getting product: %w", err)
		}

		newProducts = append(newProducts, addProduct)
	}
	_, err = processor.store.ReduceClientStockTx(ctx, db.ReduceClientStockParams{
		Client:         user,
		ProducToReduce: newProducts,
		Amount:         data["quantities"],
		Transaction:    transaction,
		AfterPaying: func() error {
			receiptTaskPayload := &GenerateReceiptAndSendEmailPayload{
				User:     user,
				Products: newProducts,
				Amount:   data["quantities"],
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
				asynq.Queue(QueueDefault),
			}

			return processor.distributor.DistributeGenerateAndSendReceipt(ctx, *receiptTaskPayload, opts...)
		},
	})
	if err != nil {
		// redirectToPythonApp(user, transaction, err)
		return fmt.Errorf("error reducing client stock: %w", err)
	}

	processor.store.ChangeStatus(ctx, db.ChangeStatusParams{
		TransactionID: transaction.TransactionID,
		Status:        true,
	})

	log.Info().
		Str("type", task.Type()).
		Bytes("body", task.Payload()).
		Str("transaction_id", transaction.TransactionID).
		Str("mpesa_receipt_number", transaction.MpesaReceiptNumber).
		Str("info", "transaction successful").
		// Str("invoice_number", invoiceGenerated.InvoiceNumber).
		Msg("tasked processed successfull")

	return nil
}