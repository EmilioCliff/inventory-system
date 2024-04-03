package worker

import (
	"context"
	"time"

	db "github.com/EmilioCliff/inventory-system/db/sqlc"
	"github.com/EmilioCliff/inventory-system/db/utils"
	"github.com/EmilioCliff/inventory-system/token"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	Start() error
	ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessGenerateAndSendInvoice(ctx context.Context, task *asynq.Task) error
	ProcessGenerateAndSendReceipt(ctx context.Context, task *asynq.Task) error
	ProcessSendResetPasswordEmail(ctx context.Context, task *asynq.Task) error
	ProcessSendSTK(ctx context.Context, task *asynq.Task) error
	ProcessMpesaCallback(ctx context.Context, task *asynq.Task) error
	ProcessTakeAndSendDBsnapshots(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server      *asynq.Server
	store       db.Store
	sender      utils.GmailSender
	tokenMaker  token.Maker
	config      utils.Config
	distributor TaskDistributor
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, sender utils.GmailSender, config utils.Config, distributor TaskDistributor) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
			QueueLow:      2,
		},
		RetryDelayFunc: CustomRetryDelayFunc,
	})

	tokenMaker, _ := token.NewPaseto(config.TOKEN_SYMMETRY_KEY)

	return &RedisTaskProcessor{
		store:       store,
		server:      server,
		sender:      sender,
		tokenMaker:  tokenMaker,
		config:      config,
		distributor: distributor,
	}
}

func CustomRetryDelayFunc(_ int, _ error, _ *asynq.Task) time.Duration {
	return 2 * time.Second
}
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(SendVerifyEmailTask, processor.ProcessSendVerifyEmail)
	mux.HandleFunc(SendResetPasswordEmailTask, processor.ProcessSendResetPasswordEmail)
	mux.HandleFunc(GenerateInvoiceAndSendEmailTask, processor.ProcessGenerateAndSendInvoice)
	mux.HandleFunc(GenerateReceiptAndSendEmailTask, processor.ProcessGenerateAndSendReceipt)
	mux.HandleFunc(SendSTKTask, processor.ProcessSendSTK)
	mux.HandleFunc(ProcessMpesaCallbackTask, processor.ProcessMpesaCallback)
	mux.HandleFunc(TakeAndSendDBsnpashotsTask, processor.ProcessTakeAndSendDBsnapshots)

	return processor.server.Start(mux)
}
