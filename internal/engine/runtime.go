package engine

import (
	"project/internal/clients"
)

type Runtime struct {
	HttpCl *clients.HTTPClient
	SshCl  *clients.SSHClient
}

func NewRuntime(sshKey string) (*Runtime, error) {
	return &Runtime{
		HttpCl: clients.NewHTTPClient(),
		SshCl:  clients.NewSSHClient(sshKey),
	}, nil
}
