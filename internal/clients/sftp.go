package clients

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SftpConnect(conn *ssh.Client) (*sftp.Client, error) {
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	return client, nil
}
