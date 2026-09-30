package agent

import (
	"encoding/json"
	"fmt"
	"project/internal/protocol"
)

type Runner interface {
	Run() (any, error)
}

type Handler func() Runner

// registry with constructors for action types
var Registry = map[string]Handler{
	(&protocol.EchoArgs{}).Action(): func() Runner { return &echo{} },
}

// get specific action from json
func Decode(raw protocol.Request[json.RawMessage]) (Runner, error) {
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
