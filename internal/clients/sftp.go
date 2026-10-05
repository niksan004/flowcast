package clients

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
)

type SFTPClient struct{}

func NewSFTPClient() *SFTPClient {
	return &SFTPClient{}
}

func (sftpCl *SFTPClient) Upload(
	conn *ssh.Client,
	src io.Reader,
	destPath string,
	mode os.FileMode,
) error {
	// sftp connection to copy agent to remote
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	dest, err := client.Create(destPath)
	if err != nil {
		return err
	}

	// write with max concurrency
	if _, err := dest.ReadFromWithConcurrency(src, -1); err != nil {
		return err
	}
	defer dest.Close()

	if err := client.Chmod(destPath, 0755); err != nil {
		return err
	}

	return nil
}
