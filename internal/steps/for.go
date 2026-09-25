package steps

import (
	"fmt"
	"github.com/expr-lang/expr"
	"golang.org/x/crypto/ssh"
)

type Looper interface {
	Loop(env map[string]any) ([]RawStep, error)
}

type ForStep struct {
	Condition string
	Do        []RawStep
}

// dummy function so this satisies StepExecutor
func (step *ForStep) Execute(sesh *ssh.Session) (any, error) {
	return nil, nil
}

func (step *ForStep) Loop(env map[string]any) ([]RawStep, error) {
	res, err := expr.Eval(step.Condition, env)
	if err != nil {
		return nil, err
	}

	val, ok := res.(bool)
	if !ok {
		return nil, fmt.Errorf("Condition did not evaluate to a bool: %v", res)
	}

	if val {
		return step.Do, nil
	}
	return nil, nil
}

func parseFor(data map[string]any) (StepExecutor, error) {
	var step ForStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return &step, nil
}
