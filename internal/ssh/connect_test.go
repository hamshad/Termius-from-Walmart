package ssh

import (
	"strings"
	"testing"

	"termius-from-walmart/internal/models"
)

func TestNormalizePemKey_EscapedNewlines(t *testing.T) {
	input := "-----BEGIN RSA PRIVATE KEY-----\\nMIIEpA...\\n-----END RSA PRIVATE KEY-----"
	result := NormalizePemKey(input)

	if !strings.Contains(result, "\n") {
		t.Error("Should convert escaped newlines to real newlines")
	}
	if strings.Contains(result, "\\n") {
		t.Error("Should not contain escaped newlines after normalization")
	}
}

func TestNormalizePemKey_WindowsLineEndings(t *testing.T) {
	input := "-----BEGIN RSA PRIVATE KEY-----\r\nMIIEpA...\r\n-----END RSA PRIVATE KEY-----"
	result := NormalizePemKey(input)

	if strings.Contains(result, "\r\n") {
		t.Error("Should normalize Windows line endings to Unix")
	}
}

func TestNormalizePemKey_TrailingNewline(t *testing.T) {
	input := "-----BEGIN RSA PRIVATE KEY-----\nMIIEpA...\n-----END RSA PRIVATE KEY-----"
	result := NormalizePemKey(input)

	if !strings.HasSuffix(result, "\n") {
		t.Error("Result should end with trailing newline for OpenSSH compatibility")
	}
}

func TestNormalizePemKey_SurroundingQuotes(t *testing.T) {
	input := `"-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----"`
	result := NormalizePemKey(input)

	if strings.HasPrefix(result, "\"") || strings.HasSuffix(strings.TrimSpace(result), "\"") {
		t.Error("Should remove surrounding quotes")
	}
}

func TestNormalizePemKey_EmptyString(t *testing.T) {
	result := NormalizePemKey("")
	if result != "" {
		t.Errorf("Empty input should produce empty output, got %q", result)
	}
}

func TestNormalizePemKey_WhitespaceOnly(t *testing.T) {
	result := NormalizePemKey("   \n\t  ")
	if result != "" {
		t.Errorf("Whitespace-only input should produce empty output, got %q", result)
	}
}

func TestNormalizePemKey_SingleLine(t *testing.T) {
	// Single-line PEM key (no newlines at all)
	input := "-----BEGIN RSA PRIVATE KEY-----MIIBogIBAAJBALRiMLAH-----END RSA PRIVATE KEY-----"
	result := NormalizePemKey(input)

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) < 3 {
		t.Errorf("Single-line PEM should be wrapped into multiple lines, got %d lines", len(lines))
	}
	if !strings.HasPrefix(lines[0], "-----BEGIN") {
		t.Error("First line should be the BEGIN header")
	}
	if !strings.HasPrefix(lines[len(lines)-1], "-----END") {
		t.Error("Last line should be the END footer")
	}
}

func TestNormalizePemKey_AlreadyValid(t *testing.T) {
	input := "-----BEGIN RSA PRIVATE KEY-----\nMIIBogIBAAJBALRi\nMLAH\n-----END RSA PRIVATE KEY-----\n"
	result := NormalizePemKey(input)

	if result != input {
		t.Errorf("Already-valid PEM should not be modified\nGot:  %q\nWant: %q", result, input)
	}
}

func TestBuildSSHCommand_DefaultSSHKeys(t *testing.T) {
	server := models.Server{
		ID:       1,
		Host:     "example.com",
		Port:     22,
		Username: "root",
	}

	cmd := BuildSSHCommand(server)
	args := cmd.Args

	if args[0] != "ssh" {
		t.Errorf("Command = %q, want %q", args[0], "ssh")
	}

	found := false
	for _, arg := range args {
		if arg == "root@example.com" {
			found = true
		}
	}
	if !found {
		t.Error("Command should contain user@host")
	}
}

func TestBuildSSHCommand_WithPassword(t *testing.T) {
	server := models.Server{
		ID:       2,
		Host:     "10.0.0.1",
		Port:     22,
		Username: "admin",
		Password: "mypassword",
	}

	cmd := BuildSSHCommand(server)
	if cmd.Args[0] != "sshpass" {
		t.Errorf("Command = %q, want %q (should use sshpass)", cmd.Args[0], "sshpass")
	}
}

func TestBuildSSHCommand_CustomPort(t *testing.T) {
	server := models.Server{
		ID:       3,
		Host:     "example.com",
		Port:     2222,
		Username: "user",
	}

	cmd := BuildSSHCommand(server)
	args := cmd.Args

	portFound := false
	for i, arg := range args {
		if arg == "-p" && i+1 < len(args) && args[i+1] == "2222" {
			portFound = true
		}
	}
	if !portFound {
		t.Error("Command should include -p 2222 for custom port")
	}
}

func TestBuildSSHCommand_WithPemKey(t *testing.T) {
	server := models.Server{
		ID:       4,
		Host:     "cloud.example.com",
		Port:     22,
		Username: "ec2-user",
		PemKey:   "-----BEGIN RSA PRIVATE KEY-----\nMIIBogIBAAJBALRi\n-----END RSA PRIVATE KEY-----",
	}

	cmd := BuildSSHCommand(server)
	args := cmd.Args

	if args[0] != "ssh" {
		t.Errorf("Command = %q, want %q", args[0], "ssh")
	}

	hasIdentity := false
	for _, arg := range args {
		if arg == "-i" {
			hasIdentity = true
		}
	}
	if !hasIdentity {
		t.Error("PEM key auth should include -i flag")
	}
}
