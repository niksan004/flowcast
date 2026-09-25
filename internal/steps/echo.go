package steps

import (
	"golang.org/x/crypto/ssh"
)

type EchoStep struct {
	Value string
}

func (step EchoStep) Execute(sesh *ssh.Session) (any, error) {
	return step.Value, nil
}

func parseEcho(data map[string]any) (StepExecutor, error) {
	var step EchoStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return step, nil
}
