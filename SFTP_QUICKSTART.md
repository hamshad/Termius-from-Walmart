# SFTP Quick Start Guide

## Getting Started

### From Server List
1. Select a server from the list
2. Press **[s]** to open SFTP browser
3. The app connects automatically using your server credentials

## Basic Commands

### Navigation
| Key | Action |
|-----|--------|
| ↑/↓ | Move between files |
| [enter] | Enter directory |
| [../] | Go to parent directory |

### File Operations
| Key | Action |
|-----|--------|
| [c] | **Copy** file (two-step: select source, select destination dir) |
| [d] | **Delete** file |
| [r] | **Rename** file |

### General
| Key | Action |
|-----|--------|
| [q]/[esc] | Back to server list |

## Examples

### Copy a File
```
1. Press [c]
2. Select source file (e.g., config.txt)
3. Navigate to destination directory
4. Confirm (press [enter] on directory)
→ File copied!
```

### Delete a File
```
1. Navigate to file
2. Press [d]
3. File is deleted
```

### Rename a File
```
1. Press [r]
2. Select file to rename
3. Enter new filename
→ File renamed!
```

## Important Notes

- SFTP uses the same authentication as SSH (password or PEM key)
- SFTP port defaults to SSH port if not configured
- Operations happen on remote server only
- All file paths use forward slashes (Unix-style)
- Parent directory is shown as "../" 

## Configuration

When adding/editing servers:
- SSH Port: Default SSH connection port (22)
- SFTP Port: Separate SFTP port (usually same as SSH port)
  - Leave empty to use SSH port automatically

## Troubleshooting

### Connection Failed
- Check server credentials
- Ensure SFTP is enabled on remote server (usually port 22)
- Verify firewall rules

### Permission Denied
- Check user permissions on remote filesystem
- Some files/directories may be restricted

### Files Not Showing
- Directory may be empty
- Check if you have read permissions
- Try navigating to a different directory
