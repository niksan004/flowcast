package agent

import (
	"encoding/json"
	"fmt"
	"project/internal/protocol"
)

type Runner interface {
	Run() (protocol.Result, error)
}

type Handler func() Runner

// registry with constructors for action types
var Registry = map[string]Handler{
	(&protocol.EchoArgs{}).Action():  func() Runner { return &echo{} },
	(&protocol.ShellArgs{}).Action(): func() Runner { return &shell{} },
	(&protocol.HttpArgs{}).Action():  func() Runner { return &httpAc{} },
}

// get specific action from json
func Decode(raw protocol.Request) (Runner, error) {
	actionConstr, exists := Registry[raw.Action]
	if !exists {
		return nil, fmt.Errorf("unknown action %s", raw.Action)
	}

	action := actionConstr()
	if err := json.Unmarshal(raw.Args, &action); err != nil {
		return nil, err
	}

	return action, nil
}
