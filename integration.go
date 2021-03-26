package flow

import "context"

//type Integration struct {
//	Name      string    `json:"name"`
//	Reference string    `json:"reference"`
//	Triggers  []Trigger `json:"triggers"`
//	Actions   []Action  `json:"actions"`
//}

// Integration as interface?
type Integration interface {
	Describer
	// It would be nice if these methods were encapsulated somewhere else
	GetActions() []Action
	GetTriggers() []Trigger
}

type Describer interface {
	Describe() Metadata
}

type Metadata struct {
	Name        string `json:"name"`
	Reference   string `json:"reference"`
	Description string `json:"description"`
}

type Installation struct {
	IntegrationReference string `json:"integration_reference`
	// TODO this probably needs config and stuff
}

type Trigger interface {
	Describer
	Emitter
	// Fields []Field `json:"fields"`
}

type Action interface {
	Describer
	Runner
	Handler
	Emitter
	//InputFields  []Field `json:"input_fields"`
	//OutputFields []Field `json:"output_fields"`
}

type Handler interface {
	Handle(inputs []Field) error
	DefineInputs() []Field
}

type Emitter interface {
	Emit(outputs []Field) error
	DefineOutputs() []Field
}

type FieldType string

const (
	Number   FieldType = "number"
	String             = "string"
	Boolean            = "boolean"
	DateTime           = "datetime"
	Complex            = "complex"
)

type Field struct {
	Key      string    `json:"key"`
	Label    string    `json:"label"`
	Required bool      `json:"required"`
	Type     FieldType `json:"type"`
	HelpText string    `json:"help_text"`
	Default  string    `json:"default"`
}

type IntegrationService interface {
	InstallIntegration(ctx context.Context, req InstallIntegrationRequest) (*Installation, error)
	UninstallIntegration(ctx context.Context, req UninstallIntegrationRequest) error
	GetIntegrations(ctx context.Context, req GetIntegrationsRequest) ([]*Integration, int, error)
}

type InstallIntegrationRequest struct {
	IntegrationReference string `json:"integration_reference"`
}

type UninstallIntegrationRequest struct {
	IntegrationReference string `json:"integration_reference"`
}

type GetIntegrationsRequest struct {
}
