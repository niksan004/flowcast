package steps

import (
	"fmt"
	"project/internal/logger"
)

type EchoStep struct {
	Value string
}

func (step EchoStep) Execute() error {
	logger.Info(step.Value)
	return nil
}

func (step EchoStep) String() string {
	return fmt.Sprintf("==== Executing EchoStep(Value=%s) ====", step.Value)
}

func parseEcho(data map[string]any) (StepExecutor, error) {
	var step EchoStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return step, nil
}
