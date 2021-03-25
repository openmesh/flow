package github

import (
	"github.com/google/go-github/v33/github"
	"github.com/openmesh/flow"
)

type PushTrigger struct {
	Payload github.PushEvent `json:"payload"`
}

func NewPushTrigger(payload github.PushEvent) flow.Trigger {
	return &PushTrigger{
		Payload: payload,
	}
}

func (p *PushTrigger) Describe() flow.Metadata {
	return flow.Metadata{
		Name:        "Push",
		Reference:   "GITHUB_PUSH",
		Description: "Trigger that occurs when one or more commits are pushed to a repository branch or tag.",
	}
}

func (p *PushTrigger) Emit(outputs []flow.Field) error {
	// map push payload to outputs
	return nil
}

func (p *PushTrigger) DefineOutputs() []flow.Field {
	panic("implement me")
}
