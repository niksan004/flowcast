package steps

import (
	"fmt"
	"github.com/expr-lang/expr"
	"golang.org/x/crypto/ssh"
)

type Brancher interface {
	Branch() ([]RawStep, error)
}

type IfStep struct {
	Condition string
	Then      []RawStep
	Else      []RawStep
}

// dummy function so this satisies StepExecutor
func (step *IfStep) Execute(sesh *ssh.Session) error {
	return nil
}

func (step *IfStep) Branch() ([]RawStep, error) {
	res, err := expr.Eval(step.Condition, map[string]string{"dummy": "env"})
	if err != nil {
		return nil, err
	}

	val, ok := res.(bool)
	if !ok {
		return nil, fmt.Errorf("condition did not evaluate to a bool: %v", val)
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
