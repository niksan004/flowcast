package agent

import (
	"bytes"
	"errors"
	"os/exec"
	"project/internal/protocol"
)

type shell struct {
	protocol.ShellArgs
}

func (action *shell) Run() (protocol.Result, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("sh", "-c", action.Cmd)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	exitCode := 0

	if err := cmd.Run(); err != nil {
		// check error type
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return nil, err
		}
		exitCode = exitErr.ExitCode()
	}

	return &protocol.ShellResult{
		Stdout:   stdout.String(),
		Stdin:    stderr.String(),
		ExitCode: exitCode,
	}, nil
}
