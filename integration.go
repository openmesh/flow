package flow

import "context"

type Integration struct {
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Key         string    `json:"key"`
	BaseURL     string    `json:"base_url"`
	Triggers    []Trigger `json:"triggers"`
	Actions     []Action  `json:"actions"`
}

type Trigger struct {
	Key         string        `json:"key"`
	Label       string        `json:"label"`
	Description string        `json:"description"`
	Endpoint    string        `json:"endpoint"`
	Method      string        `json:"method"`
	Outputs     []OutputField `json:"outputs"`
}

type Action struct {
	Key         string        `json:"key"`
	Label       string        `json:"label"`
	Description string        `json:"description"`
	Endpoint    string        `json:"endpoint"`
	Method      string        `json:"method"`
	Inputs      []InputField  `json:"inputs"`
	Outputs     []OutputField `json:"outputs"`
}

type InputField struct {
	Key         string    `json:"key"`
	Label       string    `json:"label"`
	Description string    `json:"description"`
	Required    bool      `json:"required"`
	Type        FieldType `json:"type"`
	Default     string    `json:"default"`
	Example     string    `json:"example"`
}

type OutputField struct {
	Label       string    `json:"label"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Type        FieldType `json:"type"`
	Path        string    `json:"path"`
}

//// Integration as interface?
//type Integration interface {
//	Describer
//	// It would be nice if these methods were encapsulated somewhere else
//	GetActions() []Action
//	GetTriggers() []Trigger
//}

type Describer interface {
	Describe() Metadata
}

type Metadata struct {
	Name        string `json:"name"`
	Ref         string `json:"ref"`
	Description string `json:"description"`
}

type Installation struct {
	IntegrationReference string `json:"integration_reference`
	// TODO this probably needs config and stuff
}

//type Trigger interface {
//	Describer
//	Emitter
//	// Fields []Field `json:"fields"`
//}
//
//type Action interface {
//	Describer
//	Runner
//	Handler
//	Emitter
//	//InputFields  []Field `json:"input_fields"`
//	//OutputFields []Field `json:"output_fields"`
//}

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
	FieldTypeNumber   FieldType = "number"
	FieldTypeString             = "string"
	FieldTypeBoolean            = "boolean"
	FieldTypeDateTime           = "datetime"
	FieldTypeComplex            = "complex"
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
	GetIntegrations(ctx context.Context, req GetIntegrationsRequest) ([]*Integration, int, error)
}

type GetIntegrationsRequest struct {
}
