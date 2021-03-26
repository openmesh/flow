package flow

import "context"

type Webhook struct {
	Channel string      `json:"channel"`
	Payload interface{} `json:"payload"`
}

type WebhookService interface {
	IngestWebhook(ctx context.Context, req IngestWebhookRequest) error
}

type IngestWebhookRequest struct {
	// TODO maybe think of a better name for this
	Source  string      `json:"source"`
	Payload interface{} `json:"payload"`
}
