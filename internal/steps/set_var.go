package steps

import (
	"golang.org/x/crypto/ssh"
)

type SetVarStep struct{}

func (step *SetVarStep) Execute(sesh *ssh.Session) (any, error) {
	return nil, nil
}
