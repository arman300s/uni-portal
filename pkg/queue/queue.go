package queue

import (
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

var Client *asynq.Client

func Init(redisAddr string) {
	Client = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
}

func Enqueue(taskType string, payload interface{}, delay time.Duration) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(taskType, data)
	_, err = Client.Enqueue(task, asynq.ProcessIn(delay))
	if err != nil {
		log.Println("failed to enqueue task:", err)
	}
	return err
}
