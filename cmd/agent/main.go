package main

import (
	"encoding/json"
	"fmt"
	"os"
	"project/internal/agent"
)

func main() {
	var req struct {
		Action string
		Args   map[string]any
	}

	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Println("Invalid input: %w", err)
	}

	handler, ok := agent.Registry[req.Action]
	if !ok {
		fmt.Println("Invalid action %s", req.Action)
	}

	res, err := handler(req.Args)

	var resp struct {
		Result any    `json:"result"`
		Error  string `json:"error"`
	}

	resp.Result = res
	if err != nil {
		resp.Error = err.Error()
	}

	if err := json.NewEncoder(os.Stdout).Encode(&resp); err != nil {
		fmt.Println("Invalid input: %w", err)
		os.Exit(1)
	}
}
