package worker

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypeGenerateThumbnail = "ad:generate_thumbnail"
	TypeSendWelcomeEmail  = "user:send_welcome_email"
)

type GenerateThumbnailPayload struct {
	AdID     int64  `json:"ad_id"`
	ImageURL string `json:"image_url"`
}

// NewGenerateThumbnailTask enqueues an image‐processing job.
func NewGenerateThumbnailTask(adID int64, imageURL string) (*asynq.Task, error) {
	p, err := json.Marshal(GenerateThumbnailPayload{AdID: adID, ImageURL: imageURL})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeGenerateThumbnail, p, asynq.MaxRetry(3), asynq.Timeout(30*time.Second)), nil
}

// similarly for welcome email…
