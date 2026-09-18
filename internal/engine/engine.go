package engine

import (
	"fmt"
	"gopkg.in/yaml.v3"
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
	wf  *steps.Workflow
	inv *Inventory
}

type Inventory struct {
	Hosts []Host `yaml:"inventory"`
}

type Host struct {
	Ip   string `yaml:"host"`
	User string `yaml:"user"`
}

// TODO: read step by step not whole file
func readFile(filepath string) ([]byte, error) {
	file, err := os.ReadFile(filepath)
	return file, err
}

func unmarshalHosts(file []byte) (*Inventory, error) {
	// turn file into go structs
	var inv Inventory
	if err := yaml.Unmarshal(file, &inv); err != nil {
		return nil, err
	}

	return &inv, nil
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

	inv, err := unmarshalHosts(fileContent)
	if err != nil {
		return nil, err
	}

	return &Config{
		wf:  wf,
		inv: inv,
	}, nil
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
	// iterate through hosts
	for _, host := range eng.cfg.inv.Hosts {
		client, err := eng.rt.SshCl.Connect(host.Ip, host.User)
		if err != nil {
			return fmt.Errorf("Error while trying to connect to host: %s", err.Error())
		}
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

			logger.Info(fmt.Sprintf("%s", exec))
			if err := exec.Execute(sesh); err != nil {
				sesh.Close()
				return err
			}
			sesh.Close()
		}
	}

	return nil
}
