package steps

import (
	"golang.org/x/crypto/ssh"
	"project/internal/logger"
)

type EchoStep struct {
	Value string
}

func (step EchoStep) Execute(sesh *ssh.Session) (any, error) {
	logger.Info(step.Value)
	return nil, nil
}

func parseEcho(data map[string]any) (StepExecutor, error) {
	var step EchoStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return step, nil
}
