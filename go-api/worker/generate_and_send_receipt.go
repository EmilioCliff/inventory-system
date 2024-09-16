package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const GenerateReceiptAndSendEmailTask = "task:generate_receipt_and_send_email"

type GenerateReceiptAndSendEmailPayload struct {
	User db.User `json:"username"`
	// Products    []db.Product   `json:"products"`
	// Amount      []int8         `json:"amount"`
	Transaction db.Transaction           `json:"transaction_id"`
	ReceiptData []map[string]interface{} `json:"receipt_data"`
	ByAdmin     bool                     `json:"by_admin,omitempty"`
}

func (distributor *RedisTaskDistributor) DistributeGenerateAndSendReceipt(ctx context.Context, payload GenerateReceiptAndSendEmailPayload, opt ...asynq.Option) error {
	jsonGenearatePayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(GenerateReceiptAndSendEmailTask, jsonGenearatePayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessGenerateAndSendReceipt(ctx context.Context, task *asynq.Task) error {
	var receiptDataPayload GenerateReceiptAndSendEmailPayload
	if err := json.Unmarshal(task.Payload(), &receiptDataPayload); err != nil {
		return fmt.Errorf("Failed to unmarshal payload: %w", err)
	}

	timestamp := receiptDataPayload.Transaction.TransactionID

	receiptC := map[string]string{
		"user_contact":   receiptDataPayload.User.PhoneNumber,
		"user_address":   receiptDataPayload.User.Address,
		"user_email":     receiptDataPayload.User.Email,
		"user_username":  receiptDataPayload.User.Username,
		"receipt_number": timestamp,
		"created_at":     time.Now().Format("2006-01-02"),
	}

	if receiptDataPayload.ByAdmin {
		receiptC["by_admin"] = "true"
	}

	pdfBytes, err := utils.SetReceiptVariables(receiptC, receiptDataPayload.ReceiptData)
	if err != nil {
		return fmt.Errorf("Error creating receipt: %w", err)
	}

	jsonreceiptData, err := json.Marshal(receiptDataPayload.ReceiptData)
	if err != nil {
		return fmt.Errorf("Failed to marshal receipt data: %w", err)
	}

	receiptGenerated, err := processor.store.CreateReceipt(ctx, db.CreateReceiptParams{
		ReceiptNumber:       timestamp,
		UserReceiptID:       int32(receiptDataPayload.User.UserID),
		UserReceiptUsername: receiptDataPayload.User.Username,
		ReceiptPdf:          pdfBytes,
		ReceiptData:         jsonreceiptData,
	})

	if receiptDataPayload.ByAdmin {
		err = processor.store.ChangePaymentMethodReceipt(ctx, db.ChangePaymentMethodReceiptParams{
			ReceiptID:     receiptGenerated.ReceiptID,
			PaymentMethod: "By Admin",
		})
		if err != nil {
			return fmt.Errorf("Failed to change receipt payment method: %w", err)
		}
	}

	emailBody := fmt.Sprintf(`
	<h1>Hello %s</h1>
	<p>We've received your payment. Find the receipt attached below</p>
	<h5>Thank You For Choosing Us.</h5>
	<a href="https://inventory-system-production-378e.up.railway.app/">https://inventory-system-production-378e.up.railway.app/</a>
	`, receiptDataPayload.User.Username)

	err = processor.sender.SendMail("Receipt Issued", emailBody, "application/pdf", []string{receiptDataPayload.User.Email}, nil, nil, []string{"Receipt.pdf"}, [][]byte{receiptGenerated.ReceiptPdf})
	if err != nil {
		return fmt.Errorf("Failed to send email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("email", receiptDataPayload.User.Email).
		Str("receipt_number", receiptGenerated.ReceiptNumber).
		Msg("tasked processed successfull")

	return nil
}
