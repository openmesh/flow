package flow

type EventSource struct {
}

type Runner interface {
	Run(inputs map[string]interface{}) (map[string]interface{}, error)
}

// Consider data implementation

// Things to store

// Workflow
// Steps
