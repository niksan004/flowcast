package steps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"project/internal/protocol"

	"golang.org/x/crypto/ssh"
)

func callAgent(sesh *ssh.Session, args protocol.Args, res protocol.Result) error {
	// create request to send to stdin of agent
	encodedArgs, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := json.Marshal(protocol.Request{
		Action: args.Action(),
		Args:   encodedArgs,
	})
	if err != nil {
		return err
	}
	sesh.Stdin = bytes.NewReader(req)

	var stdout, stderr bytes.Buffer
	sesh.Stdout = &stdout
	sesh.Stderr = &stderr

	// execute agent with prepared input
	if err := sesh.Run(protocol.AgentRemotePath); err != nil {
		return fmt.Errorf("agent failed: %w: stderr: %s", err, stderr.String())
	}

	// capture response
	var resp protocol.Response

	// json response is in stdout separated with a new line
	// from the other output and ending in a new line
	out := strings.TrimRight(stdout.String(), "\n\r")
	jsonLine := ""
	if i := strings.LastIndexByte(out, '\n'); i >= 0 {
		jsonLine = out[i+1:]
	}

	if err := json.Unmarshal([]byte(jsonLine), &resp); err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	if err := json.Unmarshal(resp.Result, res); err != nil {
		return err
	}
	return nil
}
