package github

import (
	"github.com/google/go-github/v33/github"
	"github.com/openmesh/flow"
)

type integration struct {
	client *github.Client
}

func NewGitHubIntegration(client *github.Client) flow.Integration {
	return integration{
		client: client,
	}
}

func (i integration) Describe() flow.Metadata {
	return flow.Metadata{
		Name:        "GitHub",
		Reference:   "GITHUB",
		Description: "OpenMesh integration with the GitHub API",
	}
}

func (i integration) GetActions() []flow.Action {
	panic("implement me")
}

func (i integration) GetTriggers() []flow.Trigger {
	panic("implement me")
}
