package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload SendEmailVerifyPayload,
		opt ...asynq.Option,
	) error
	DistributeGenerateAndSendInvoice(
		ctx context.Context,
		payload GenerateInvoiceAndSendEmailPayload,
		opt ...asynq.Option,
	) error
	DistributeGenerateAndSendReceipt(
		ctx context.Context,
		payload GenerateReceiptAndSendEmailPayload,
		opt ...asynq.Option,
	) error
	DistributeSendResetPasswordEmail(
		ctx context.Context,
		payload SendResetPasswordEmail,
		opt ...asynq.Option,
	) error
	DistributeSendSTK(
		ctx context.Context,
		payload SendSTKPayload,
		opts ...asynq.Option,
	) error
	DistributeProcessMpesaCallback(
		ctx context.Context,
		payload ProcessMpesaCallbackPayload,
		opts ...asynq.Option,
	) error
	DistributeTakeAndSendDBsnapshots(
		ctx context.Context,
		word string,
		opts ...asynq.Option,
	) error
	DistributeSendRequestToAdmin(
		ctx context.Context,
		payload RequestStockPayload,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
