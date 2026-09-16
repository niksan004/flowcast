package steps

import (
	"gopkg.in/yaml.v3"
)

type Workflow struct {
	Steps []RawStep `yaml:"steps"`
}

type RawStep struct {
	Name string         `yaml:"name"`
	Data map[string]any `yaml:",inline"`
}

func mapToStruct(data map[string]any, out any) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, out)
}

func UnmarshalWorkflow(file []byte) (*Workflow, error) {
	// turn file into go structs
	var wf Workflow
	if err := yaml.Unmarshal(file, &wf); err != nil {
		return nil, err
	}

	return &wf, nil
}
