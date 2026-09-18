package steps

import (
	"golang.org/x/crypto/ssh"
)

type StepExecutor interface {
	Execute(sesh *ssh.Session) error
}

type StepFactory func(data map[string]any) (StepExecutor, error)

var Registry = map[string]StepFactory{
	"echo":  parseEcho,
	"http":  parseHTTP,
	"shell": parseShell,
}
