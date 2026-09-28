package agent

import ()

type Handler func(args map[string]any) (any, error)

var Registry = map[string]Handler{
	"echo": handleEcho,
}
