package steps

import (
	"bytes"
	"os/exec"
	"project/internal/logger"
)

type ShellStep struct {
	Cmd string
}

func (step *ShellStep) Execute() error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("sh", "-c", step.Cmd)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
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
