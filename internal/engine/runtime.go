package engine

import (
	"project/internal/clients"
)

type Runtime struct {
	HttpCl *clients.HTTPClient
	SshCl  *clients.SSHClient
	SftpCl *clients.SFTPClient
}

func NewRuntime(sshKey string) (*Runtime, error) {
	return &Runtime{
		HttpCl: clients.NewHTTPClient(),
		SshCl:  clients.NewSSHClient(sshKey),
		SftpCl: clients.NewSFTPClient(),
	}, nil
}
