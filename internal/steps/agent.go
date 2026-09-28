package steps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
)

const AgentPath = "/tmp/flowcast-agent"

func callAgent(sesh *ssh.Session, action string, args map[string]any) (any, error) {
	payload, err := json.Marshal(map[string]any{"action": action, "args": args})
	fmt.Println(string(payload))
	if err != nil {
		return nil, err
	}

	sesh.Stdin = bytes.NewReader(payload)

	var stdout bytes.Buffer
	sesh.Stdout = &stdout

	if err := sesh.Run(AgentPath); err != nil {
		return nil, err
	}

	var resp struct {
		Result any    `json:"result"`
		Error  string `json:"error"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &resp); err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	return resp.Result, nil
}
