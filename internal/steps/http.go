package steps

import (
	"golang.org/x/crypto/ssh"
	"project/internal/protocol"
)

type HttpStep struct {
	protocol.HttpArgs `yaml:",inline"`
}

func (step *HttpStep) Execute(sesh *ssh.Session) (any, error) {
	var res protocol.HttpResult
	if err := callAgent(sesh, &step.HttpArgs, &res); err != nil {
		return nil, err
	}
	return res, nil
}
