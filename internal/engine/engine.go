package engine

import (
	"fmt"
	"os"
	"project/internal/clients"
	"project/internal/logger"
	"project/internal/steps"
)

type Runtime struct {
	HttpCl *clients.HTTPClient
	SshCl  *clients.SSHClient
}

func NewRuntime() (*Runtime, error) {
	return &Runtime{
		HttpCl: clients.NewHTTPClient(),
		SshCl:  clients.NewSSHClient(),
	}, nil
}

type Config struct {
	wf *steps.Workflow
}

// TODO: read step by step not whole file
func readFile(filepath string) ([]byte, error) {
	file, err := os.ReadFile(filepath)
	return file, err
}

func NewConfig(filepath string) (*Config, error) {
	fileContent, err := readFile(filepath)
	if err != nil {
		return nil, err
	}

	wf, err := steps.UnmarshalWorkflow(fileContent)
	if err != nil {
		return nil, err
	}

	return &Config{wf: wf}, nil
}

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

func (eng *Engine) Run() error {
	// execute steps
	for _, step := range eng.cfg.wf.Steps {
		// get step factory for specific step
		fact, exists := steps.Registry[step.Name]
		if !exists {
			return fmt.Errorf("Unknown step: %s", step.Name)
		}

		// get specific step e.g. EchoStep
		exec, err := fact(step.Data)
		if err != nil {
			return fmt.Errorf("Error while executing step: %s", step.Name)
		}

		logger.Info(fmt.Sprintf("%s", exec))
		if err := exec.Execute(); err != nil {
			return err
		}
	}

	return nil
}
