package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const SendSTKTask = "task:send_stk_push"

type SendSTKPayload struct {
	User            db.User `json:"user"`
	Amount          int64   `json:"amount"`
	TransactionData []byte  `json:"transaction_data"`
}

func (distributor *RedisTaskDistributor) DistributeSendSTK(ctx context.Context, payload SendSTKPayload, opts ...asynq.Option) error {
	jsonSendPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Couldn't marshal send_stkPayload: %w", err)
	}

	task := asynq.NewTask(SendSTKTask, jsonSendPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue send_stk: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessSendSTK(ctx context.Context, task *asynq.Task) error {
	var sendSTKPayload SendSTKPayload
	if err := json.Unmarshal(task.Payload(), &sendSTKPayload); err != nil {
		return fmt.Errorf("Failed to unmarshal send_stk payload: %w", err)
	}

	rescDescription, trasactionID, err := utils.SendSTK(strconv.Itoa(int(sendSTKPayload.Amount)), sendSTKPayload.User.UserID, sendSTKPayload.User.PhoneNumber)
	if err != nil {
		return fmt.Errorf("Failed to send stk : %w", err)
	}

	// intAmount, err := strconv.ParseInt(sendSTKPayload.Amount, 0, 64)
	// if err != nil {
	// 	return fmt.Errorf("Failed to parse amount to int64: %w", err)
	// }

	transaction, err := processor.store.CreateTransaction(ctx, db.CreateTransactionParams{
		TransactionID:     trasactionID,
		TransactionUserID: int32(sendSTKPayload.User.UserID),
		Amount:            int32(sendSTKPayload.Amount),
		PhoneNumber:       sendSTKPayload.User.PhoneNumber,
		DataSold:          sendSTKPayload.TransactionData,
		ResultDescription: rescDescription,
	})
	if err != nil {
		log.Info().Msg("failed to create transaction")
		return fmt.Errorf("failed to create new transaction: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("transaction", transaction.TransactionID).
		Str("username", sendSTKPayload.User.Username).
		Msg("tasked processed successfull")

	return nil
}
