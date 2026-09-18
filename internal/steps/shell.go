package steps

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"project/internal/logger"
)

type ShellStep struct {
	Cmd string
}

func (step *ShellStep) Execute(sesh *ssh.Session) error {
	var stdout, stderr bytes.Buffer
	sesh.Stdout = &stdout
	sesh.Stderr = &stderr

	err := sesh.Run(step.Cmd)
	logger.Info(stdout.String() + "\n" + stderr.String())
	return err
}

func (step *ShellStep) String() string {
	return logger.StringifyStruct(step)
}

func parseShell(data map[string]any) (StepExecutor, error) {
	var step ShellStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return &step, nil
}
