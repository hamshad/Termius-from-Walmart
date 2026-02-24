package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"termius-from-walmart/internal/models"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPManager handles SFTP operations for the UI layer.
type SFTPManager struct {
	client *sftp.Client
	conn   *ssh.Client
}

// ConnectSFTP creates a new SFTP connection to the given server.
func ConnectSFTP(server *models.Server) (*SFTPManager, error) {
	config := &ssh.ClientConfig{
		User:            server.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if server.PemKey != "" {
		normalized := normalizePemForSSH(server.PemKey)
		signer, err := ssh.ParsePrivateKey([]byte(normalized))
		if err != nil {
			return nil, fmt.Errorf("failed to parse PEM key: %v", err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if server.Password != "" {
		config.Auth = []ssh.AuthMethod{ssh.Password(server.Password)}
	}

	port := server.SFTPPort
	if port == 0 {
		port = server.Port
		if port == 0 {
			port = 22
		}
	}

	sshAddr := fmt.Sprintf("%s:%d", server.Host, port)
	conn, err := ssh.Dial("tcp", sshAddr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create SFTP client: %v", err)
	}

	return &SFTPManager{
		client: client,
		conn:   conn,
	}, nil
}

// normalizePemForSSH reuses the ssh package's normalizer.
func normalizePemForSSH(pem string) string {
	// Inline normalization to avoid circular import
	clean := strings.TrimSpace(pem)
	clean = strings.ReplaceAll(clean, "\\r\\n", "\n")
	clean = strings.ReplaceAll(clean, "\\n", "\n")
	clean = strings.ReplaceAll(clean, "\r\n", "\n")
	clean = strings.Trim(clean, "\"")

	if strings.Contains(clean, "-----BEGIN") && strings.Contains(clean, "-----END") && !strings.Contains(clean, "\n") {
		beginIdx := strings.Index(clean, "-----BEGIN")
		endIdx := strings.Index(clean, "-----END")
		if beginIdx >= 0 && endIdx > beginIdx {
			headerCloseRel := strings.Index(clean[beginIdx+len("-----BEGIN"):], "-----")
			if headerCloseRel >= 0 {
				headerEnd := beginIdx + len("-----BEGIN") + headerCloseRel + len("-----")
				header := clean[beginIdx:headerEnd]
				footer := clean[endIdx:]
				body := strings.TrimSpace(clean[headerEnd:endIdx])
				body = strings.ReplaceAll(body, " ", "")
				body = strings.ReplaceAll(body, "\t", "")
				body = strings.ReplaceAll(body, "\r", "")
				body = strings.ReplaceAll(body, "\n", "")
				var wrapped []string
				for i := 0; i < len(body); i += 64 {
					end := i + 64
					if end > len(body) {
						end = len(body)
					}
					wrapped = append(wrapped, body[i:end])
				}
				clean = header + "\n" + strings.Join(wrapped, "\n") + "\n" + footer
			}
		}
	}

	if clean != "" && !strings.HasSuffix(clean, "\n") {
		clean += "\n"
	}

	return clean
}

// Close closes the SFTP and SSH connections.
func (sm *SFTPManager) Close() error {
	if sm.client != nil {
		sm.client.Close()
	}
	if sm.conn != nil {
		sm.conn.Close()
	}
	return nil
}

// ListFiles lists files in the given remote directory.
func (sm *SFTPManager) ListFiles(path string) ([]os.FileInfo, error) {
	return sm.client.ReadDir(path)
}

// UploadFile uploads a local file to the remote server.
func (sm *SFTPManager) UploadFile(localPath, remotePath string) error {
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer localFile.Close()

	remoteFile, err := sm.client.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %v", err)
	}
	defer remoteFile.Close()

	if _, err := io.Copy(remoteFile, localFile); err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// DownloadFile downloads a file from the remote server.
func (sm *SFTPManager) DownloadFile(remotePath, localPath string) error {
	remoteFile, err := sm.client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, remoteFile); err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	return nil
}

// DeleteFile deletes a file on the remote server.
func (sm *SFTPManager) DeleteFile(path string) error {
	return sm.client.Remove(path)
}

// homeDir returns the user's home directory.
func homeDir() string {
	return os.Getenv("HOME")
}
