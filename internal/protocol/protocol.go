package protocol

import (
// "encoding/json"
)

const AgentRemotePath = "/tmp/flowcast-agent"

type Request[T any] struct {
	Action string `json:"action"`
	Args   T      `json:"args"`
}

type Response struct {
	Result any    `json:"result"`
	Error  string `json:"error"`
}

type Args interface {
	Action() string
}

type EchoArgs struct {
	Value string `json:"value"`
}

func (*EchoArgs) Action() string { return "echo" }
