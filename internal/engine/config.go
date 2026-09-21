package engine

import (
	"gopkg.in/yaml.v3"
	"os"
)

// TODO: add variables to Config
type Config struct {
	wf  *Workflow
	inv *Inventory
}

type Inventory struct {
	Hosts []Host `yaml:"inventory"`
}

type Host struct {
	Ip   string `yaml:"host"`
	User string `yaml:"user"`
	Port string `yaml:"port"`
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

	wf, err := unmarshalWorkflow(fileContent)
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
