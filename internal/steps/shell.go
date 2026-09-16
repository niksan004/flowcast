package steps

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"project/internal/logger"
)

type ShellStep struct {
	Cmd string
}

func (step *ShellStep) Execute() error {
	var out bytes.Buffer
	cmd := exec.Command("sh", "-c", step.Cmd)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	logger.Info("\n" + out.String())
	return err
}

func (step *ShellStep) String() string {
	return fmt.Sprintf("==== Executing ShellStep(Cmd=%s) ====", step.Cmd)
}

func parseShell(data map[string]any) (StepExecutor, error) {
	var step ShellStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return &step, nil
}
