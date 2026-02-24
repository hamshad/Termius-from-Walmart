package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"termius-from-walmart/internal/models"
)

// NormalizePemKey converts escaped newlines, trims surrounding quotes/space,
// normalizes line endings, and ensures a trailing newline. This ensures PEM
// keys pasted from various sources (JSON, clipboard, single-line) work with
// OpenSSH.
func NormalizePemKey(pem string) string {
	clean := strings.TrimSpace(pem)

	// Handle JSON-style escaped newlines
	clean = strings.ReplaceAll(clean, "\\r\\n", "\n")
	clean = strings.ReplaceAll(clean, "\\n", "\n")

	// Normalize Windows line endings
	clean = strings.ReplaceAll(clean, "\r\n", "\n")

	// Remove surrounding quotes if pasted with them
	clean = strings.Trim(clean, "\"")

	// If the key is on a single line, re-wrap it into a proper PEM block
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

	// Ensure trailing newline for OpenSSH parser compatibility
	if clean != "" && !strings.HasSuffix(clean, "\n") {
		clean += "\n"
	}

	return clean
}

// BuildSSHCommand returns an *exec.Cmd configured for connecting to the given
// server via SSH. It handles PEM key auth, password auth (via sshpass), and
// fallback to system SSH keys.
func BuildSSHCommand(server models.Server) *exec.Cmd {
	args := []string{
		fmt.Sprintf("%s@%s", server.Username, server.Host),
		"-p", strconv.Itoa(server.Port),
	}

	// If PEM key is provided, save it to a temporary file
	if server.PemKey != "" {
		tempDir := filepath.Join(os.TempDir(), "termius-from-walmart-keys")
		os.MkdirAll(tempDir, 0700)

		keyFile := filepath.Join(tempDir, fmt.Sprintf("key_%d.pem", server.ID))
		cleanKey := NormalizePemKey(server.PemKey)

		if err := os.WriteFile(keyFile, []byte(cleanKey), 0600); err == nil {
			args = append([]string{
				"-i", keyFile,
				"-o", "StrictHostKeyChecking=no",
				"-o", "UserKnownHostsFile=/dev/null",
				"-o", "IdentitiesOnly=yes",
			}, args...)
			return exec.Command("ssh", args...)
		}
	}

	// If password is provided, use sshpass
	if server.Password != "" {
		return exec.Command("sshpass", append([]string{"-p", server.Password, "ssh"}, args...)...)
	}

	// Default: use system SSH keys
	return exec.Command("ssh", args...)
}
