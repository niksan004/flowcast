package agent

import (
	"fmt"
	"project/internal/protocol"
)

type echo struct {
	protocol.EchoArgs
}

func (action *echo) Run() (any, error) {
	fmt.Println(action.Value)
	return action.Value, nil
}
