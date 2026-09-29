package steps

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/crypto/ssh"
	"project/internal/protocol"
)

func callAgent(sesh *ssh.Session, args protocol.Args) (protocol.Response, error) {
	// create request to send to stdin of agent
	req, err := json.Marshal(protocol.Request[protocol.Args]{
		Action: args.Action(),
		Args:   args,
	})
	if err != nil {
		return protocol.Response{}, err
	}
	sesh.Stdin = bytes.NewReader(req)

	var stdout bytes.Buffer
	sesh.Stdout = &stdout

	// execute agent with prepared input
	if err := sesh.Run(protocol.AgentRemotePath); err != nil {
		return protocol.Response{}, err
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
		return protocol.Response{}, err
	}

	if resp.Error != "" {
		return protocol.Response{}, errors.New(resp.Error)
	}
	return resp, nil
}
