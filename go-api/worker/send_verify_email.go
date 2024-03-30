package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const SendVerifyEmailTask = "task:send_verify_email"

type SendEmailVerifyPayload struct {
	Username string `json:"username"`
}

func (distributor RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload SendEmailVerifyPayload, opt ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(SendVerifyEmailTask, jsonPayload, opt...)
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

func (processor *RedisTaskProcessor) ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payLoad SendEmailVerifyPayload
	if err := json.Unmarshal(task.Payload(), &payLoad); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	user, err := processor.store.GetUserByUsename(ctx, payLoad.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No user found: %w", err)
		}
		return fmt.Errorf("Internal error: %w", err)
	}

	accessToken, err := processor.tokenMaker.CreateToken(user.Username, (24 * time.Hour))
	if err != nil {
		return fmt.Errorf("Internal error: %w", err)
	}

	// Use real project ulr below
	resetPasswordLink := fmt.Sprintf("%v/resetit?token=%v", "http://127.0.0.1:5000", accessToken) // URL + TOKEN for passwordreset
	emailBody := fmt.Sprintf(`
	<h1>Hello %s</h1>
	<p>A new account has been created for the Kokomed Supplies System. Please create your Password and Login to check it out</p>
	<a href="%s" style="display:inline-block; padding:10px 20px; background-color:#007BFF; color:#fff; text-decoration:none; border-radius:5px;">Create Password</a>
	<h5>The link is valid for 10 Minutes</h5>
`, user.Username, resetPasswordLink)

	err = processor.sender.SendMail("Create Password", emailBody, []string{user.Email}, nil, nil, "", nil)
	if err != nil {
		return fmt.Errorf("faild to send email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("body", task.Payload()).
		Str("email", user.Email).
		Msg("tasked processed successfull")

	return nil
}
