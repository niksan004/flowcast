package clients

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	sshKey string
}

func NewSSHClient(sshKey string) *SSHClient {
	return &SSHClient{sshKey: sshKey}
}

func (c *SSHClient) Connect(host string, user string, port string) (*ssh.Client, error) {
	// handle port
	if port == "" {
		port = "22"
	}

	// handle 2 cases of private ssh key
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := c.sshKey
	if !strings.Contains(c.sshKey, "/") {
		path = filepath.Join(homeDir, ".ssh", c.sshKey)
	}

	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	cfg := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	return ssh.Dial("tcp", host+":"+port, cfg)
}
