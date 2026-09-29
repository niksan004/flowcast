package steps

import (
	"golang.org/x/crypto/ssh"
	"project/internal/protocol"
)

type EchoStep struct {
	Value string
}

func (step *EchoStep) Execute(sesh *ssh.Session) (any, error) {
	return callAgent(sesh, &protocol.EchoArgs{Value: step.Value})
}
