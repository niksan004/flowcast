package steps

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/ssh"
	"project/internal/logger"
)

type ShellStep struct {
	Cmd string
}

func (step *ShellStep) Execute(sesh *ssh.Session) (any, error) {
	var stdout, stderr bytes.Buffer
	sesh.Stdout = &stdout
	sesh.Stderr = &stderr
	exitCode := 0

	if err := sesh.Run(step.Cmd); err != nil {
		// check error type
		var exitErr *ssh.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitStatus()
		}
	}

	logger.Info(stdout.String() + "\n" + stderr.String())

	return map[string]any{
		"stdout":   stdout.String(),
		"stderr":   stderr.String(),
		"exitCode": exitCode,
	}, nil
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
