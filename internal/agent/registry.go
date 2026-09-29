package agent

import (
	"encoding/json"
	"project/internal/protocol"
)

func typed[T any](fn func(T) (any, error)) Handler {
	return func(raw json.RawMessage) (any, error) {
		var args T
		if err := json.Unmarshal(raw, &args); err != nil {
			return nil, err
		}
		return fn(args)
	}
}

type Handler func(args json.RawMessage) (any, error)

var Registry = map[string]Handler{
	(&protocol.EchoArgs{}).Action(): typed(handleEcho),
}
