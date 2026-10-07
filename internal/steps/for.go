package steps

import (
	"fmt"
	"github.com/expr-lang/expr"
	"golang.org/x/crypto/ssh"
	"time"
)

const (
	defaultMaxIterations = 1000
)

type Looper interface {
	Loop(env map[string]any) ([]RawStep, error)
}

type ForStep struct {
	Condition     string        `yaml:"condition"`
	Do            []RawStep     `yaml:"do"`
	MaxIterations int           `yaml:"max_iterations"`
	Delay         time.Duration `yaml:"delay"`
}

// dummy function so this satisies StepExecutor
func (step *ForStep) Execute(sesh *ssh.Session) (any, error) {
	return nil, nil
}

func (step *ForStep) Run(env map[string]any, runBody func([]RawStep) error) error {
	if step.MaxIterations == 0 {
		step.MaxIterations = defaultMaxIterations
	}

	for i := 0; i < step.MaxIterations; i++ {
		res, err := expr.Eval(step.Condition, env)
		if err != nil {
			return err
		}

		val, ok := res.(bool)
		if !ok {
			return fmt.Errorf("condition did not evaluate to a bool: %v", res)
		}

		// break loop if condition if false
		if !val {
			break
		}

		// run loop if condition is true
		if err := runBody(step.Do); err != nil {
			return fmt.Errorf("error while running body: %w", err)
		}
		time.Sleep(step.Delay)
	}

	return nil
}
