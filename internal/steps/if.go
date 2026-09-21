package steps

import (
	"fmt"
	"github.com/expr-lang/expr"
	"golang.org/x/crypto/ssh"
)

type Brancher interface {
	Branch(env map[string]any) ([]RawStep, error)
}

type IfStep struct {
	Condition string
	Then      []RawStep
	Else      []RawStep
}

// dummy function so this satisies StepExecutor
func (step *IfStep) Execute(sesh *ssh.Session) (any, error) {
	return nil, nil
}

func (step *IfStep) Branch(env map[string]any) ([]RawStep, error) {
	res, err := expr.Eval(step.Condition, env)
	if err != nil {
		return nil, err
	}

	val, ok := res.(bool)
	if !ok {
		return nil, fmt.Errorf("Condition did not evaluate to a bool: %v", val)
	}

	if val {
		return step.Then, nil
	}
	return step.Else, nil
}

func parseIf(data map[string]any) (StepExecutor, error) {
	var step IfStep

	if err := mapToStruct(data, &step); err != nil {
		return nil, err
	}

	return &step, nil
}
