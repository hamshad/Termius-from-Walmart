# Phase 01 Plan 01: Restructure & Harden Summary

## One-liner
Modularized 1613-line monolith into 4 internal packages with 28 passing tests and confirmation dialogs

## Tasks Completed

| # | Task | Commit | Key Files |
|---|------|--------|-----------|
| 1 | Split monolith into modular packages | 7236b99 | internal/models/server.go, internal/storage/config.go, internal/ssh/connect.go, internal/ui/*.go |
| 2 | Fix deprecated ioutil APIs + remove debug logging | 7236b99 | (handled during restructure — all new files use os.ReadFile/WriteFile/ReadDir) |
| 3 | Fix Makefile build command | 9bbb298 | Makefile |
| 4 | Add unit tests for storage, SSH, PEM, import/export | 53adb7c | internal/models/server_test.go, internal/storage/config_test.go, internal/ssh/connect_test.go |
| 5 | Add confirmation dialogs for destructive actions | 7236b99 | internal/ui/views.go (ConfirmView), internal/ui/model.go (ShowConfirm) |

## Architecture Changes

### Before
```
main.go       — 1613 lines (all UI, logic, storage, SSH)
sftp.go       — 225 lines
storage.go    — 48 lines
```

### After
```
main.go                          — 18 lines (entry point only)
internal/models/server.go        — 39 lines (data types)
internal/models/server_test.go   — 50 lines (5 tests)
internal/storage/config.go       — 118 lines (load, save, import, export)
internal/storage/config_test.go  — 216 lines (11 tests)
internal/ssh/connect.go          — 99 lines (PEM normalization, SSH command builder)
internal/ssh/connect_test.go     — 168 lines (12 tests)
internal/ui/model.go             — 272 lines (model, init, update dispatch, view dispatch)
internal/ui/styles.go            — 28 lines (shared lipgloss styles)
internal/ui/file_item.go         — 42 lines (file picker item + delegate)
internal/ui/view_list.go         — 103 lines (server list view)
internal/ui/view_form.go         — 207 lines (add/edit form + validation)
internal/ui/view_menu.go         — 197 lines (import/export menu + file picker)
internal/ui/sftp.go              — 170 lines (SFTP connection manager)
internal/ui/view_sftp.go         — 270 lines (SFTP split-screen view)
internal/ui/views.go             — 260 lines (all render functions + confirm dialog)
```

## Test Coverage

| Package | Tests | Status |
|---------|-------|--------|
| internal/models | 5 | PASS |
| internal/storage | 11 | PASS |
| internal/ssh | 12 | PASS |
| **Total** | **28** | **ALL PASS** |

### Key Test Scenarios
- Config load/save round-trip with file permissions verification
- Import/Export JSON round-trip (data integrity preserved across devices)
- PEM key normalization: escaped newlines, Windows CRLF, single-line rewrap, surrounding quotes, empty input
- SSH command builder: password auth (sshpass), PEM key auth (-i flag), custom port, default keys

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] .gitignore was blocking all .md and .json files**
- **Found during:** Planning infrastructure creation
- **Issue:** `.gitignore` had `*.md` and `config.json` patterns that blocked all markdown and json files
- **Fix:** Updated .gitignore to only ignore build artifacts and runtime-specific files
- **Commit:** 8a72030

**2. [Rule 2 - Missing Functionality] Duplicate PEM normalization in SFTP**
- **Found during:** Task 1 (restructure)
- **Issue:** SFTP connection needed PEM normalization but importing from internal/ssh would create cleaner architecture
- **Fix:** Inlined minimal PEM normalization in ui/sftp.go to avoid circular dependency, kept full version in internal/ssh
- **Commit:** 7236b99

## Decisions Made

- **Package structure:** `internal/models`, `internal/storage`, `internal/ssh`, `internal/ui` — standard Go internal package pattern
- **Confirmation dialog:** Default cursor on "No" (safe default) with left/right arrow navigation
- **SFTP manager stays in ui package:** It manages connection lifecycle tied to UI state, not reusable independently

## Verification

- `go build .` — PASS
- `go vet ./...` — PASS  
- `go test ./...` — 28/28 PASS
- `make build` — PASS

## Self-Check: PASSED
