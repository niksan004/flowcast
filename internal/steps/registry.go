package steps

import (
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

type RawStep struct {
	Name string         `yaml:"name"`
	Data map[string]any `yaml:",inline"`
}

type StepExecutor interface {
	Execute(sesh *ssh.Session) error
}

type StepFactory func(data map[string]any) (StepExecutor, error)

var Registry = map[string]StepFactory{
	"echo":  parseEcho,
	"http":  parseHTTP,
	"shell": parseShell,
	"if":    parseIf,
}

func mapToStruct(data map[string]any, out any) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, out)
}
