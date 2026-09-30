package agent

import (
	"fmt"
	"project/internal/protocol"
)

type echo struct {
	protocol.EchoArgs
}

func (action *echo) Run() (protocol.Result, error) {
	fmt.Println(action.Value)
	return &protocol.EchoResult{Value: action.Value}, nil
}
