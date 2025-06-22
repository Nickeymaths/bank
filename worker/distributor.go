package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskVeriyEmail(
		ctx context.Context,
		payload *PayloadTaskSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpts asynq.RedisClientOpt) TaskDistributor {
	redisClient := asynq.NewClient(redisOpts)
	return &RedisTaskDistributor{
		client: redisClient,
	}
}
