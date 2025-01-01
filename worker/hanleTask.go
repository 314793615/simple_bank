package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "github.com/314793615/simplebank/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type Processor interface {
	ProcessTask(ctx context.Context, task *asynq.Task) error
}

type TaskProcessor struct {
	store  db.Store
	server *asynq.Server
}

func NewTaskProcessor(opt asynq.RedisClientOpt, store db.Store) *TaskProcessor {
	server := asynq.NewServer(opt, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
		},
	})
	return &TaskProcessor{
		store:  store,
		server: server,
	}
}

func (t *TaskProcessor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	payload := task.Payload()
	var taskPayload SendEmailPayload
	if err := json.Unmarshal(payload, &taskPayload); err != nil {
		return fmt.Errorf("failed to unmarshall the task payload: %w", err)
	}
	username := taskPayload.Username
	user, err := t.store.GetUser(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("the user %s does not exist", username)
		}
		return fmt.Errorf("failed to get user %s: %w", username, err)
	}
	// send email
	log.Info().Str("username", user.Username).Msg("successfully process task")
	return nil
}

func (t *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(sendEmailTask, t.ProcessTask)
	return t.server.Start(mux)
}
