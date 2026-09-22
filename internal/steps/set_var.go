package steps

import (
	"golang.org/x/crypto/ssh"
)

type SetVarStep struct{}

func (step SetVarStep) Execute(sesh *ssh.Session) (any, error) {
	return nil, nil
}

func parseSetVar(data map[string]any) (StepExecutor, error) {
	var step SetVarStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return step, nil
}
