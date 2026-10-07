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

func (step *IfStep) Run(env map[string]any, runBody func([]RawStep) error) error {
	res, err := expr.Eval(step.Condition, env)
	if err != nil {
		return err
	}

	val, ok := res.(bool)
	if !ok {
		return fmt.Errorf("Condition did not evaluate to a bool: %v", res)
	}

	if val {
		return runBody(step.Then)
	}
	return runBody(step.Else)
}
