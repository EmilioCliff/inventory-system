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

const GenerateInvoiceAndSendEmailTask = "task:generate_invoice_and_send_email"

type GenerateInvoiceAndSendEmailPayload struct {
	User     db.User      `json:"username"`
	Products []db.Product `json:"products"`
	Amount   []int64      `json:"amount"`
}

func (distributor *RedisTaskDistributor) DistributeGenerateAndSendInvoice(ctx context.Context, payload GenerateInvoiceAndSendEmailPayload, opt ...asynq.Option) error {
	jsonGenearatePayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(GenerateInvoiceAndSendEmailTask, jsonGenearatePayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("body", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessGenerateAndSendInvoice(ctx context.Context, task *asynq.Task) error {
	var invoiceDataPayload GenerateInvoiceAndSendEmailPayload
	if err := json.Unmarshal(task.Payload(), &invoiceDataPayload); err != nil {
		return fmt.Errorf("Failed to unmarshal payload: %w", err)
	}

	invoiceData := []map[string]interface{}{
		{
			"user_contact": invoiceDataPayload.User.PhoneNumber,
			"user_address": invoiceDataPayload.User.Address,
			"user_email":   invoiceDataPayload.User.Email,
		},
	}

	for index, addProduct := range invoiceDataPayload.Products {
		invoiceData = append(invoiceData, map[string]interface{}{
			"productID":       float64(addProduct.ProductID),
			"productName":     addProduct.ProductName,
			"productQuantity": invoiceDataPayload.Amount[index],
			"totalBill":       int32(invoiceDataPayload.Amount[index]) * addProduct.UnitPrice,
		})
	}

	timestamp := time.Now().Format("20060102150405")

	invoiceC := map[string]string{
		"invoice_number":   timestamp,
		"created_at":       time.Now().Format("2006-01-02"),
		"invoice_username": invoiceDataPayload.User.Username,
	}

	pdfBytes, err := utils.SetInvoiceVariables(invoiceC, invoiceData)
	if err != nil {
		return fmt.Errorf("Error creating invoice: %w", err)
	}

	jsonInvoiceData, err := json.Marshal(invoiceData)
	if err != nil {
		return fmt.Errorf("Failed to marshal invoice data: %w", err)
	}

	invoiceGenerated, err := processor.store.CreateInvoice(ctx, db.CreateInvoiceParams{
		InvoiceNumber:       timestamp,
		UserInvoiceID:       int32(invoiceDataPayload.User.UserID),
		InvoiceData:         jsonInvoiceData,
		UserInvoiceUsername: invoiceDataPayload.User.Username,
		InvoicePdf:          pdfBytes,
	})

	emailBody := fmt.Sprintf(`
	<h1>Hello %s</h1>
	<p>We've issued products. Find the invoice attached below</p>
	<h5>Thank You For Choosing Us.</h5>
	<a href="https://inventory-system-production-378e.up.railway.app/">https://inventory-system-production-378e.up.railway.app/</a>
	`, invoiceDataPayload.User.Username)

	err = processor.sender.SendMail("Invoice Issued", emailBody, "application/pdf", []string{invoiceDataPayload.User.Email}, nil, nil, []string{"Invoice.pdf"}, [][]byte{invoiceGenerated.InvoicePdf})
	if err != nil {
		return fmt.Errorf("Failed to send email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("body", task.Payload()).
		Str("email", invoiceDataPayload.User.Email).
		Str("invoice_number", invoiceGenerated.InvoiceNumber).
		Msg("tasked processed successfull")

	return nil
}
