package clients

import (
	"golang.org/x/crypto/ssh"
	"time"
)

type SSHClient struct {
	client *ssh.ClientConfig
}

func NewSSHClient() *SSHClient {
	// TODO add logging
	return &SSHClient{
		client: &ssh.ClientConfig{
			User: "",
			Auth: []ssh.AuthMethod{
				ssh.Password("yourpassword"),
				// or: ssh.PublicKeys(signer)
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Second * 5,
		},
	}
}
