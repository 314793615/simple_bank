package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	sendEmailTask = "send_email_task"
)

type Distributor interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload SendEmailPayload, opts ...asynq.Option) error
}

type RedisDistributor struct {
	Client *asynq.Client
}

type SendEmailPayload struct {
	Username string `json:"username"`
}

func NewRedisTaskDistributor(opts asynq.RedisClientOpt) Distributor {
	return &RedisDistributor{
		Client: asynq.NewClient(opts),
	}
}

func (r *RedisDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload SendEmailPayload, opts ...asynq.Option) error {
	taskPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	task := asynq.NewTask(sendEmailTask, taskPayload, opts...)

	info, err := r.Client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("send verify email task")
	return nil
}
