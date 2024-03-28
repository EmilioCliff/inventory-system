package worker

import "github.com/hibiken/asynq"

type TaskDistributor interface {
}

type RedisTaskDistributor struct {
	client *asynq.Client
}
