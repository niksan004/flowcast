package clients

import (
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

type SSHClient struct{}

func NewSSHClient() *SSHClient {
	return &SSHClient{}
}

func (c *SSHClient) Connect(host string, user string) (*ssh.Client, error) {
	key, err := os.ReadFile("../../.ssh/id_ed25519")
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)

	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	return ssh.Dial("tcp", host+":22", cfg)
}
