package steps

import (
	"golang.org/x/crypto/ssh"
	"project/internal/protocol"
)

type EchoStep struct {
	protocol.EchoArgs `yaml:",inline"`
}

func (step *EchoStep) Execute(sesh *ssh.Session) (any, error) {
	var res protocol.EchoResult
	if err := callAgent(sesh, &step.EchoArgs, &res); err != nil {
		return nil, err
	}
	return res, nil
}
