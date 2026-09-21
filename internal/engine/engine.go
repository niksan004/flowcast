package engine

import (
	"fmt"
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
func RunSteps(stepsSlice []steps.RawStep, client *ssh.Client) error {
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
			nextSteps, err := brancher.Branch()
			if err != nil {
				return err
			}

			// recursively run the steps after branching
			if err := RunSteps(nextSteps, client); err != nil {
				return err
			}
		} else {
			// create session for step
			sesh, err := client.NewSession()
			if err != nil {
				return fmt.Errorf("Error while creating session for step: %s", step.Name)
			}

			// execute step
			if err := exec.Execute(sesh); err != nil {
				sesh.Close()
				return err
			}
			sesh.Close()
		}
	}

	return nil
}

// TODO: run this in parallel for each host
func (eng *Engine) RunEngine() error {
	// iterate through hosts
	for _, host := range eng.cfg.inv.Hosts {
		client, err := eng.rt.SshCl.Connect(host.Ip, host.User, host.Port)
		if err != nil {
			return fmt.Errorf("Error while trying to connect to host: %s", err.Error())
		}
		logger.Info("{Connected to} " + logger.StringifyStruct(host))
		defer client.Close()

		if err := RunSteps(eng.cfg.wf.Steps, client); err != nil {
			return err
		}
	}

	return nil
}
