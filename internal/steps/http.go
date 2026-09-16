package steps

import (
	"fmt"
	"project/internal/logger"
	"strings"
)

type HTTPStep struct {
	Url    string
	Method string
}

func (step *HTTPStep) Execute() error {
	switch strings.ToUpper(step.Method) {
	case "GET":
		logger.Info(fmt.Sprintf("GET req to: %s", step.Url))
	case "POST":
		logger.Info(fmt.Sprintf("POST req to: %s", step.Url))
	}
	return nil
}

func (step *HTTPStep) String() string {
	return fmt.Sprintf("==== Executing HTTPStep(Url=%s, Method=%s) ====", step.Url, step.Method)
}

func parseHTTP(data map[string]any) (StepExecutor, error) {
	var step HTTPStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return &step, nil
}
