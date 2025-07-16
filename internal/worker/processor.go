package worker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"
)

type Processor struct {
	srv *asynq.Server
}

func NewProcessor(redisAddr string) *Processor {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	return &Processor{srv: srv}
}

func (p *Processor) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeGenerateThumbnail, p.handleGenerateThumbnail)
	mux.HandleFunc(TypeSendWelcomeEmail, p.handleSendWelcomeEmail)
	return p.srv.Run(mux)
}

func (p *Processor) handleGenerateThumbnail(ctx context.Context, t *asynq.Task) error {
	var payload GenerateThumbnailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	// TODO: fetch image, resize & upload thumb, persist URL in DB
	log.Printf("Generating thumbnail for ad %d from %s", payload.AdID, payload.ImageURL)
	return nil
}

func (p *Processor) handleSendWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	// …send email…
	return nil
}
