package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const SendRequestStockTask = "task:send_request_stock"

type RequestStockPayload struct {
	UsernameID int32   `json:"user_id"`
	Products   []int32 `json:"products"`
	Quantities []int32 `json:"quantities"`
}

func (distributor *RedisTaskDistributor) DistributeSendRequestToAdmin(ctx context.Context, payload RequestStockPayload, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error marshaling payload: %w", err)
	}

	task := asynq.NewTask(SendRequestStockTask, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("body", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessSendRequestStock(ctx context.Context, task *asynq.Task) error {
	var payload RequestStockPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("Failed to unmarshal payload: %w", err)
	}

	var dataToSend []map[string]interface{}
	for idx, productID := range payload.Products {
		product, err := processor.store.GetProduct(ctx, int64(productID))
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("No product found: %w", err)
			}
			return fmt.Errorf("Internal error: %w", err)
		}

		dataToSend = append(dataToSend, map[string]interface{}{
			"ProductName": product.ProductName,
			"Quantity":    payload.Quantities[idx],
		})
	}

	user, err := processor.store.GetUser(ctx, int64(payload.UsernameID))
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No user found: %w", err)
		}
		return fmt.Errorf("Internal error: %w", err)
	}

	emailBody := fmt.Sprintf(`
	<h1>Hello Joan</h1>
	<p>You have received a new Stock Request from %s, %v - %v. Find the more details below.</p>
	%s
`, user.Username, user.PhoneNumber, user.Email, generateHTMLTable(dataToSend))

	err = processor.sender.SendMail("Request Stock", emailBody, "application/json", []string{"jcherono8@gmail.com"}, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("type", "SendRequestStock").
		Str("email", user.Email).
		Msg("tasked processed successfully")

	return nil
}

func generateHTMLTable(products []map[string]interface{}) string {
	const tableTemplate = `
	<table border="1" style="border-collapse: collapse; width: 70%; text-align: center;">
		<tr style="background-color: #333; color: white; font-weight: bold; text-transform: uppercase;">
			<td style="padding: 10px;">Product</td>
			<td style="padding: 10px;">Quantity</td>
		</tr>
		{{range .}}
		<tr>
			<td style="padding: 10px;">{{.ProductName}}</td>
			<td style="padding: 10px;">{{.Quantity}}</td>
		</tr>
		{{end}}
	</table>
	`
	tmpl, err := template.New("table").Parse(tableTemplate)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse table template")
		return ""
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, products)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute table template")
		return ""
	}

	return buf.String()
}
