package engine

import (
	"project/internal/clients"
)

type Runtime struct {
	HttpCl *clients.HTTPClient
	SshCl  *clients.SSHClient
}

func NewRuntime() (*Runtime, error) {
	return &Runtime{
		HttpCl: clients.NewHTTPClient(),
		SshCl:  clients.NewSSHClient(),
	}, nil
}
