package steps

import (
	"golang.org/x/crypto/ssh"
	"project/internal/protocol"
)

type ShellStep struct {
	protocol.ShellArgs `yaml:",inline"`
}

func (step *ShellStep) Execute(sesh *ssh.Session) (any, error) {
	var res protocol.ShellResult
	if err := callAgent(sesh, &step.ShellArgs, &res); err != nil {
		return nil, err
	}
	return res, nil
}
