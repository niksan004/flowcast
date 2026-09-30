package protocol

import "encoding/json"

const AgentRemotePath = "/tmp/flowcast-agent"

type Request struct {
	Action string          `json:"action"`
	Args   json.RawMessage `json:"args"`
}

type Response struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

type Args interface {
	Action() string
}

type EchoArgs struct {
	Value string `json:"value" yaml:"value"`
}

func (*EchoArgs) Action() string { return "echo" }

type Result interface {
	isResult()
}

type EchoResult struct {
	Value string `json:"value"`
}

func (*EchoResult) isResult() {}
