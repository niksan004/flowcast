package engine

import (
	"gopkg.in/yaml.v3"
	"project/internal/steps"
)

type Workflow struct {
	Steps []steps.RawStep `yaml:"steps"`
}

func unmarshalWorkflow(file []byte) (*Workflow, error) {
	// turn file into go structs
	var wf Workflow
	if err := yaml.Unmarshal(file, &wf); err != nil {
		return nil, err
	}

	return &wf, nil
}
