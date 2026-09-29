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
	src string,
	dest string,
	mode os.FileMode,
) error {
	// sftp connection to copy agent to remote
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := client.Create(dest)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	defer out.Close()

	if err := client.Chmod(dest, 0755); err != nil {
		return err
	}

	return nil
}
