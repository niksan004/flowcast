package steps

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

type RawStep struct {
	Name   string            `yaml:"name"`
	SaveAs map[string]string `yaml:"save_as"`
	Data   map[string]any    `yaml:",inline"`
}

type StepExecutor interface {
	Execute(sesh *ssh.Session) (any, error)
}

type StepFactory func() StepExecutor

// registry of constructors for all the step types
var Registry = map[string]StepFactory{
	"echo":    func() StepExecutor { return &EchoStep{} },
	"set_var": func() StepExecutor { return &SetVarStep{} },
	"http":    func() StepExecutor { return &HTTPStep{} },
	"shell":   func() StepExecutor { return &ShellStep{} },
	"if":      func() StepExecutor { return &IfStep{} },
	"for":     func() StepExecutor { return &ForStep{} },
}

func Decode(name string, data map[string]any) (StepExecutor, error) {
	stepConstr, exists := Registry[name]
	if !exists {
		return nil, fmt.Errorf("unknown step %s", name)
	}

	step := stepConstr()
	if err := mapToStruct(data, step); err != nil {
		return nil, err
	}

	return step, nil
}

func mapToStruct(data map[string]any, out StepExecutor) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, out)
}
