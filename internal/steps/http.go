package steps

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"strings"
)

type HTTPStep struct {
	Url    string
	Method string
}

func (step *HTTPStep) Execute(sesh *ssh.Session) (any, error) {
	switch strings.ToUpper(step.Method) {
	case "GET":
		return fmt.Sprintf("GET req to: %s", step.Url), nil
	case "POST":
		return fmt.Sprintf("POST req to: %s", step.Url), nil
	}
	return nil, nil
}
