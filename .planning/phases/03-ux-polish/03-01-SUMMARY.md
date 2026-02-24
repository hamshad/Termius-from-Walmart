---
phase: "03"
plan: "03-01"
subsystem: ui
tags: [ux, async, spinner, metadata, help-overlay]
dependency-graph:
  requires: [02-01]
  provides: [async-sftp, file-metadata, help-overlay, message-timer]
  affects: [model, views, file-item, sftp]
tech-stack:
  added: []
  patterns: [tea.Cmd async, tick-based spinner, message generation counter]
key-files:
  created:
    - internal/ui/messages.go
  modified:
    - internal/ui/model.go
    - internal/ui/view_list.go
    - internal/ui/view_sftp.go
    - internal/ui/views.go
    - internal/ui/file_item.go
    - internal/ui/view_menu.go
decisions:
  - Spinner uses manual tick (80ms) with custom frame array instead of bubbles/spinner to avoid extra model nesting
  - Message auto-clear uses generation counter pattern to avoid race between rapid messages
  - FileItem changed from string alias to struct with metadata fields
  - FileDelegate.ShowMetadata flag distinguishes SFTP lists (show metadata) from file picker lists (no metadata)
  - Help overlay reuses ModalStyle with wider width (55 chars) for readability
metrics:
  duration: ~5min
  completed: 2026-02-24
---

# Phase 03 Plan 01: UX Polish Summary

Async SFTP connections with spinner, SSH feedback, file metadata display, non-blocking file transfers, keyboard shortcut help overlay, and timed message auto-clear.

## What Was Done

### Task 1: Async SFTP Connection with Spinner
- Added `ConnectingView` state with animated braille spinner (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
- SFTP connection now runs asynchronously via `connectSFTPCmd` returning `SFTPConnectMsg`
- Spinner ticks at 80ms via `SpinnerTickMsg`
- Connection can be cancelled with Esc during connecting phase
- On success: transitions to SFTPView with "Connected to {server}" message
- On failure: returns to ListView with error message (5s auto-clear)

### Task 2: SSH Connection Feedback
- Added "Connecting to {user}@{host}:{port}..." message displayed before `tea.ExecProcess` hand-off
- Provides visual feedback that connection is being established

### Task 3: File Metadata in SFTP
- Rewrote `FileItem` from `type FileItem string` to a struct with Name, Size, Mode, ModTime, IsDir fields
- Added `NewFileItemFromInfo` (for SFTP remote files from os.FileInfo) and `NewFileItemFromDirEntry` (for local files from os.DirEntry)
- `FileDelegate.ShowMetadata` flag enables 2-line rendering with size, permissions, and mod time
- SFTP local/remote lists show metadata; file picker lists do not
- Added `humanizeBytes` helper (e.g., "4.2K", "1.3M")

### Task 4: Async File Transfers
- `performCopy` now returns `tea.Cmd` instead of blocking
- Transfer runs in background via `transferFileCmd` goroutine
- `TransferCompleteMsg` carries result back to Update loop
- Progress indicator shows during transfer, file lists refresh on completion
- Duplicate transfers blocked (checks `IsTransferring` flag)

### Task 5: Keyboard Shortcut Help Overlay
- New `HelpView` state toggled by `?` key from any non-input view
- Organized by sections: Server List, SFTP Browser, Forms & Editors, Global
- Rendered as centered modal overlay using ModalStyle
- Closed with `?`, `esc`, or `q`

### Task 6: Timed Message Auto-Clear
- Messages auto-clear after 3 seconds (success) or 5 seconds (errors)
- Uses generation counter pattern (`MessageTimer`) to prevent newer messages from being cleared by older timers
- `ClearMessageMsg` only clears if no newer message has been posted

## New File

- `internal/ui/messages.go` — All custom tea.Msg types (SFTPConnectMsg, TransferCompleteMsg, SpinnerTickMsg, ClearMessageMsg) and their associated tea.Cmd constructors

## Deviations from Plan

None — plan executed exactly as written.

## Commits

| Hash | Message |
|------|---------|
| cfd7456 | feat(03-01): async SFTP connection with spinner, SSH feedback, file metadata, async transfers, help overlay, message auto-clear |

## Verification

- `go build -v .` — compiles without errors
- `go test ./...` — all 28 tests pass (models: 5, storage: 11, ssh: 12)
