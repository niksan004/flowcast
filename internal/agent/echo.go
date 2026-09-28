package agent

import ()

func handleEcho(args map[string]any) (any, error) {
	value, _ := args["value"].(string)
	// fmt.Println(value)
	return value, nil
}
