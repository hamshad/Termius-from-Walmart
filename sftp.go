package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPManager handles SFTP operations
type SFTPManager struct {
	client *sftp.Client
	conn   *ssh.Client
}

// ConnectSFTP creates a new SFTP connection
func ConnectSFTP(server *Server) (*SFTPManager, error) {
	// Setup SSH config
	config := &ssh.ClientConfig{
		User:            server.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Warning: insecure for production
	}

	// Use PEM key if available
	if server.PemKey != "" {
		normalized := normalizePemKey(server.PemKey)
		signer, err := ssh.ParsePrivateKey([]byte(normalized))
		if err != nil {
			return nil, fmt.Errorf("failed to parse PEM key: %v", err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else if server.Password != "" {
		// Use password authentication
		config.Auth = []ssh.AuthMethod{ssh.Password(server.Password)}
	}

	// Determine SFTP port
	port := server.SFTPPort
	if port == 0 {
		port = server.Port
		if port == 0 {
			port = 22
		}
	}

	// Connect to SSH server
	sshAddr := fmt.Sprintf("%s:%d", server.Host, port)
	conn, err := ssh.Dial("tcp", sshAddr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
	}

	// Create SFTP client
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

// Close closes the SFTP connection
func (sm *SFTPManager) Close() error {
	if sm.client != nil {
		sm.client.Close()
	}
	if sm.conn != nil {
		sm.conn.Close()
	}
	return nil
}

// ListFiles lists files in a directory
func (sm *SFTPManager) ListFiles(path string) ([]os.FileInfo, error) {
	return sm.client.ReadDir(path)
}

// CopyFile copies a file from source to destination on remote server
func (sm *SFTPManager) CopyFile(src, dst string) error {
	// Open source file
	srcFile, err := sm.client.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := sm.client.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	// Copy content
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}

	return nil
}

// DeleteFile deletes a file
func (sm *SFTPManager) DeleteFile(path string) error {
	return sm.client.Remove(path)
}

// RenameFile renames a file
func (sm *SFTPManager) RenameFile(oldPath, newPath string) error {
	return sm.client.Rename(oldPath, newPath)
}

// CreateDirectory creates a directory
func (sm *SFTPManager) CreateDirectory(path string) error {
	return sm.client.Mkdir(path)
}

// GetWorkingDirectory returns the current directory
func (sm *SFTPManager) GetWorkingDirectory() (string, error) {
	return sm.client.Getwd()
}

// ChangeDirectory changes the current directory
func (sm *SFTPManager) ChangeDirectory(path string) error {
	_, err := sm.client.Stat(path)
	return err
}

// UploadFile uploads a local file to the remote server
func (sm *SFTPManager) UploadFile(localPath, remotePath string) error {
	// Open local file
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer localFile.Close()

	// Create remote file
	remoteFile, err := sm.client.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %v", err)
	}
	defer remoteFile.Close()

	// Copy content
	if _, err := io.Copy(remoteFile, localFile); err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// DownloadFile downloads a file from the remote server
func (sm *SFTPManager) DownloadFile(remotePath, localPath string) error {
	// Open remote file
	remoteFile, err := sm.client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer remoteFile.Close()

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	// Copy content
	if _, err := io.Copy(localFile, remoteFile); err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	return nil
}

// FormatFileList returns formatted file list for display
func FormatFileList(files []os.FileInfo) []string {
	var formatted []string

	for _, f := range files {
		if f.IsDir() {
			formatted = append(formatted, f.Name()+"/")
		} else {
			formatted = append(formatted, f.Name())
		}
	}

	return formatted
}

// TrimPathSlash removes trailing slash from path
func TrimPathSlash(path string) string {
	return strings.TrimSuffix(path, "/")
}

// JoinPath safely joins path components
func JoinPath(base, element string) string {
	if strings.HasPrefix(element, "/") {
		return element
	}
	if base == "/" {
		return "/" + element
	}
	return base + "/" + element
}

// GetParentPath returns the parent directory path
func GetParentPath(path string) string {
	if path == "/" {
		return "/"
	}
	parentPath := filepath.Dir(path)
	if !strings.HasSuffix(parentPath, "/") && parentPath != "/" {
		parentPath += "/"
	}
	return parentPath
}
