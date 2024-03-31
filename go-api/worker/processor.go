package worker

import (
	"context"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/token"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessGenerateAndSendInvoice(ctx context.Context, task *asynq.Task) error
	ProcessGenerateAndSendReceipt(ctx context.Context, task *asynq.Task) error
	ProcessSendResetPasswordEmail(ctx context.Context, task *asynq.Task) error
	ProcessSendSTK(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server     *asynq.Server
	store      db.Store
	sender     utils.GmailSender
	tokenMaker token.Maker
	config     utils.Config
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, sender utils.GmailSender, config utils.Config) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
	})

	tokenMaker, _ := token.NewPaseto(config.TOKEN_SYMMETRY_KEY)

	return &RedisTaskProcessor{
		store:      store,
		server:     server,
		sender:     sender,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(SendVerifyEmailTask, processor.ProcessSendVerifyEmail)
	mux.HandleFunc(SendResetPasswordEmailTask, processor.ProcessSendResetPasswordEmail)
	mux.HandleFunc(GenerateInvoiceAndSendEmailTask, processor.ProcessGenerateAndSendInvoice)
	mux.HandleFunc(GenerateReceiptAndSendEmailTask, processor.ProcessGenerateAndSendReceipt)
	mux.HandleFunc(SendSTKTask, processor.ProcessSendSTK)

	return processor.server.Start(mux)
}
