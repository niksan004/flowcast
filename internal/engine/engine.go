package engine

import (
	"fmt"
	"github.com/expr-lang/expr"
	"golang.org/x/crypto/ssh"
	"project/internal/logger"
	"project/internal/steps"
)

type Engine struct {
	cfg *Config
	rt  *Runtime
}

func NewEngine(cfg *Config, rt *Runtime) *Engine {
	return &Engine{
		cfg: cfg,
		rt:  rt,
	}
}

// execute steps for each host
func RunSteps(stepsSlice []steps.RawStep, client *ssh.Client, env map[string]any) error {
	for _, step := range stepsSlice {
		// get step factory for specific step
		fact, exists := steps.Registry[step.Name]
		if !exists {
			return fmt.Errorf("Unknown step: %s", step.Name)
		}

		// get specific step(StepExecutor) e.g. EchoStep
		exec, err := fact(step.Data)
		if err != nil {
			return fmt.Errorf("Error while executing step: %s", step.Name)
		}

		logger.Info("{Executing} " + logger.StringifyStruct(exec))
		// check if there is branching in the workflow
		if brancher, isBrancher := exec.(steps.Brancher); isBrancher {
			nextSteps, err := brancher.Branch(env)
			if err != nil {
				return err
			}

			// recursively run the steps after branching
			if err := RunSteps(nextSteps, client, env); err != nil {
				return err
			}
		} else {
			// create session for step
			sesh, err := client.NewSession()
			if err != nil {
				return fmt.Errorf("Error while creating session for step: %s", step.Name)
			}

			// execute step
			returnVals, err := exec.Execute(sesh)
			if err != nil {
				sesh.Close()
				return err
			}
			sesh.Close()

			// save workflow variables if there are any
			if len(step.SaveAs) != 0 {
				for varName, path := range step.SaveAs {
					temp_env := map[string]any{"result": returnVals}
					val, err := expr.Eval(path, temp_env)
					if err != nil {
						return err
					}
					env[varName] = val
				}
			}
			logger.Info("Env: " + logger.StringifyStruct(env))
		}
	}

	return nil
}

// TODO: run this in parallel for each host
func (eng *Engine) RunEngine() error {
	// iterate through hosts
	for _, host := range eng.cfg.inv.Hosts {
		// environment
		env := map[string]any{}

		// client used for ssh
		client, err := eng.rt.SshCl.Connect(host.Ip, host.User, host.Port)
		if err != nil {
			return fmt.Errorf("Error while trying to connect to host: %s", err.Error())
		}
		logger.Info("{Connected to} " + logger.StringifyStruct(host))
		defer client.Close()

		// run steps
		if err := RunSteps(eng.cfg.wf.Steps, client, env); err != nil {
			return err
		}
	}

	return nil
}
