package engine

import (
	"fmt"
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

// TODO: run this in parallel for each host
func (eng *Engine) Run() error {
	// iterate through hosts
	for _, host := range eng.cfg.inv.Hosts {
		client, err := eng.rt.SshCl.Connect(host.Ip, host.User, host.Port)
		if err != nil {
			return fmt.Errorf("Error while trying to connect to host: %s", err.Error())
		}
		logger.Info("{Connected to} " + logger.StringifyStruct(host))
		defer client.Close()

		// execute steps for each host
		for _, step := range eng.cfg.wf.Steps {
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

			// create session for step
			sesh, err := client.NewSession()
			if err != nil {
				return fmt.Errorf("Error while creating session for step: %s", step.Name)
			}

			logger.Info("{Executing} " + logger.StringifyStruct(exec))
			if err := exec.Execute(sesh); err != nil {
				sesh.Close()
				return err
			}
			sesh.Close()
		}
	}

	return nil
}
