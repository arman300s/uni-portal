package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/arman300s/uni-portal/pkg/tasks"
	"github.com/hibiken/asynq"
)

func main() {
	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSendWelcomeEmail, func(ctx context.Context, t *asynq.Task) error {
		var p tasks.SendWelcomeEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return err
		}
		return tasks.ExecuteSendWelcomeEmail(ctx, p)
	})

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run worker: %v", err)
	}
}
