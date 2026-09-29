package agent

import (
	"fmt"
	"project/internal/protocol"
)

func handleEcho(args protocol.EchoArgs) (any, error) {
	fmt.Println(args.Value)
	return args.Value, nil
}
