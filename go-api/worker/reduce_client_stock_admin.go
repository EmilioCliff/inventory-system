package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const ReduceClientStockAdminTask = "task:reduce_client_stock_admin"

type ReduceClientStockByAdminPayload struct {
	Amount             int64        `json:"amount"`
	ProducToReduce     []db.Product `json:"productstoadd"`
	Quantities         []int64      `json:"quantities"`
	PhoneNumber        string       `json:"phone_number"`
	MpesaReceiptNumber string       `json:"mpesa_receipt_number"`
	Description        string       `json:"description"`
	UserID             int32        `json:"user_id"`
	TransactionData    []byte       `json:"transaction_data"`
}

func (distributor *RedisTaskDistributor) DistributeSendReduceClientStockAdmin(ctx context.Context, payload ReduceClientStockByAdminPayload, opts ...asynq.Option) error {
	jsonProcessMpesaPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal reduce stock by admin payload: %w", err)
	}

	task := asynq.NewTask(ReduceClientStockAdminTask, jsonProcessMpesaPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue reduce stock by admin task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessReduceClientStockByAdmin(ctx context.Context, task *asynq.Task) error {
	var reduceStockByAdminPayload ReduceClientStockByAdminPayload
	if err := json.Unmarshal(task.Payload(), &reduceStockByAdminPayload); err != nil {
		return fmt.Errorf("failed to unmarshal reduce stock by admin payload: %w", err)
	}

	user, err := processor.store.GetUser(ctx, int64(reduceStockByAdminPayload.UserID))
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	transaction, err := processor.store.ReduceClientStockByAdminTx(ctx, db.ReduceClientStockByAdmin{
		Amount:             reduceStockByAdminPayload.Amount,
		PhoneNumber:        reduceStockByAdminPayload.PhoneNumber,
		MpesaReceiptNumber: reduceStockByAdminPayload.MpesaReceiptNumber,
		Description:        reduceStockByAdminPayload.Description,
		UserID:             reduceStockByAdminPayload.UserID,
		TransactionData:    reduceStockByAdminPayload.TransactionData,
	})
	if err != nil {
		return fmt.Errorf("failed to reduce stock by admin transaction: %w", err)
	}

	err = processor.store.ReduceClientStockTx(ctx, db.ReduceClientStockParams{
		ClientID:       user.UserID,
		ProducToReduce: reduceStockByAdminPayload.ProducToReduce,
		Amount:         reduceStockByAdminPayload.Quantities,
		Transaction:    transaction,
		AfterPaying: func(data []map[string]interface{}) error {
			receiptTaskPayload := &GenerateReceiptAndSendEmailPayload{
				User:        user,
				Transaction: transaction,
				ReceiptData: data,
				ByAdmin:     true,
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
		return fmt.Errorf("error reducing client stock: %w", err)
	}

	err = processor.store.ChangeStatus(ctx, db.ChangeStatusParams{
		TransactionID: transaction.TransactionID,
		Status:        true,
	})
	if err != nil {
		return fmt.Errorf("error changing transaction data: %w", err)
	}

	err = processor.store.ChangePaymentMethod(ctx, db.ChangePaymentMethodParams{
		TransactionID: transaction.TransactionID,
		PaymentMethod: "By Admin",
	})
	if err != nil {
		return fmt.Errorf("error changing payment method data: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("receipt_generated_no", transaction.TransactionID).
		Str("mpesa_receipt_number", transaction.MpesaReceiptNumber).
		Str("info", "transaction successful").
		Msg("tasked processed successfull")

	return nil
}
