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

const SendResetPasswordEmailTask = "task:send_reset_password_email_task"

type SendResetPasswordEmail struct {
	Email string `json:"email"`
}

func (distributor *RedisTaskDistributor) DistributeSendResetPasswordEmail(
	ctx context.Context,
	payload SendResetPasswordEmail,
	opt ...asynq.Option,
) error {
	jsonResetPasswordPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal reset password payload: %w", err)
	}

	task := asynq.NewTask(SendResetPasswordEmailTask, jsonResetPasswordPayload, opt...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("Failed to enqueue reset password task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessSendResetPasswordEmail(ctx context.Context, task *asynq.Task) error {
	var payLoad SendResetPasswordEmail
	if err := json.Unmarshal(task.Payload(), &payLoad); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	user, err := processor.store.GetUserByEmail(ctx, payLoad.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No user found: %w", err)
		}
		return fmt.Errorf("Internal error: %w", err)
	}

	accessToken, err := processor.tokenMaker.CreateToken(user.Username, (10 * time.Minute))
	if err != nil {
		return fmt.Errorf("Internal error: %w", err)
	}

	resetPasswordLink := fmt.Sprintf("%v/resetit?token=%v", processor.config.PUBLIC_URL, accessToken)
	emailBody := fmt.Sprintf(`
	<h1>Hello %s</h1>
	<p>We received a request to reset your password. Click the link below to reset it:</p>
	<a href="%s" style="display:inline-block; padding:10px 20px; background-color:#007BFF; color:#fff; text-decoration:none; border-radius:5px;">Reset Password</a>
	<h5>The link is valid for 10 Minutes</h5>
	<a href="https://inventory-system-production-378e.up.railway.app/">https://inventory-system-production-378e.up.railway.app/</a>
`, user.Username, resetPasswordLink)

	err = processor.sender.SendMail("Reset Password", emailBody, "application/pdf", []string{user.Email}, nil, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Str("email", user.Email).
		Msg("tasked processed successfull")

	return nil
}
