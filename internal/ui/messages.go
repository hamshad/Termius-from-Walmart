package ui

import (
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"termius-from-walmart/internal/models"
)

// ── SFTP Connection Messages ───────────────────────────────────────────

// SFTPConnectMsg is returned when an SFTP connection attempt completes.
type SFTPConnectMsg struct {
	Manager *SFTPManager
	Err     error
}

// connectSFTPCmd returns a tea.Cmd that connects to SFTP asynchronously.
func connectSFTPCmd(server *models.Server) tea.Cmd {
	return func() tea.Msg {
		mgr, err := ConnectSFTP(server)
		return SFTPConnectMsg{Manager: mgr, Err: err}
	}
}

// ── File Transfer Messages ─────────────────────────────────────────────

// TransferCompleteMsg is returned when a file transfer finishes.
type TransferCompleteMsg struct {
	FileName string
	Err      error
	IsUpload bool
}

// transferFileCmd performs a file transfer asynchronously.
func transferFileCmd(mgr *SFTPManager, srcFile, dstFile string, isUpload bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		if isUpload {
			err = mgr.UploadFile(srcFile, dstFile)
		} else {
			err = mgr.DownloadFile(srcFile, dstFile)
		}
		return TransferCompleteMsg{
			FileName: filepath.Base(srcFile),
			Err:      err,
			IsUpload: isUpload,
		}
	}
}

// ── Spinner Tick ───────────────────────────────────────────────────────

// SpinnerTickMsg drives the connecting spinner animation.
type SpinnerTickMsg time.Time

func spinnerTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinnerTickMsg(t)
	})
}

// ── Message Timer ──────────────────────────────────────────────────────

const secondDuration = time.Second

// ClearMessageMsg signals that the status message should be cleared.
type ClearMessageMsg struct{}

// clearMessageAfter returns a command that clears the message after a delay.
func clearMessageAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return ClearMessageMsg{}
	})
}
