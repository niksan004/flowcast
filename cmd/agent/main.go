package main

import (
	"encoding/json"
	"fmt"
	"os"
	"project/internal/agent"
	"project/internal/protocol"
)

func main() {
	var req protocol.Request[json.RawMessage]

	// receive stdin from controller
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Println("invalid input: %w", err)
	}

	// find correct handler
	handler, err := agent.Decode(req)
	if err != nil {
		fmt.Println("error wile decoding action %w", err)
	}

	res, err := handler.Run()

	// build json response
	var resp protocol.Response
	resp.Result = res
	if err != nil {
		resp.Error = err.Error()
	}

	// print new line to separate json output
	fmt.Fprintln(os.Stdout)

	if err := json.NewEncoder(os.Stdout).Encode(&resp); err != nil {
		fmt.Println("Invalid input: %w", err)
		os.Exit(1)
	}
}
