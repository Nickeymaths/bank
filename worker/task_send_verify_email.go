package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	db "github.com/Nickeymaths/bank/db/sqlc"
	"github.com/Nickeymaths/bank/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type PayloadTaskSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskVeriyEmail(
	ctx context.Context,
	payload *PayloadTaskSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshaled payload %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueued task %w", err)
	}

	log.Info().Str("type", task.Type()).
		Bytes("payload", jsonPayload).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadTaskSendVerifyEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal task payload %w", err)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user doesn't exist: %w", err)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	arg := db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create verify email entry: %w", err)
	}

	subject := "Welcome to our bank service"
	verifyUrl := fmt.Sprintf("http://localhost:4000/v1/verify_email?id=%d&secret_code=%s", verifyEmail.ID, arg.SecretCode)
	content := fmt.Sprintf(`<h1> Hello, %s</h1><br/><p>Please click <a href="%s">here</a> to verify your email</p></br>`, user.Username, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.Send(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", task.Type()).
		Str("email", user.Email).Msg("processed task")
	return nil
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
