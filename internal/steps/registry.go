package steps

import (
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

type StepFactory func(data map[string]any) (StepExecutor, error)

func parse[T any, PT interface {
	*T
	StepExecutor
}]() StepFactory {
	return func(data map[string]any) (StepExecutor, error) {
		var step T

		if err := mapToStruct(data, &step); err != nil {
			return nil, err
		}

		return PT(&step), nil
	}
}

var Registry = map[string]StepFactory{
	"set_var": parse[SetVarStep](),
	"echo":    parse[EchoStep](),
	"http":    parse[HTTPStep](),
	"shell":   parse[ShellStep](),
	"if":      parse[IfStep](),
	"for":     parse[ForStep](),
}

func mapToStruct(data map[string]any, out any) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, out)
}
