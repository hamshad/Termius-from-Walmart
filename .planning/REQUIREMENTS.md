# Requirements: Termius from Walmart v1.0

## Functional Requirements

- [x] **CONN-01**: Add VPS/SSH connections with name, host, port, username, password, PEM key
- [x] **CONN-02**: Edit existing server connections
- [x] **CONN-03**: Delete server connections
- [x] **CONN-04**: Locally persist connections in JSON format at ~/.termius-from-walmart/config.json
- [x] **CONN-05**: Connect to servers via SSH (password, PEM key, or system SSH keys)
- [ ] **CONN-06**: Confirmation dialog before destructive actions (delete)
- [x] **IE-01**: Export all server configs to portable JSON file
- [x] **IE-02**: Import server configs from JSON file on another device
- [ ] **IE-03**: File picker for import/export with directory navigation
- [x] **SFTP-01**: SFTP split-screen file browser (local + remote)
- [x] **SFTP-02**: Upload files from local to remote
- [x] **SFTP-03**: Download files from remote to local
- [x] **SFTP-04**: Delete files on local and remote
- [x] **PEM-01**: Multiline PEM key editor with save/cancel
- [x] **PEM-02**: PEM key normalization (escaped newlines, line wrapping, format validation)

## Non-Functional Requirements

- [ ] **ARCH-01**: Modular codebase — separate packages for UI, SSH, storage, models
- [ ] **ARCH-02**: No deprecated standard library APIs (ioutil removed)
- [ ] **ARCH-03**: No debug logging in production code
- [ ] **TEST-01**: Unit tests for storage (load, save, config integrity)
- [ ] **TEST-02**: Unit tests for import/export (JSON parsing, edge cases)
- [ ] **TEST-03**: Unit tests for PEM key normalization
- [ ] **BUILD-01**: Makefile builds all .go files (not just main.go)
- [ ] **UX-01**: Compelling, aesthetic TUI with themed color palette
- [ ] **UX-02**: Responsive SFTP split-screen layout
- [ ] **UX-03**: Loading/connection state indicators
- [ ] **SEC-01**: Encrypted credential storage with master password
- [ ] **SEC-02**: Keychain integration (macOS Keychain / Linux secret-service)

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| CONN-01 | Pre-existing | Done |
| CONN-02 | Pre-existing | Done |
| CONN-03 | Pre-existing | Done |
| CONN-04 | Pre-existing | Done |
| CONN-05 | Pre-existing | Done |
| CONN-06 | Phase 01 | Pending |
| IE-01 | Pre-existing | Done |
| IE-02 | Pre-existing | Done |
| IE-03 | Pre-existing | Done |
| SFTP-01 | Pre-existing | Done |
| SFTP-02 | Pre-existing | Done |
| SFTP-03 | Pre-existing | Done |
| SFTP-04 | Pre-existing | Done |
| PEM-01 | Pre-existing | Done |
| PEM-02 | Pre-existing | Done |
| ARCH-01 | Phase 01 | Pending |
| ARCH-02 | Phase 01 | Pending |
| ARCH-03 | Phase 01 | Pending |
| TEST-01 | Phase 01 | Pending |
| TEST-02 | Phase 01 | Pending |
| TEST-03 | Phase 01 | Pending |
| BUILD-01 | Phase 01 | Pending |
| UX-01 | Phase 02 | Pending |
| UX-02 | Phase 02 | Pending |
| UX-03 | Phase 03 | Pending |
| SEC-01 | Phase 04 | Pending |
| SEC-02 | Phase 04 | Pending |
