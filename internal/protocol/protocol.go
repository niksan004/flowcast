package protocol

import (
	"encoding/json"
	"net/http"
)

const AgentRemotePath = "/tmp/flowcast-agent"

type Request struct {
	Action string          `json:"action"`
	Args   json.RawMessage `json:"args"`
}

type Response struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}

// arg definitions
type Args interface {
	Action() string
}

type EchoArgs struct {
	Value string `json:"value" yaml:"value"`
}

func (*EchoArgs) Action() string { return "echo" }

type ShellArgs struct {
	Cmd string `json:"cmd" yaml:"cmd"`
}

func (*ShellArgs) Action() string { return "shell" }

type HttpArgs struct {
	Url    string      `json:"url" yaml:"url"`
	Method string      `json:"method" yaml:"method"`
	Body   string      `json:"body" yaml:"body"`
	Header http.Header `json:"header" yaml:"header"`
}

func (*HttpArgs) Action() string { return "http" }

// result definitions
type Result interface {
	isResult()
}

type EchoResult struct {
	Value string `json:"value"`
}

func (*EchoResult) isResult() {}

type ShellResult struct {
	Stdout   string `json:"stdout"`
	Stdin    string `json:"stdin"`
	ExitCode int    `json:"exitCode"`
}

func (*ShellResult) isResult() {}

type HttpResult struct {
	StatusCode int         `json:"statusCode"`
	Header     http.Header `json:"header"`
	Body       string      `json:"body"`
}

func (*HttpResult) isResult() {}
